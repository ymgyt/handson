extern crate alloc;

use crate::{hpet::global_timestamp, info, result::Result, x86::busy_loop_hint};
use core::{
    fmt,
    future::Future,
    panic::Location,
    pin::Pin,
    ptr::null,
    sync::atomic::AtomicBool,
    task::{Context, Poll, RawWaker, RawWakerVTable, Waker},
    time::Duration,
};

use alloc::{boxed::Box, collections::VecDeque};

pub struct Task<T> {
    future: Pin<Box<dyn Future<Output = Result<T>>>>,
    created_at_file: &'static str,
    created_at_line: u32,
}

impl<T> Task<T> {
    #[track_caller]
    pub fn new(future: impl Future<Output = Result<T>> + 'static) -> Task<T> {
        Task {
            future: Box::pin(future),
            created_at_file: Location::caller().file(),
            created_at_line: Location::caller().line(),
        }
    }
    fn poll(&mut self, context: &mut Context) -> Poll<Result<T>> {
        self.future.as_mut().poll(context)
    }
}
impl<T> fmt::Debug for Task<T> {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "Task({}:{})", self.created_at_file, self.created_at_line)
    }
}

fn no_op_raw_waker() -> RawWaker {
    fn no_op(_: *const ()) {}
    fn clone(_: *const ()) -> RawWaker {
        no_op_raw_waker()
    }
    let vtable = &RawWakerVTable::new(clone, no_op, no_op, no_op);
    RawWaker::new(null::<()>(), vtable)
}

pub fn no_op_waker() -> Waker {
    unsafe { Waker::from_raw(no_op_raw_waker()) }
}

pub fn block_on<T>(future: impl Future<Output = Result<T>> + 'static) -> Result<T> {
    let mut task = Task::new(future);
    loop {
        let waker = no_op_waker();
        let mut context = Context::from_waker(&waker);
        match task.poll(&mut context) {
            Poll::Ready(result) => {
                break result;
            }
            Poll::Pending => busy_loop_hint(),
        }
    }
}

pub struct Executor {
    task_queue: Option<VecDeque<Task<()>>>,
}

impl Executor {
    pub const fn new() -> Self {
        Self { task_queue: None }
    }
    fn task_queue(&mut self) -> &mut VecDeque<Task<()>> {
        if self.task_queue.is_none() {
            self.task_queue = Some(VecDeque::new());
        }
        self.task_queue.as_mut().unwrap()
    }
    pub fn enqueue(&mut self, task: Task<()>) {
        self.task_queue().push_back(task)
    }
    pub fn run(mut executor: Self) -> ! {
        info!("Executor status running...");
        loop {
            let task = executor.task_queue().pop_front();
            if let Some(mut task) = task {
                let waker = no_op_waker();
                let mut context = Context::from_waker(&waker);
                match task.poll(&mut context) {
                    Poll::Ready(result) => {
                        info!("Task completed: {:?}: {:?}", task, result);
                    }
                    Poll::Pending => {
                        executor.task_queue().push_back(task);
                    }
                }
            }
        }
    }
}

impl Default for Executor {
    fn default() -> Self {
        Self::new()
    }
}

#[derive(Default)]
pub struct Yield {
    polled: AtomicBool,
}
impl Future for Yield {
    type Output = ();

    fn poll(self: Pin<&mut Self>, _cx: &mut Context<'_>) -> Poll<Self::Output> {
        if self
            .polled
            .fetch_or(true, core::sync::atomic::Ordering::Relaxed)
        {
            Poll::Ready(())
        } else {
            Poll::Pending
        }
    }
}

pub async fn yield_execution() {
    Yield::default().await
}

pub struct TimeoutFuture {
    timeout: Duration,
}
impl TimeoutFuture {
    pub fn new(duration: Duration) -> Self {
        Self {
            timeout: global_timestamp() + duration,
        }
    }
}

impl Future for TimeoutFuture {
    type Output = ();

    fn poll(self: Pin<&mut Self>, _cx: &mut Context<'_>) -> Poll<Self::Output> {
        if self.timeout < global_timestamp() {
            Poll::Ready(())
        } else {
            Poll::Pending
        }
    }
}

use std::{
    cell::UnsafeCell,
    mem::MaybeUninit,
    sync::{
        atomic::{AtomicU8, Ordering::*},
        Arc,
    },
    thread,
};

static NOT_READY: u8 = 0;
static READY: u8 = 1;

struct Channel<T> {
    message: UnsafeCell<MaybeUninit<T>>,
    state: AtomicU8,
}

struct Sender<T> {
    ch: Arc<Channel<T>>,
}

struct Receiver<T> {
    ch: Arc<Channel<T>>,
}

impl<T> Sender<T> {
    fn send(self, message: T) {
        if self.ch.state.load(Acquire) != NOT_READY {
            panic!("invalid state");
        }
        unsafe {
            (*self.ch.message.get()).write(message);
        }
        self.ch.state.store(READY, Release);
    }
}

impl<T> Receiver<T> {
    fn receive(self) -> T {
        if self.ch.state.load(Acquire) != READY {
            panic!("invalid state")
        }
        unsafe { (*self.ch.message.get()).assume_init_read() }
    }
}

unsafe impl<T> Sync for Channel<T> where T: Send {}

impl<T> Channel<T> {
    fn new() -> Self {
        Self {
            message: UnsafeCell::new(MaybeUninit::uninit()),
            state: AtomicU8::new(NOT_READY),
        }
    }
    fn split() -> (Sender<T>, Receiver<T>) {
        let ch = Channel::new();
        let ch = Arc::new(ch);
        (Sender { ch: ch.clone() }, Receiver { ch: ch })
    }
}

fn main() {
    let (sender, receiver) = Channel::split();
    thread::scope(move |s| {
        s.spawn(move || {
            sender.send("Hello");
        });
    });

    let m = receiver.receive();
    println!("{m}");
}

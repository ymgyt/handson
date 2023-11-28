use std::{
    sync::atomic::{AtomicU32, Ordering::*},
    thread,
    time::Duration,
};

#[cfg(not(target_os = "linux"))]
compile_error!("Linux only.");

fn wait(a: &AtomicU32, expected: u32) {
    unsafe {
        libc::syscall(
            libc::SYS_futex,
            a as *const AtomicU32,
            libc::FUTEX_WAIT,
            expected,
            std::ptr::null::<libc::timespec>(),
        );
    }
}

fn wake_one(a: &AtomicU32) {
    unsafe {
        libc::syscall(libc::SYS_futex, a as *const AtomicU32, libc::FUTEX_WAKE, 1);
    }
}
fn main() {
    let a = AtomicU32::new(0);

    thread::scope(|s| {
        s.spawn(|| {
            thread::sleep(Duration::from_secs(3));
            a.store(1, Relaxed);
            wake_one(&a);
        });

        println!("Waiting...");
        while a.load(Relaxed) == 0 {
            wait(&a, 0);
        }
        println!("Done!");
    });
}

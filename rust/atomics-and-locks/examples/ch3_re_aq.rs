use std::{
    sync::atomic::{AtomicBool, AtomicU32, Ordering::*},
    thread,
    time::Duration,
};

static DATA: AtomicU32 = AtomicU32::new(0);
static READY: AtomicBool = AtomicBool::new(false);

fn main() {
    thread::spawn(|| {
        DATA.store(123, Relaxed);
        READY.store(true, Release);
    });

    while !READY.load(Acquire) {
        thread::sleep(Duration::from_millis(100));
        println!("waiting...");
    }

    println!("{}", DATA.load(Relaxed));
    assert_eq!(DATA.load(Relaxed), 123);
}

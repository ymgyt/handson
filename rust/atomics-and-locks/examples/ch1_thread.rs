use std::thread;

fn main() {
    let t1 = thread::spawn(f);
    let t2 = thread::spawn(f);

    println!("Hello from main thread");

    t1.join().unwrap();
    t2.join().unwrap();
}

fn f() {
    let id = thread::current().id();
    println!("My threadId: {id:?}");
}

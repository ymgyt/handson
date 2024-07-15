use builder_macro::Builder;

#[derive(Builder)]
struct Example {
    name: String,
}

fn main() {
    println!("Hello, world!");
}

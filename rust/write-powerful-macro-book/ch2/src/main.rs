fn add_one(n: i32) -> i32 {
    n + 1
}

fn stringify(n: i32) -> String {
    n.to_string()
}

fn prefix_with<'a>(prefix: &'a str) -> impl Fn(String) -> String + 'a {
    move |x| format!("{}{}", prefix, x)
}

mod macros;
use macros::compose;

fn main() {
    let composed = compose!(add_one, stringify, prefix_with("Result: "));

    println!("{}", composed(3));
}

use private_macro::compose;
use private_macro::private;

private! {
    struct Example {
        string_value: String,
        number_value: i32,
    }
}

fn add_one(n: i32) -> i32 {
    n + 1
}

fn to_string(n: i32) -> String {
    n.to_string()
}

fn with_prefix(s: String) -> String {
    format!("prefix_{s}")
}

fn main() {
    let composed = compose!(add_one.to_string.with_prefix);

    println!("{}", composed(2));
}

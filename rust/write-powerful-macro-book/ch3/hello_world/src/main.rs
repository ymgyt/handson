use hello_world_macro::Hello;

#[derive(Hello)]
struct Example {}

fn main() {
    let e = Example {};
    e.hello_world();
}

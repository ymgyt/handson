use cgp::prelude::*;

#[cgp_component {
    name: HelloComponent,
    provider: HelloProvider
}]
trait CanHello {
    fn hello(&self) -> String;
}

mod impls {
    use super::*;
    pub struct DebugHello;

    impl<Context> HelloProvider<Context> for DebugHello
    where
        Context: std::fmt::Debug,
    {
        fn hello(context: &Context) -> String {
            format!("Hello {context:?}")
        }
    }

    pub struct DisplayHello;

    impl<Context> HelloProvider<Context> for DisplayHello
    where
        Context: std::fmt::Display,
    {
        fn hello(context: &Context) -> String {
            format!("Hello {context}")
        }
    }
}

#[derive(Debug)]
struct User {
    name: &'static str,
}
impl std::fmt::Display for User {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        f.write_str(self.name)
    }
}

struct UserComponents;

impl HasComponents for User {
    type Components = UserComponents;
}

delegate_components! {
    UserComponents {
        HelloComponent: impls::DisplayHello,
    }
}

fn main() {
    let user = User { name: "ymgyt" };

    println!("{}", user.hello());
}

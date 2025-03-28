use proc_macro::TokenStream;
use quote::quote;
use syn::DeriveInput;

#[proc_macro_derive(Hello)]
pub fn anything(item: TokenStream) -> TokenStream {
    let ast = match syn::parse::<DeriveInput>(item) {
        Ok(input) => input,
        Err(err) => return TokenStream::from(err.to_compile_error()),
    };
    let name = ast.ident;
    let add_hello_world = quote! {
        impl #name {
            fn hello_world(&self) {
                println!("Hello World")
            }
        }
    };
    add_hello_world.into()
}

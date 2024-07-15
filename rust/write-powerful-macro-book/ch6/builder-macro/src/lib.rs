use proc_macro::TokenStream;

#[proc_macro_derive(Builder)]
pub fn builder(item: TokenStream) -> TokenStream {
    builder_code::create_builder(item.into()).into()
}

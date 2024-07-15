use proc_macro::TokenStream;
use proc_macro2::{Span, TokenStream as TokenStream2};
use quote::{quote, ToTokens};
use syn::{
    parse::{Parse, ParseStream},
    parse_macro_input,
    punctuated::Punctuated,
    Data, DataStruct, DeriveInput, Fields, FieldsNamed, Ident, Token,
};

fn generated_methods(ast: &DeriveInput) -> Vec<TokenStream2> {
    let named_fields = match ast.data {
        Data::Struct(DataStruct {
            fields: Fields::Named(FieldsNamed { ref named, .. }),
            ..
        }) => named,
        _ => unimplemented!("only workds for structs with named fields"),
    };

    named_fields
        .iter()
        .map(|f| {
            let field_name = f.ident.as_ref().unwrap();
            let type_name = &f.ty;
            let method_name = Ident::new(&format!("get_{field_name}"), Span::call_site());

            quote!(
                fn #method_name(&self) -> &#type_name {
                    &self.#field_name
                }
            )
        })
        .collect()
}

#[proc_macro]
pub fn private(item: TokenStream) -> TokenStream {
    let item_as_stream: TokenStream2 = item.clone().into();

    let ast = parse_macro_input!(item as DeriveInput);
    let name = &ast.ident;
    let methods = generated_methods(&ast);

    quote! {
        #item_as_stream

        impl #name {
            #(#methods)*
        }
    }
    .into()
}

struct ComposeInput {
    expressions: Punctuated<Ident, Token![.]>,
}

impl Parse for ComposeInput {
    fn parse(input: ParseStream) -> syn::Result<Self> {
        Ok(ComposeInput {
            expressions: Punctuated::<Ident, Token![.]>::parse_terminated(input).unwrap(),
        })
    }
}

impl ToTokens for ComposeInput {
    fn to_tokens(&self, tokens: &mut TokenStream2) {
        let mut total = None;
        let mut as_idents: Vec<&Ident> = self.expressions.iter().collect();
        let last_ident = as_idents.pop().unwrap();

        as_idents.iter().rev().for_each(|i| {
            if let Some(current_total) = &total {
                total = Some(quote!(
                    compose_two(#i, #current_total)
                ));
            } else {
                total = Some(quote!(
                    compose_two(#i, #last_ident)
                ));
            }
        });
        total.to_tokens(tokens)
    }
}

#[proc_macro]
pub fn compose(item: TokenStream) -> TokenStream {
    let ci: ComposeInput = parse_macro_input!(item);

    quote!(
        {
            fn compose_two<FIRST, SECOND, THIRD, F, G>(first: F, second: G) -> impl Fn(FIRST) -> THIRD
                where
                    F: Fn(FIRST) -> SECOND,
                    G: Fn(SECOND) -> THIRD,
            {
                move |x| second(first(x))
            }
            #ci
        }
    ).into()
}

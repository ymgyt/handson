use wasm_bindgen::prelude::*;

#[wasm_bindgen]
extern "C" {
    pub fn alert(s: &str);
}

#[wasm_bindgen]
pub fn say_hello(name: &str, whom: &str) {
    alert(&format!("Hello, {} from {}", name, whom));
}

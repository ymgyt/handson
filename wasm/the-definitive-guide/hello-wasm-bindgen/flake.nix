{
  description = "wasm-pack setup";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    rust-overlay.url = "github:oxalica/rust-overlay";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { nixpkgs, rust-overlay, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system: 
    let
      pkgs = import nixpkgs {
        inherit system;
        overlays = [ rust-overlay.overlays.default ];
      };
    in
    {
      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          cargo
          wasm-pack
          (rust-bin.stable.latest.default.override {
            extensions = [ "rust-src" ];
            targets = [ "wasm32-unknown-unknown" ];
          })
        ];
      };
    });
}
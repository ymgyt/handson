{
  description = "OS 1000 lines flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/release-23.05";
    nixpkgs-unstable.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { self
    , nixpkgs
    , nixpkgs-unstable
    , flake-utils
    }:

    flake-utils.lib.eachDefaultSystem (system:
    let
      rscvPkgs = import nixpkgs-unstable {
        localSystem = "${system}";
        crossSystem = {
          config = "riscv32-unknown-linux-gnu";
        };
      };
      pkgs = import nixpkgs { inherit system; };
    in
    {
      devShells.default = pkgs.mkShell {
        # Tryied
        # llvmPackages_16.libcxxClang
        packages = with pkgs; [
          rscvPkgs.buildPackages.clang
        ];
        shellHook = ''
        '';
      };
    });
}

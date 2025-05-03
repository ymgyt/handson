{
  description = "development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.11";
    nixpkgs-unstable.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      nixpkgs-unstable,
      flake-utils,
    }:

    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
        pkgsUnstable = import nixpkgs-unstable { inherit system; };

        devPackages =
          with pkgs;
          [
            cargo-audit
            cargo-release
            cargo-machete
            fd
            git-cliff
            graphviz
            jsonnet
            ponysay
            renovate
            shfmt
            taplo
            typos
          ]
          ++ (with pkgsUnstable; [
            cargo-insta
            just
            terraform
            yamlfmt
          ])
          ++ pkgs.lib.optionals pkgs.stdenv.isDarwin [ ];
      in
      {
        devShells.default = pkgs.mkShell {
          packages = devPackages;
          shellHook = ''
            # 共通の環境変数をセットする
          '';
        };
      }
    );
}

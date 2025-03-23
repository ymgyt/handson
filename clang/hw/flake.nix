{
  description = "HelloWorld";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        packages.hello = pkgs.stdenv.mkDerivation {
          pname = "hello";
          version = "1.0";
          src = ./.;
          buildInputs = [ pkgs.gcc ];
          unpackPhase = ''
            cp  $src/hello.c .
          '';
          buildPhase = ''
            gcc hello.c -o hello
          '';
          installPhase = ''
            mkdir -p $out/bin
            cp hello $out/bin/
          '';
          meta = {
            description = "A simple HelloWorld C program";
          };
        };

        # アプリ定義（nix run で実行可能にする）
        apps.${system}.hello = {
          type = "app";
          program = "${self.packages.hello}/bin/hello";
        };
      }
    );
}

{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:

    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {

      devShells."${system}".default = pkgs.mkShell {
        packages = with pkgs; [
          minicom
          libusb1
          SDL2
          cargo-generate
          cargo-hf2
        ];
        shellHook = ''
          exec nu
        '';
      };
    };
}

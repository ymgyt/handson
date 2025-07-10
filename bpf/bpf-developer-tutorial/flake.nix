{
  description = "BPF development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # BPF toolchain
            # clang
            llvmPackages_20.clang-unwrapped
            llvm
            libbpf
            linuxHeaders

          ];

          shellHook = ''
            echo "BPF development environment loaded"
            echo "Available tools:"
            echo "  - clang ($(clang --version | head -1))"
            echo "  - libbpf"
            echo "  - kernel headers"
            echo "linux headers - ${pkgs.linuxHeaders}/include"
            echo "bpf   headers - ${pkgs.libbpf}/include"

            # Set header paths for BPF compilation
            export C_INCLUDE_PATH="${pkgs.linuxHeaders}/include:${pkgs.libbpf}/include"
          '';
        };
      }
    );
}

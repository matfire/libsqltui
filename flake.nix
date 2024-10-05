{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};
      in {
        devShell = pkgs.mkShell {
          buildInputs = [
            pkgs.go
          ];
        };
        packages.default = pkgs.buildGoModule {
          pname = "libsqltui";
          version = "0.1.0";
          src = ./.;
          module = "github.com/matfire/libsqltui";
          vendorHash = "sha256-cVkroKJlU+s9dIRuNSbKAk0evpwYTSoG6ZvtQzdRUaE=";
        };
 
      });
}

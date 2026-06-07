{
  description = "Severity-aware validation for Go with multiple outcome levels";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
    systems.url = "github:nix-systems/default";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{
      self,
      nixpkgs,
      flake-parts,
      systems,
      treefmt-nix,
    }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import systems;

      imports = [
        treefmt-nix.flakeModule
      ];

      perSystem =
        {
          config,
          pkgs,
          ...
        }:
        {
          treefmt = {
            projectRootFile = "go.mod";
            programs = {
              gofumpt.enable = true;
              nixfmt.enable = true;
            };
          };

          devShells = {
            default = pkgs.mkShell {
              name = "businessrules-dev";

              packages = [
                pkgs.go
                pkgs.golangci-lint
                pkgs.gopls
                pkgs.delve
                pkgs.gosec
                pkgs.gotools
                pkgs.gofumpt
              ];

              GOWORK = "off";

              shellHook = ''
                echo "businessrules dev shell"
                echo "Go: $(go version)"
              '';
            };

            ci = pkgs.mkShellNoCC {
              packages = [
                pkgs.go
                pkgs.golangci-lint
              ];

              GOWORK = "off";
            };
          };

          checks.fmt = pkgs.runCommand "businessrules-fmt-check" { nativeBuildInputs = [ pkgs.go ]; } ''
            cd ${builtins.path { path = ./.; name = "businessrules"; }}
            test -z "$(gofmt -l .)" || (echo "Files need formatting:"; gofmt -l .; exit 1)
            touch $out
          '';
        };
    };
}

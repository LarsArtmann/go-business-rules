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
      flake-parts,
      systems,
      treefmt-nix,
      ...
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
              goimports.enable = true;
              nixfmt.enable = true;
            };
          };

          checks.format = config.treefmt.build.check self;
          devShells = {
            default = pkgs.mkShell {
              name = "businessrules-dev";

              packages = [
                pkgs.go_1_26
                pkgs.golangci-lint
                pkgs.gopls
                pkgs.delve
                pkgs.gosec
                pkgs.gotools
                pkgs.gofumpt
              ];

              GOWORK = "off";
              GOEXPERIMENT = "jsonv2";
              GOPRIVATE = "github.com/larsartmann/*,github.com/LarsArtmann/*";
              GONOSUMDB = "github.com/larsartmann/*,github.com/LarsArtmann/*";
              GONOPROXY = "github.com/larsartmann/*,github.com/LarsArtmann/*";

              shellHook = ''
                echo "businessrules dev shell"
                echo "Go: $(go version)"
              '';
            };

            ci = pkgs.mkShellNoCC {
              packages = [
                pkgs.go_1_26
                pkgs.golangci-lint
              ];

              GOWORK = "off";
              GOEXPERIMENT = "jsonv2";
              GOPRIVATE = "github.com/larsartmann/*,github.com/LarsArtmann/*";
              GONOSUMDB = "github.com/larsartmann/*,github.com/LarsArtmann/*";
              GONOPROXY = "github.com/larsartmann/*,github.com/LarsArtmann/*";
            };
          };
        };
    };
}

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
        let
          goEnv = {
            GOEXPERIMENT = "jsonv2";
            GOWORK = "off";
            GOPRIVATE = "github.com/larsartmann/*,github.com/LarsArtmann/*";
            GONOSUMDB = "github.com/larsartmann/*,github.com/LarsArtmann/*";
            GONOPROXY = "github.com/larsartmann/*,github.com/LarsArtmann/*";
          };
        in
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

          # Build + vet + test + lint every Go module in this repository,
          # including the nested ones that the root commands do not cover.
          packages.check-all = pkgs.writeShellApplication {
            name = "go-business-rules-check-all";

            runtimeInputs = [
              pkgs.go_1_26
              pkgs.golangci-lint
            ];

            text = ''
              set -o errexit
              set -o pipefail

              export GOEXPERIMENT="${goEnv.GOEXPERIMENT}"
              export GOWORK="${goEnv.GOWORK}"
              export GOPRIVATE="${goEnv.GOPRIVATE}"
              export GONOSUMDB="${goEnv.GONOSUMDB}"
              export GONOPROXY="${goEnv.GONOPROXY}"

              root="$PWD"

              for module in . adapters/cqrslite examples/sse listeners/otel; do
                echo "=== module: $module ==="
                cd "$root/$module"
                go build ./...
                go vet ./...
                go test ./...
                golangci-lint run --timeout 5m
              done

              echo "ALL MODULES GREEN"
            '';

            meta = {
              description = "Build, vet, test, and lint every Go module in this repository";
              license = pkgs.lib.licenses.mit;
              mainProgram = "go-business-rules-check-all";
            };
          };

          apps.check-all = {
            type = "app";
            program = pkgs.lib.getExe config.packages.check-all;
          };

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
                config.packages.check-all
              ];

              inherit (goEnv)
                GOWORK
                GOEXPERIMENT
                GOPRIVATE
                GONOSUMDB
                GONOPROXY
                ;

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

              inherit (goEnv)
                GOWORK
                GOEXPERIMENT
                GOPRIVATE
                GONOSUMDB
                GONOPROXY
                ;
            };
          };
        };
    };
}

{
  description = "Severity-aware validation for Go with multiple outcome levels";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{
      self,
      flake-parts,
      treefmt-nix,
      ...
    }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      # Inline systems: nixpkgs 26.11 dropped x86_64-darwin, which github:nix-systems/default still lists.
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "aarch64-darwin"
      ];

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
          # The go.mod floors are go 1.27, so every Go here must be >= 1.27.1
          # (a go-cqrs-lite dependency requires it): goimports shells out to the
          # `go` binary and would otherwise try to download a toolchain, which
          # dies in the hermetic treefmt sandbox.
          go127Gotools = pkgs.gotools.override {
            buildGoModule = pkgs.buildGo127Module;
            go = pkgs.go_1_27;
          };
          # Mirrors /home/lars/projects/SystemNix/pkgs/govalid.nix; keep the
          # rev in sync when bumping there. Needed in the devShell so BuildFlow
          # runs govalid with this shell's go (>= 1.27.1) instead of the
          # system binary's older bundled environment.
          govalid = pkgs.buildGo127Module {
            pname = "govalid";
            version = "0-unstable-2026-09-17";

            src = pkgs.fetchFromGitHub {
              owner = "sivchari";
              repo = "govalid";
              rev = "8d6700c031967fa871a0e1739f507ab2e19f4615";
              hash = "sha256-yA2lMdy6HKgPkd0+yqNWJdAC7Jxwtmsgif6s2Q6LDRM=";
            };

            subPackages = [ "cmd/govalid" ];

            doCheck = false;

            vendorHash = "sha256-fKvE4wGU8PQbzgxTnUaRNqbTy6JlzDMBWcWGy9uUTqo=";

            meta = {
              description = "Type-safe struct validation code generator for Go";
              homepage = "https://github.com/sivchari/govalid";
              license = pkgs.lib.licenses.mit;
              maintainers = [
                {
                  name = "Lars Artmann";
                  github = "LarsArtmann";
                }
              ];
              platforms = pkgs.lib.platforms.all;
              mainProgram = "govalid";
            };
          };
        in
        {
          treefmt = {
            projectRootFile = "go.mod";
            programs = {
              gofumpt.enable = true;
              goimports = {
                enable = true;
                package = go127Gotools;
              };
              nixfmt.enable = true;
            };
          };

          checks.format = config.treefmt.build.check self;

          # Build + vet + test + lint every Go module in this repository,
          # including the nested ones that the root commands do not cover.
          packages.check-all = pkgs.writeShellApplication {
            name = "go-business-rules-check-all";

            runtimeInputs = [
              pkgs.go_1_27
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
              homepage = "https://github.com/LarsArtmann/go-business-rules";
              license = pkgs.lib.licenses.mit;
              mainProgram = "go-business-rules-check-all";
              maintainers = [
                {
                  name = "Lars Artmann";
                  github = "LarsArtmann";
                }
              ];
              platforms = pkgs.lib.platforms.unix;
            };
          };

          apps.check-all = {
            type = "app";
            program = pkgs.lib.getExe config.packages.check-all;
            meta.description = "Build, vet, test, and lint every Go module in this repository";
          };

          devShells = {
            default = pkgs.mkShell {
              name = "businessrules-dev";

              packages = [
                pkgs.go_1_27
                pkgs.golangci-lint
                pkgs.go-licenses
                govalid
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
                pkgs.go_1_27
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

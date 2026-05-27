{
  description = "Severity-aware validation for Go with multiple outcome levels";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      pkgsFor = system: nixpkgs.legacyPackages.${system};
    in
    {
      devShells = forAllSystems (
        system:
        let
          pkgs = pkgsFor system;
        in
        {
          default = pkgs.mkShell {
            name = "businessrules-dev";

            packages = [
              pkgs.go
              pkgs.golangci-lint
              pkgs.gopls
              pkgs.delve
              pkgs.just
              pkgs.gosec
              pkgs.gotools
              pkgs.gofumpt
              pkgs.nixfmt
            ];

            GOWORK = "off";

            shellHook = ''
              echo "businessrules dev shell"
              echo "Go: $(go version)"
              echo ""
              echo "  just test        Run tests with race detection"
              echo "  just lint        Run golangci-lint"
              echo "  just check       Run all quality checks"
              echo "  nix flake check  Run hermetic checks"
            '';
          };
        }
      );

      checks = forAllSystems (
        system:
        let
          pkgs = pkgsFor system;
          src = builtins.path {
            path = ./.;
            name = "businessrules";
          };
        in
        {
          fmt = pkgs.runCommand "businessrules-fmt-check" { nativeBuildInputs = [ pkgs.go ]; } ''
            cd ${src}
            test -z "$(gofmt -l .)" || (echo "Files need formatting:"; gofmt -l .; exit 1)
            touch $out
          '';
        }
      );

      formatter = forAllSystems (system: (pkgsFor system).nixfmt);
    };
}

{
  description = "Generate QR codes in your terminal — one standalone binary";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
      # Keep in sync with the git tag when releasing (see the release runbook).
      version = "0.1.0";
    in
    {
      packages = forAllSystems (pkgs: rec {
        qrcli = pkgs.buildGoModule {
          pname = "qrcli";
          inherit version;
          src = self;
          # Re-compute whenever go.mod/go.sum change:
          # set to nixpkgs.lib.fakeHash, build, copy the "got:" hash.
          vendorHash = "sha256-xGPzODAJOls8RyyYdoEbPqz63i4oex41QsF1HNpmWAc=";
          env.CGO_ENABLED = 0;
          ldflags = [ "-s" "-w" "-X main.version=${version}" ];
          # buildGoModule runs `go test ./...` by default (doCheck = true),
          # so `nix build` doubles as a test gate.
          meta = with pkgs.lib; {
            description = "Generate QR codes in your terminal — one standalone binary";
            homepage = "https://github.com/CaporalDead/qrcli";
            license = licenses.mit;
            mainProgram = "qrcli";
          };
        };
        default = qrcli;
      });

      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            go
            golangci-lint
            goreleaser
            gnumake
          ];
        };
      });
    };
}

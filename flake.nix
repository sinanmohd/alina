{
  inputs.nixpkgs.url = "github:NixOs/nixpkgs/nixos-unstable";

  outputs =
    inputs@{ self, nixpkgs }:
    let
      lib = nixpkgs.lib;

      forSystem =
        f: system:
        f {
          inherit system;
          pkgs = import nixpkgs { inherit system; };
        };
      supportedSystems = lib.platforms.unix;
      forAllSystems = f: lib.genAttrs supportedSystems (forSystem f);

      version =
        if self ? shortRev then
          self.shortRev
        else if self ? dirtyShortRev then
          self.dirtyShortRev
        else
          "not-a-gitrepo";
    in
    {
      packages = forAllSystems (
        { system, pkgs }:
        {
          frontend = pkgs.callPackage ./nix/packages/frontend.nix {
            inherit version;
          };

          alina = pkgs.callPackage ./nix/packages/alina.nix {
            inherit version;
            frontend = self.packages.${system}.frontend;
          };
          default = self.packages.${system}.alina;
        }
      );

      nixosModules = {
        alina = import ./nix/module.nix inputs;
        default = self.nixosModules.alina;
      };

      devShells = forAllSystems (
        { system, pkgs }:
        {
          alina = pkgs.callPackage ./nix/shell.nix {
            frontend = self.packages.${system}.frontend;
            alina = self.packages.${system}.alina;
          };
          default = self.devShells.${system}.alina;
        }
      );
    };
}

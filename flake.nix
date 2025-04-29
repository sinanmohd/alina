{
  inputs.nixpkgs.url = "github:NixOs/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }: let
    lib = nixpkgs.lib;

    forSystem = f: system: f {
      inherit system;
      pkgs = import nixpkgs { inherit system; };
    };

    supportedSystems = lib.platforms.unix;
    forAllSystems = f: lib.genAttrs supportedSystems (forSystem f);
  in {
    devShells = forAllSystems ({ system, pkgs, ... }: {
      default = pkgs.mkShell {
        name = "dev";

        buildInputs = with pkgs; [
          go
          gopls
          sqlc

          nodejs
          vue-language-server
          nixfmt-rfc-style
          tailwindcss-language-server
        ];
	shellHook = ''
          export PS1="\033[0;31m[alina]\033[0m $PS1"
          export npm_config_nodedir=${pkgs.nodePackages.nodejs}
        '';
      };
    });
  };
}

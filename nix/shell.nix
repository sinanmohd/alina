{
  mkShell,
  frontend,
  alina,
  nixfmt-rfc-style,
  typescript-language-server,
  vue-language-server,
  tailwindcss-language-server,
  gopls,
  sqlc,
  air,
  nodePackages,
}:

mkShell {
  inputsFrom = [
    frontend
    alina
  ];

  buildInputs = [
    gopls
    sqlc
    air

    nixfmt-rfc-style

    vue-language-server
    typescript-language-server
    tailwindcss-language-server
  ];

  shellHook = ''
    export PS1="\033[0;31m[alina]\033[0m $PS1"
    export NODE_ENV="development"

    export npm_config_nodedir=${nodePackages.nodejs}
    pnpm config set store-dir ~/.local/share/pnpm
  '';
}

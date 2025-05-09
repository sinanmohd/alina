{
  lib,
  stdenv,
  nodejs,
  pnpm,
  version,
}:

stdenv.mkDerivation (finalAttrs: {
  inherit version;
  pname = "alina-frontend";

  src = lib.cleanSourceWith {
    filter =
      name: type:
      lib.cleanSourceFilter name type
      && !(builtins.elem (baseNameOf name) [
        "nix"
        "flake.nix"
      ]);

    src = ../../frontend;
  };

  buildPhase = ''
    export NODE_OPTIONS="--max_old_space_size=16384"
    export NUXT_TELEMETRY_DISABLED=1
    export npm_config_nodedir=${nodejs}

    pnpm exec nuxt generate
  '';

  installPhase = ''
    mkdir -p $out/share
    cp -r .output/public $out/share/www
  '';

  nativeBuildInputs = [
    pnpm
    pnpm.configHook
    nodejs
  ];

  pnpmDeps = pnpm.fetchDeps {
    inherit (finalAttrs) pname version src;
    hash = "sha256-z29qhXy7qsrUXiOz+AzBK3MJ3/FAbH/1PUXZdvWJXa0=";
  };

  meta = {
    description = "Your frenly neighbourhood file sharing website.";
    homepage = "https://github.com/sinanmohd/alina";
    platforms = lib.platforms.unix;
    license = lib.licenses.agpl3Plus;
    maintainers = with lib.maintainers; [ sinanmohd ];
  };
})

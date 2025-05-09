{
  lib,
  buildGoModule,
  version,
  frontend,
}:

buildGoModule (finalAttrs: {
  inherit version;
  pname = "alina";

  src = lib.cleanSourceWith {
    filter =
      name: type:
      lib.cleanSourceFilter name type
      && !(builtins.elem (baseNameOf name) [
        "nix"
        "flake.nix"
      ]);

    src = ../../backend;
  };

  nativeBuildInputs = [ frontend ];
  postConfigure = ''
    cp -rv ${frontend}/share/www internal/server/frontend
  '';

  vendorHash = "sha256-N/7WvfOp3CJnQYfAiHMpL+Y21PhOB1B110cGmLXqnMc=";

  meta = {
    description = "Your frenly neighbourhood file sharing website.";
    homepage = "https://github.com/sinanmohd/alina";
    platforms = lib.platforms.unix;
    license = lib.licenses.agpl3Plus;
    mainProgram = "alina";
    maintainers = with lib.maintainers; [ sinanmohd ];
  };
})

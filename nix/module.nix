inputs:
{
  config,
  lib,
  pkgs,
  ...
}:

let
  cfg = config.services.alina;
  inherit (pkgs.stdenv.hostPlatform) system;

  configFormat = pkgs.formats.toml { };
  configFile = configFormat.generate "alina.toml" cfg.settings;

  defaultEnvs = {
    ALINA_CONFIG = "${configFile}";
  };
in
{
  meta.maintainers = with lib.maintainers; [ sinanmohd ];

  options.services.alina = {
    enable = lib.mkEnableOption "alina";
    package = lib.mkOption {
      type = lib.types.package;
      description = "The alina package to use.";
      default = inputs.self.packages.${system}.alina;
    };

    port = lib.mkOption {
      type = lib.types.port;
      default = 8008;
      description = "The port alina should be reachable from.";
    };
    environment = lib.mkOption {
      default = { };
      type = lib.types.attrsOf lib.types.str;
    };
    environmentFile = lib.mkOption {
      type = lib.types.nullOr lib.types.path;
      example = "/var/lib/alina/secrets";
      default = null;
      description = ''
        Secrets may be passed to the service without adding them to the world-readable Nix store using this option.
      '';
    };
    settings = lib.mkOption {
      inherit (configFormat) type;
      default = { };
      description = ''
        Configuration options for alina.
      '';
    };
  };

  config = lib.mkIf cfg.enable {
    environment.systemPackages = [ cfg.package ];
    services.alina.settings.server.port = lib.mkDefault cfg.port;

    # This service stores a potentially large amount of data.
    # Running it as a dynamic user would force chown to be run everytime the
    # service is restarted on a potentially large number of files.
    # That would cause unnecessary and unwanted delays.
    users = {
      groups.alina = { };
      users.alina = {
        isSystemUser = true;
        group = "alina";
      };
    };

    systemd.services.alina = {
      description = "Your frenly neighbourhood file sharing website.";
      wantedBy = [ "multi-user.target" ];
      after = [ "network-online.target" ];
      environment = defaultEnvs // cfg.environment;

      serviceConfig = {
        Type = "simple";
        StateDirectory = "alina";
        Restart = "on-failure";
        EnvironmentFile = lib.mkIf (cfg.environmentFile != null) cfg.environmentFile;
        ExecStart = lib.getExe cfg.package;
      };
    };
  };
}

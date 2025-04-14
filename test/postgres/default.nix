{ pkgs, ... }:

{
  config = {
    networking.firewall.allowedTCPPorts = [ 5432 ];

    services.postgresql = {
      enable = true;
      enableTCPIP = true;
      authentication = pkgs.lib.mkOverride 10 ''
        #type database DBuser origin-address auth-method
        local all      all     trust
        # ipv4
        host  all      all     0.0.0.0/0   trust
        # ipv6
        host  all      all     ::0/0        trust
      '';
      initialScript = ./setup.sql;
    };

    users = {
      groups.database = {};
    };
  };
}

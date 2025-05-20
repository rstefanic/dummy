{ pkgs, ... }:

{
  environment.defaultPackages = [
    (pkgs.callPackage ./.. { })
  ];

  environment.etc."dummy.yml".text = ''
    server:
      host: "postgres"
      name: "postgres"
      user: "postgres"
  '';
}

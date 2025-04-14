{ pkgs, ... }:

{
  environment.defaultPackages = [
    (pkgs.callPackage ./.. { })
  ];
}

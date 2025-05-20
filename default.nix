{
  pkgs ? import <nixpkgs> {}
}:

pkgs.buildGoModule {
  pname = "dummy";
  version = "0.1.0";

  # In 'nix develop', we don't need a copy of the source tree
  # in the Nix store.
  src = ./.;
  vendorHash = "sha256-6bbLkCnQy/LWUqMGktHxZ6JVGoWVJ0CtCyKONQZk9q0=";
}

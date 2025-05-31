{
  pkgs ? import <nixpkgs> {}
}:

pkgs.buildGoModule {
  pname = "dummy";
  version = "0.1.0";

  # In 'nix develop', we don't need a copy of the source tree
  # in the Nix store.
  src = ./.;
  vendorHash = "sha256-Y+E58mWniwSjDQOV4gF7LqTm/4GGNygsw7Ils+atAsk=";
}

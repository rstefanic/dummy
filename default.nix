{
  pkgs ? import <nixpkgs> {}
}:

pkgs.buildGoModule {
  pname = "dummy";
  version = "0.1.0";

  # In 'nix develop', we don't need a copy of the source tree
  # in the Nix store.
  src = ./.;
  vendorHash = "sha256-W/aH6ax3G5HGyL8MZjLFZ53d8E1mRI1smyRNhLcF9wc=";
}

{ pkgs ? import (builtins.fetchTarball {
    url = "https://github.com/nixos/nixpkgs/archive/ccd458053b0e.tar.gz";
    sha256 = "1qwnmmb2p7mj1h1ffz1wvkr1v55qbhzvxr79i3a15blq622r4al9";
  }) { } 
}:

pkgs.buildGoModule {
  pname = "neuron-language-server";
  version = "0.1.1";

  src = ./.;

  vendorSha256 = "0pjjkw0633l8qbvwzy57rx76zjn3w3kf5f7plxnpxih9zj0q258l";
}

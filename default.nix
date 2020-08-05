{ pkgs ? import (builtins.fetchTarball {
    url = "https://github.com/nixos/nixpkgs/archive/ccd458053b0e.tar.gz";
    sha256 = "1qwnmmb2p7mj1h1ffz1wvkr1v55qbhzvxr79i3a15blq622r4al9";
  }) { } 
}:

pkgs.buildGoModule {
  pname = "neuron-language-server";
  version = "0.1";

  src = ./.;

  vendorSha256 = "02dajl4l3c8522ik2hmiq8cx4kj4h2ykx8l7qsal5xznx9pqbs7i";
}

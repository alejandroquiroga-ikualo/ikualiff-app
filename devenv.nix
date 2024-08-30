{ pkgs, lib, config, inputs, ... }:

let
  pkgs-unstable = import inputs.nixpkgs-unstable { system = pkgs.stdenv.system; };
in {
  languages.go = {
    enable = true;
    package = pkgs-unstable.go;
  };

  packages = [
    pkgs.postgresql
  ];
}


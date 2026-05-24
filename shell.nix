{
  pkgs ? import <nixpkgs> { },
}:
pkgs.mkShellNoCC {
  buildInputs = with pkgs; [
    go
    gopls
    golangci-lint
  ];
}

{ pkgs, lib, config, inputs, ... }:

{
  # https://devenv.sh/basics/
  env.GREET = "devenv";
  env.POSTGRESQL_URL = "postgres://vhodges:@localhost:5432/kuamua_development?sslmode=disable&search_path=public";

  # https://devenv.sh/packages/
  packages = [ 
    pkgs.git 
    pkgs.sqlc
    pkgs.bruno
    pkgs.k6
  ];

  # https://devenv.sh/languages/
  languages.go.enable = true;

  # https://devenv.sh/processes/
  # processes.cargo-watch.exec = "cargo-watch";

  # https://devenv.sh/services/
  services.postgres.enable = true;
  services.postgres.initialDatabases = [
    { name = "kuamua_development"; }
  ];
  services.postgres.listen_addresses = "127.0.0.1";

  # https://devenv.sh/scripts/
  scripts.hello.exec = ''
    echo hello from $GREET
  '';

  enterShell = ''
    hello
    git --version
    go version
  '';

  # https://devenv.sh/tasks/
  # tasks = {
  #   "myproj:setup".exec = "mytool build";
  #   "devenv:enterShell".after = [ "myproj:setup" ];
  # };

  # https://devenv.sh/tests/
  enterTest = ''
    echo "Running tests"
    git --version | grep --color=auto "${pkgs.git.version}"
  '';

  # https://devenv.sh/pre-commit-hooks/
  # pre-commit.hooks.shellcheck.enable = true;

  # See full reference at https://devenv.sh/reference/options/
}

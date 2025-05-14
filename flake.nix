{
  description = "dummy";

  inputs.nixpkgs.url = "nixpkgs/nixos-24.11";

  outputs = { self, nixpkgs }:
    let
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          inherit (self.checks.${system}.default) driverInteractive;
          dummy = pkgs.callPackage ./. { };
        });

      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [ go gopls gotools go-tools ];
          };
        });

      checks = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
          expected = pkgs.lib.strings.concatStrings (
            pkgs.lib.strings.intersperse "\n" [
              "-- host: postgres"
              "-- name: postgres"
              "-- user: postgres"
              "-- pass: "
              "-- table: todos"
              "-- seed: 1"
              ""
              "INSERT INTO todos (id,task,complete,created_at) VALUES (DEFAULT,'Change itself still I world without that myself how below.',false,'%d-%11d-%11d'),(DEFAULT,'How bored light what nearby regularly children constantly judge turkey.',false,'%d-%07d-%07d'),(DEFAULT,'Year never you hand her hand certain aircraft village skip.',false,'%d-%17d-%17d');\n"
            ]
          );
        in
        {
          default = pkgs.nixosTest {
            name = "Integration Test";

            nodes = {
              postgres = import ./test/postgres;
              client = import ./test/client.nix;
            };

            testScript = ''
              start_all()
              postgres.wait_for_open_port(5432)

              expected = """${expected}"""
              actual = client.succeed('dummy --host postgres --name postgres --user postgres --table todos --count 3 --seed 1')

              assert expected == actual, "actual:\n" + actual + "\nexpected:\n" + expected
            '';
          };
        }
      );

      defaultPackage = forAllSystems (system: self.packages.${system}.dummy);
    };
}

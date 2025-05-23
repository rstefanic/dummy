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
              "-- seed: 1"
              ""
              "INSERT INTO todos (id,task,complete,created_at) VALUES (DEFAULT,'Change.',false,'2015-04-29'),(DEFAULT,'Without.',true,'1954-01-27'),(DEFAULT,'Tomorrow.',true,'1917-05-22');"
              ""
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
              actual = client.succeed('dummy --path /etc/dummy.yml')

              assert expected == actual, "actual:\n" + actual + "\nexpected:\n" + expected
            '';
          };
        }
      );
    };
}

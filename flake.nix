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
              "host:  postgres"
              "name:  postgres"
              "user:  postgres"
              "pass:  "
              "table:  todos"
              "seed:  1"
              "INSERT INTO todos (id,task,complete,created_at) VALUES (9172393864939720632,'Their turn have her formerly honour while Bismarckian example product.',false,'%d-%29d-%29d'),(424577099883569842,'My government not regularly children constantly judge turkey whoa what.',false,'%d-%25d-%25d'),(5560765190427624264,'Later include hand certain aircraft village skip his he yesterday.',true,'%d-%06d-%06d');\n"
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

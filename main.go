package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/goccy/go-yaml"
	_ "github.com/lib/pq"

	"dummy/postgresql"
)

func main() {
	var (
		path      string
		tableName string
		count     int
		seed      int
	)

	flag.StringVar(&path, "path", "dummy.yml", "Path to the configuration yaml file.")
	flag.StringVar(&tableName, "table", "", "The table that you want to create data dummy for.")
	flag.IntVar(&count, "count", 0, "The number of rows of dummy data to generate.")
	flag.IntVar(&seed, "seed", rand.Int(), "Set the seeder used to generate the output.")
	flag.Parse()

	configFile, err := os.ReadFile(path)
	if err != nil {
		panic("could not find file specified")
	}

	var config Config
	yaml.Unmarshal(configFile, &config)

	// If there was no seed set in the config file, randomize it
	if config.Options.Seed == 0 {
		config.Options.Seed = seed
	}

	if tableName == "" && len(config.Tables) == 0 {
		panic("argument \"table\" is required")
	}

	if tableName == "" {
		tableName = config.Tables[0].Name
	}

	// Try to use the one from the config if we're missing the `count` program argument
	if count == 0 {
		count = config.Tables[0].Count
	}

	if count == 0 {
		panic("\"count\" is 0 -- no data to generate")
	}

	var table = postgresql.NewTable(tableName)

	gofakeit.Seed(config.Options.Seed)

	psqlDb, err := postgresql.New(config.Server.User, config.Server.Password, config.Server.Host, config.Server.Name)
	if err != nil {
		panic("(postgresql.New): " + err.Error())
	}
	defer psqlDb.Close()

	columns, err := psqlDb.GetTableColumns(tableName)
	if err != nil {
		panic("(postgresql.GetTableData): " + err.Error())
	}
	table.Columns = columns

	err = table.FillMetadata()
	if err != nil {
		panic(err)
	}

	err = table.CreateData(count)
	if err != nil {
		panic(err)
	}

	if !config.Options.HideInputComment {
		fmt.Println("-- host:", config.Server.Host)
		fmt.Println("-- name:", config.Server.Name)
		fmt.Println("-- user:", config.Server.User)
		fmt.Println("-- pass:", config.Server.Password)
		fmt.Println("-- table:", tableName)
		fmt.Println("-- seed:", config.Options.Seed)
		fmt.Println("")
	}

	fmt.Println(table.ToPsqlStatement())
}

type Config struct {
	Server struct {
		Host     string `yaml:"host"`
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	}
	Options struct {
		Seed             int  `yaml:"seed"`
		HideInputComment bool `yaml:"hideInputComments"`
	}
	Tables []struct {
		Name  string `yaml:"name"`
		Count int    `yaml:"count"`
	}
}

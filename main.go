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
		path string
		seed int
	)

	flag.StringVar(&path, "path", "dummy.yml", "Path to the configuration yaml file.")
	flag.IntVar(&seed, "seed", rand.Int(), "Set the seeder used to generate the output.")
	flag.Parse()

	configFile, err := os.ReadFile(path)
	if err != nil {
		panic("could not find file specified")
	}

	var config Config
	yaml.Unmarshal(configFile, &config)

	// If there was no seed set in the config file, use the randomized one
	if config.Options.Seed == 0 {
		config.Options.Seed = seed
	}
	gofakeit.Seed(config.Options.Seed)

	if !config.Options.HideInputComment {
		fmt.Println("-- host:", config.Server.Host)
		fmt.Println("-- name:", config.Server.Name)
		fmt.Println("-- user:", config.Server.User)
		fmt.Println("-- pass:", config.Server.Password)
		fmt.Println("-- seed:", config.Options.Seed)
		fmt.Println("")
	}

	for i, table := range config.Tables {
		if i > 0 {
			fmt.Print("\n\n")
		}
		var t = postgresql.NewTable(table.Name)

		psqlDb, err := postgresql.New(config.Server.User, config.Server.Password, config.Server.Host, config.Server.Name)
		if err != nil {
			panic("(postgresql.New): " + err.Error())
		}
		defer psqlDb.Close()

		columns, err := psqlDb.GetTableColumns(t.Name)
		if err != nil {
			panic("(postgresql.GetTableData): " + err.Error())
		}
		t.Columns = columns

		err = t.FillMetadata()
		if err != nil {
			panic(err)
		}

		// Set the defualt count to 10 if it isn't set
		var count int
		if table.Count == 0 {
			count = 10
		} else {
			count = table.Count
		}

		err = t.CreateData(count)
		if err != nil {
			panic(err)
		}

		fmt.Println(t.ToPsqlStatement())
	}
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

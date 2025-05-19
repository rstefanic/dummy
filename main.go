package main

import (
	"flag"
	"fmt"
	"math/rand/v2"

	"github.com/brianvoe/gofakeit/v7"
	_ "github.com/lib/pq"

	"dummy/postgresql"
)

func main() {
	var (
		host              string
		name              string
		user              string
		password          string
		tableName         string
		count             int
		seed              int
		hideInputComments bool
	)

	flag.StringVar(&host, "host", "127.0.0.1", "The host to connect to.")
	flag.StringVar(&name, "name", "postgres", "Name of the database to connect to.")
	flag.StringVar(&user, "user", "root", "The user/role in the DB to connect with.")
	flag.StringVar(&password, "password", "", "The password of the user/role to connect with.")
	flag.StringVar(&tableName, "table", "", "The table that you want to create data dummy for.")
	flag.IntVar(&count, "count", 10, "The number of rows of dummy data to generate.")
	flag.IntVar(&seed, "seed", rand.Int(), "Set the seeder used to generate the output.")
	flag.BoolVar(&hideInputComments, "hide-input-comments", false, "Write the input information as a comment in the output.")
	flag.Parse()

	if tableName == "" {
		panic("argument \"table\" is required")
	}
	var table = postgresql.NewTable(tableName)

	gofakeit.Seed(seed)

	psqlDb, err := postgresql.New(user, password, host, name)
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

	if !hideInputComments {
		fmt.Println("-- host:", host)
		fmt.Println("-- name:", name)
		fmt.Println("-- user:", user)
		fmt.Println("-- pass:", password)
		fmt.Println("-- table:", tableName)
		fmt.Println("-- seed:", seed)
		fmt.Println("")
	}

	fmt.Println(table.ToPsqlStatement())
}

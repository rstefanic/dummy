package main

import (
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	var (
		host     string
		name     string
		user     string
		password string
		table    string
	)

	flag.StringVar(&host, "host", "127.0.0.1", "The host to connect to.")
	flag.StringVar(&name, "name", "postgres", "Name of the database to connect to.")
	flag.StringVar(&user, "user", "root", "The user/role in the DB to connect with.")
	flag.StringVar(&password, "password", "", "The password of the user/role to connect with.")
	flag.StringVar(&table, "table", "", "The table that you want to query for its schema.")
	flag.Parse()

	fmt.Println("host: ", host)
	fmt.Println("name: ", name)
	fmt.Println("user: ", user)
	fmt.Println("pass: ", password)
	fmt.Println("table: ", table)

	if table == "" {
		panic("argument \"table\" is required")
	}

	connString := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", user, password, host, name)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic("(sql.Open): " + err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic("(db.Query): " + err.Error())
	}

	fmt.Println("All is good! Closing connection...")
}

package main

import (
	"database/sql"
	"flag"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	_ "github.com/lib/pq"

	t "dummy/table"
)

func main() {
	var (
		host      string
		name      string
		user      string
		password  string
		tableName string
		count     int
		seed      int
	)

	flag.StringVar(&host, "host", "127.0.0.1", "The host to connect to.")
	flag.StringVar(&name, "name", "postgres", "Name of the database to connect to.")
	flag.StringVar(&user, "user", "root", "The user/role in the DB to connect with.")
	flag.StringVar(&password, "password", "", "The password of the user/role to connect with.")
	flag.StringVar(&tableName, "table", "", "The table that you want to create data dummy for.")
	flag.IntVar(&count, "count", 10, "The number of rows of dummy data to generate.")
	flag.IntVar(&seed, "seed", 0, "Set the seeder used to generate the output.")
	flag.Parse()

	fmt.Println("host: ", host)
	fmt.Println("name: ", name)
	fmt.Println("user: ", user)
	fmt.Println("pass: ", password)
	fmt.Println("table: ", tableName)
	fmt.Println("seed: ", seed)

	if tableName == "" {
		panic("argument \"table\" is required")
	}
	var table = t.New(tableName)

	gofakeit.Seed(seed)

	connString := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", user, password, host, name)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic("(sql.Open): " + err.Error())
	}
	defer db.Close()

	rows, err := db.Query(`SELECT
		column_name, ordinal_position, column_default, is_nullable, data_type,
		character_maximum_length, character_octet_length, numeric_precision,
		numeric_precision_radix, numeric_scale, datetime_precision, udt_name,
		is_self_referencing, is_identity, identity_generation, identity_start,
		identity_increment, identity_maximum, identity_minimum, is_updatable
		FROM information_schema.columns WHERE table_name = $1 ORDER BY ordinal_position`, tableName)

	if err != nil {
		panic("(db.Query): " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var c t.Column
		err := rows.Scan(&c.Name, &c.OrdinalPosition, &c.ColumnDefault, &c.IsNullable,
			&c.DataType, &c.CharacterMaximumLength, &c.CharacterOctetLength,
			&c.NumericPrecision, &c.NumericPrecisionRadix, &c.NumericScale,
			&c.DatetimePrecision, &c.UdtName, &c.IsSelfReferencing,
			&c.IsIdentity, &c.IdentityGeneration, &c.IdentityStart,
			&c.IdentityIncrement, &c.IdentityMaximum, &c.IdentityMinimum,
			&c.IsUpdateable)

		if err != nil {
			panic("(rows.Scan): " + err.Error())
		}

		table.Columns = append(table.Columns, c)
	}

	err = table.FillMetadata()
	if err != nil {
		panic(err)
	}

	err = table.CreateData(count)
	if err != nil {
		panic(err)
	}

	fmt.Println(table.ToPsqlStatement())
}

package main

import (
	"database/sql"
	"flag"
	"fmt"

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
	)

	flag.StringVar(&host, "host", "127.0.0.1", "The host to connect to.")
	flag.StringVar(&name, "name", "postgres", "Name of the database to connect to.")
	flag.StringVar(&user, "user", "root", "The user/role in the DB to connect with.")
	flag.StringVar(&password, "password", "", "The password of the user/role to connect with.")
	flag.StringVar(&tableName, "table", "", "The table that you want to create data dummy for.")
	flag.Parse()

	fmt.Println("host: ", host)
	fmt.Println("name: ", name)
	fmt.Println("user: ", user)
	fmt.Println("pass: ", password)
	fmt.Println("table: ", tableName)

	if tableName == "" {
		panic("argument \"table\" is required")
	}

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
		FROM information_schema.columns WHERE table_name = $1`, tableName)

	if err != nil {
		panic("(db.Query): " + err.Error())
	}
	defer rows.Close()

	var table = t.New(tableName)

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

	err = table.CreateData(10)
	if err != nil {
		panic(err)
	}

	fmt.Println(table.ToPsqlStatement())
}

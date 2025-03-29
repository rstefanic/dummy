package main

import (
	"database/sql"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-faker/faker/v4"
	_ "github.com/lib/pq"
)

type ColumnInfo struct {
	name                   string
	ordinalPosition        int
	columnDefault          sql.NullString
	isNullable             string
	dataType               string
	characterMaximumLength sql.NullInt32
	characterOctetLength   sql.NullInt32
	numericPrecision       sql.NullInt32
	numericPrecisionRadix  sql.NullInt32
	numericScale           sql.NullInt32
	datetimePrecision      sql.NullInt16
	udtName                string
	isSelfReferencing      string
	isIdentity             string
	identityGeneration     sql.NullString
	identityStart          sql.NullInt32
	identityIncrement      sql.NullInt32
	identityMaximum        sql.NullInt32
	identityMinimum        sql.NullInt32
	isUpdateable           string
}

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

	rows, err := db.Query(`SELECT
		column_name, ordinal_position, column_default, is_nullable, data_type,
		character_maximum_length, character_octet_length, numeric_precision,
		numeric_precision_radix, numeric_scale, datetime_precision, udt_name,
		is_self_referencing, is_identity, identity_generation, identity_start,
		identity_increment, identity_maximum, identity_minimum, is_updatable
		FROM information_schema.columns WHERE table_name = $1`, table)

	if err != nil {
		panic("(db.Query): " + err.Error())
	}
	defer rows.Close()

	var tableInfo []ColumnInfo

	for rows.Next() {
		var c ColumnInfo
		err := rows.Scan(&c.name, &c.ordinalPosition, &c.columnDefault, &c.isNullable,
			&c.dataType, &c.characterMaximumLength, &c.characterOctetLength,
			&c.numericPrecision, &c.numericPrecisionRadix, &c.numericScale,
			&c.datetimePrecision, &c.udtName, &c.isSelfReferencing,
			&c.isIdentity, &c.identityGeneration, &c.identityStart,
			&c.identityIncrement, &c.identityMaximum, &c.identityMinimum,
			&c.isUpdateable)

		if err != nil {
			panic("(rows.Scan): " + err.Error())
		}

		tableInfo = append(tableInfo, c)
	}

	fmt.Println(generatePsqlInsertStatement(tableInfo, table))
}

func generatePsqlInsertStatement(tableInfo []ColumnInfo, tableName string) string {
	var output strings.Builder
	output.WriteString("INSERT INTO ")
	output.WriteString(tableName)
	output.WriteString(" (")

	for i, col := range tableInfo {
		if i > 0 {
			output.WriteRune(',')
		}

		output.WriteString(col.name)
	}

	output.WriteString(") VALUES (")

	for i, col := range tableInfo {
		var value string

		if i > 0 {
			output.WriteRune(',')
		}

		switch col.udtName {
		case "int4":
			intVal, err := faker.RandomInt(0, 128, 1)
			value = strconv.Itoa(intVal[0])
			if err != nil {
				panic("(faker.RandomInt for \"int4\"): " + err.Error())
			}
		case "text":
			var sentence strings.Builder
			sentence.WriteRune('\'')
			sentence.WriteString(faker.Sentence())
			sentence.WriteRune('\'')
			value = sentence.String()
		case "_text":
			var sentence strings.Builder
			sentence.WriteString("'{\"")
			sentence.WriteString(faker.Sentence())
			sentence.WriteString("\"}'")
			value = sentence.String()
		case "timestamp":
			var timestamp strings.Builder
			timestamp.WriteRune('\'')
			timestamp.WriteString(faker.Timestamp())
			timestamp.WriteRune('\'')
			value = timestamp.String()
		case "json", "jsonb":
			value = `'{ "a": "1", "b": "2" }'`
		case "bool":
			intVal, err := faker.RandomInt(1, 2, 1)
			if err != nil {
				panic("(faker.RandomInt for \"boolean\"): " + err.Error())
			}

			boolean := (intVal[0] % 2) == 0
			if boolean {
				value = "true"
			} else {
				value = "false"
			}
		}

		output.WriteString(value)
	}

	output.WriteString(");")
	return output.String()
}

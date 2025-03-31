package table

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/go-faker/faker/v4"
)

type Table struct {
	Name       string
	Columns    []Column
	InsertRows []map[string]string
}

func New(name string) *Table {
	return &Table{
		Name: name,
	}
}

type Column struct {
	Name                   string
	OrdinalPosition        int
	ColumnDefault          sql.NullString
	IsNullable             string
	DataType               string
	CharacterMaximumLength sql.NullInt32
	CharacterOctetLength   sql.NullInt32
	NumericPrecision       sql.NullInt32
	NumericPrecisionRadix  sql.NullInt32
	NumericScale           sql.NullInt32
	DatetimePrecision      sql.NullInt16
	UdtName                string
	IsSelfReferencing      string
	IsIdentity             string
	IdentityGeneration     sql.NullString
	IdentityStart          sql.NullInt32
	IdentityIncrement      sql.NullInt32
	IdentityMaximum        sql.NullInt32
	IdentityMinimum        sql.NullInt32
	IsUpdateable           string
}

func (t *Table) CreateData(count int) error {
	for range count {
		row := make(map[string]string)

		for _, col := range t.Columns {
			var value string

			switch col.UdtName {
			case "int4":
				intVal, err := faker.RandomInt(0, 128, 1)
				value = strconv.Itoa(intVal[0])
				if err != nil {
					return errors.New("(faker.RandomInt for \"int4\"): " + err.Error())
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
					return errors.New("(faker.RandomInt for \"boolean\"): " + err.Error())
				}

				boolean := (intVal[0] % 2) == 0
				if boolean {
					value = "true"
				} else {
					value = "false"
				}
			default:
				return errors.New("UDT Currently unsupported: " + col.UdtName)
			}

			row[col.Name] = value
		}

		t.InsertRows = append(t.InsertRows, row)
	}

	return nil
}

func (t *Table) ToPsqlStatement() string {
	var output strings.Builder
	output.WriteString("INSERT INTO ")
	output.WriteString(t.Name)
	output.WriteString(" (")

	for i, col := range t.Columns {
		if i > 0 {
			output.WriteRune(',')
		}

		output.WriteString(col.Name)
	}

	output.WriteString(") VALUES ")

	for i := range len(t.InsertRows) {
		if i > 0 {
			output.WriteRune(',')
		}

		output.WriteRune('(')

		for j, col := range t.Columns {
			value := t.InsertRows[i][col.Name]
			if j > 0 {
				output.WriteRune(',')
			}

			output.WriteString(value)
		}

		output.WriteRune(')')
	}

	output.WriteRune(';')
	return output.String()
}

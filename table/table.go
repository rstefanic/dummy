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
	metadata   metadata
	Columns    []Column
	InsertRows [][]string
}

func New(name string) *Table {
	return &Table{
		Name: name,
	}
}

type metadata struct {
	identityColumns []int
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

func (t *Table) FillMetadata() error {
	if len(t.Columns) == 0 {
		return errors.New("Columns on table " + t.Name + " is empty")
	}

	for i, col := range t.Columns {
		if col.IsIdentity == "YES" {
			t.metadata.identityColumns = append(t.metadata.identityColumns, i)
		}
	}

	return nil
}

func (t *Table) CreateData(count int) error {
	for range count {
		var row []string

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

			row = append(row, value)
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

		current := t.InsertRows[i]
		for j, row := range current {
			if j > 0 {
				output.WriteRune(',')
			}

			output.WriteString(row)
		}

		output.WriteRune(')')
	}

	output.WriteRune(';')
	return output.String()
}

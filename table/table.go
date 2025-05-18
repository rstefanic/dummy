package table

import (
	"database/sql"
	"errors"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
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

	return nil
}

func (t *Table) CreateData(count int) error {
	for range count {
		var row []string

		for _, col := range t.Columns {
			if col.IsIdentity == "YES" {
				row = append(row, "DEFAULT")
				continue
			}

			var value string
			value, err := fakeData(col.DataType, col.UdtName)
			if err != nil {
				return err
			}

			row = append(row, value)
		}

		t.InsertRows = append(t.InsertRows, row)
	}

	return nil
}

func fakeData(datatype, udt string) (string, error) {
	switch datatype {
	case "ARRAY":
		underlyingDt, err := udtToPsqlDatatype(udt)
		if err != nil {
			return "", err
		}

		value, err := fakeData(underlyingDt, "")
		if err != nil {
			return "", err
		}

		var array strings.Builder
		array.WriteString("ARRAY[")
		array.WriteString(value)
		array.WriteString("]")
		return array.String(), nil
	case "bigint":
		bigIntVal := gofakeit.Int64()
		return strconv.FormatInt(bigIntVal, 10), nil
	case "bigserial":
		// NOTE: IntRange may cut off the upper bounds of a int64
		bigSerialVal := int64(gofakeit.IntRange(1, math.MaxInt64))
		return strconv.FormatInt(bigSerialVal, 10), nil
	case "boolean":
		boolVal := gofakeit.Bool()
		if boolVal {
			return "true", nil
		} else {
			return "false", nil
		}
	case "integer":
		intVal := gofakeit.Int16()
		return strconv.FormatInt(int64(intVal), 10), nil
	case "json", "jsonb":
		var jo gofakeit.JSONOptions

		// Use gofakeit to create random JSON fields
		err := gofakeit.Struct(&jo)
		if err != nil {
			return "", err
		}

		// Overwrite the fields to force this to be an object
		jo.Indent = false
		jo.RowCount = 1
		jo.Type = "object"

		jsonRaw, err := gofakeit.JSON(&jo)
		if err != nil {
			return "", err
		}

		var json strings.Builder
		json.WriteRune('\'')
		json.WriteString(strings.ReplaceAll(string(jsonRaw), "'", "''")) // escape single quotes
		json.WriteRune('\'')
		return json.String(), err
	case "serial":
		serialVal := gofakeit.IntRange(1, math.MaxInt32)
		return strconv.FormatInt(int64(serialVal), 10), nil
	case "smallserial":
		smallSerialVal := gofakeit.IntRange(1, math.MaxInt16)
		return strconv.FormatInt(int64(smallSerialVal), 10), nil
	case "text":
		var sentence strings.Builder
		sentence.WriteRune('\'')
		sentence.WriteString(strings.ReplaceAll(gofakeit.Sentence(1), "'", "''")) // escape single quotes
		sentence.WriteRune('\'')
		return sentence.String(), nil
	case "timestamp", "timestamp with time zone", "timestamp without time zone":
		var timestamp strings.Builder
		timestamp.WriteRune('\'')
		timestamp.WriteString(gofakeit.Date().Format(time.DateOnly))
		timestamp.WriteRune('\'')
		return timestamp.String(), nil
	default:
		return "", errors.New("Datatype currently unsupported: " + datatype + "(" + udt + ")")
	}
}

func udtToPsqlDatatype(udt string) (string, error) {
	switch udt {
	case "_text", "text":
		return "text", nil
	default:
		return "", errors.New("Unknown UDT to datatype mapping: " + udt)
	}
}

func (t *Table) ToPsqlStatement() string {
	var output strings.Builder

	output.WriteString("INSERT INTO ")
	output.WriteString(t.Name)
	output.WriteString(" (")

	// Write out the column names
	{
		written := 0
		for i, col := range t.Columns {
			// Skip identity columns since we've already
			// generated the data without this column
			if slices.Contains(t.metadata.identityColumns, i) {
				continue
			}

			if written > 0 {
				output.WriteRune(',')
			}

			output.WriteString(col.Name)
			written += 1
		}
	}

	output.WriteString(") VALUES ")

	// Build the main part of the insert statement from the generated data
	for i := range len(t.InsertRows) {
		if i > 0 {
			output.WriteRune(',')
		}

		// Build the current row
		{
			output.WriteRune('(')
			written := 0
			for _, row := range t.InsertRows[i] {
				if written > 0 {
					output.WriteRune(',')
				}

				output.WriteString(row)
				written += 1
			}

			output.WriteRune(')')
		}
	}

	output.WriteRune(';')
	return output.String()
}

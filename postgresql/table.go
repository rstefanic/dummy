package postgresql

import (
	"database/sql"
	"errors"
	"slices"
	"strings"
)

type Table struct {
	Name       string
	metadata   metadata
	Columns    []Column
	InsertRows [][]string
}

func NewTable(name string) *Table {
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

func (t *Table) Validate() error {
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

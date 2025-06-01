package table

import (
	"errors"
	"regexp"
	"strings"

	"dummy/commands"
	"dummy/generate"

	. "dummy/sqldatabase/column"
	. "dummy/sqldatabase/foreignkeyrelation"
)

type Table struct {
	Name       string
	Metadata   Metadata
	Columns    []Column
	InsertRows [][]string
}

func NewTable(name string) *Table {
	return &Table{
		Name: name,
		Metadata: Metadata{
			CustomData: make(map[string]string),
		},
	}
}

type Metadata struct {
	CustomData      map[string]string
	IdentityColumns []int
}

func (t *Table) Validate(cmds commands.TableCommands, fks []ForeignKeyRelation) error {
	if len(t.Columns) == 0 {
		return errors.New("Columns on table " + t.Name + " is empty")
	}

	for _, col := range t.Columns {
		name := col.Name
		cmd, ok := cmds.Columns[name]
		if ok {

			// If it's a text column, ensure that the value requested is something supported
			if col.UdtName == "text" && cmd != "" {
				if !regexp.MustCompile(`(?i)(company|firstname|lastname|name|uuid)`).MatchString(cmd) {
					return errors.New("Column '" + name + "' is not a text column and cannot generate a \"" + cmd + "\" for it.")
				}

				t.Metadata.CustomData[name] = strings.ToLower(cmd)
			}
		}

		// Check if the column we're working with has a FK constraint
		for _, fk := range fks {
			if fk.ColumnName == col.Name {
				if col.IsNullable == "NO" {
					return errors.New("Column '" + col.Name + "' has a FK constraint named '" + fk.ConstraintName + "' that is not nullable")
				}
			}
		}
	}

	return nil
}

func (t *Table) GuessCustomTextFieldGenerators() {
	customData := t.Metadata.CustomData

	for _, col := range t.Columns {
		_, exists := customData[col.Name]
		if exists {
			continue
		}

		if col.DataType != "text" {
			continue
		}

		// Take a look at the column name to see if we can't guess the correct data
		colName := col.Name
		switch true {
		case regexp.MustCompile(`(?i)^first[_-]*name$`).MatchString(colName):
			fallthrough
		case regexp.MustCompile(`(?i)^given[_-]*name$`).MatchString(colName):
			customData[colName] = "firstname"
			continue
		case regexp.MustCompile(`(?i)^last[_-]*name$`).MatchString(colName):
			fallthrough
		case regexp.MustCompile(`(?i)^family[_-]*name$`).MatchString(colName):
			customData[colName] = "lastname"
			continue
		case regexp.MustCompile(`(?i)^(company|firm|business|corporation|establishment|organization|institution)[-_]*(name)?$`).MatchString(colName):
			customData[colName] = "company"
			continue
		}

		// Try to see if the table name provides a clue
		tableName := t.Name
		switch true {
		case regexp.MustCompile(`(?i)users`).MatchString(tableName):
			switch true {
			case regexp.MustCompile(`(?i)^id$`).MatchString(colName):
				customData[colName] = "uuid"
			case regexp.MustCompile(`(?i)^name$`).MatchString(colName):
				customData[colName] = "name"
			default:
				// No guesses to offer
				continue
			}
		case regexp.MustCompile(`(?i)(company|firm|business|corporation|establishment|organization|institution)`).MatchString(tableName):
			switch true {
			case regexp.MustCompile(`(?i)^id$`).MatchString(colName):
				customData[colName] = "uuid"
			case regexp.MustCompile(`(?i)^name$`).MatchString(colName):
				customData[colName] = "company"
			default:
				// No guesses to offer
				continue
			}
		}
	}
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
			value, err := generate.FakeData(col.DataType, col.UdtName, col.Name, &t.Metadata.CustomData)
			if err != nil {
				return err
			}

			row = append(row, value)
		}

		t.InsertRows = append(t.InsertRows, row)
	}

	return nil
}

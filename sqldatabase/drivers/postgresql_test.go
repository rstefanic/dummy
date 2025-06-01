package drivers

import (
	"strings"
	"testing"

	"dummy/commands"
	. "dummy/sqldatabase/column"
	. "dummy/sqldatabase/foreignkeyrelation"
	. "dummy/sqldatabase/table"
)

func createFakeColumn(name string, ordinalPosition int, isNullable bool, udtName string, isIdentity bool) *Column {
	var col *Column

	col = &Column{}
	col.Name = name
	col.OrdinalPosition = ordinalPosition
	col.DataType = udtName // this is not the same, but good enough for testing
	col.UdtName = udtName
	col.IsSelfReferencing = "NO"

	if isNullable {
		col.IsNullable = "YES"
	} else {
		col.IsNullable = "NO"
	}

	if isIdentity {
		col.IsIdentity = "YES"
	} else {
		col.IsIdentity = "NO"
	}

	return col
}

func createFakeTable(tableName string) *Table {
	table := NewTable(tableName)

	idCol := createFakeColumn("id", 1, false, "int4", true)
	table.Columns = append(table.Columns, *idCol)

	nameCol := createFakeColumn("name", 1, false, "text", false)
	table.Columns = append(table.Columns, *nameCol)

	createdAtCol := createFakeColumn("created_at", 1, false, "timestamp without time zone", false)
	table.Columns = append(table.Columns, *createdAtCol)

	return table
}

func TestToPsqlStatement(t *testing.T) {
	driver := PostgresqlDriver{database: nil}

	var tblCmds commands.TableCommands
	table := createFakeTable("fake_table")
	table.Validate(tblCmds, make([]ForeignKeyRelation, 0))

	table.InsertRows = append(table.InsertRows, []string{"DEFAULT", "Bill Bob", "2025-04-12 10:00:00 UTC"})
	table.InsertRows = append(table.InsertRows, []string{"DEFAULT", "Jim George", "2025-04-12 10:00:00 UTC"})

	actual := driver.InsertStatement(table)
	expected := "INSERT INTO fake_table (id,name,created_at) VALUES (DEFAULT,Bill Bob,2025-04-12 10:00:00 UTC),(DEFAULT,Jim George,2025-04-12 10:00:00 UTC);"

	if strings.Compare(actual, expected) != 0 {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, actual)
	}
}

package table

import (
	"strings"
	"testing"
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

	if isIdentity{
		col.IsIdentity = "YES"
	} else {
		col.IsIdentity = "NO"
	}

	return col
}

func createFakeTable(tableName string) *Table {
	table := New(tableName)

	idCol := createFakeColumn("id", 1, false, "int4", true)
	table.Columns = append(table.Columns, *idCol)

	nameCol := createFakeColumn("name", 1, false, "text", false)
	table.Columns = append(table.Columns, *nameCol)

	createdAtCol := createFakeColumn("created_at", 1, false, "timestamp", false)
	table.Columns = append(table.Columns, *createdAtCol)

	return table
}

func TestCreateData(t *testing.T) {
	table := createFakeTable("fake_table")
	err := table.CreateData(3)
	if err != nil {
		t.Errorf(`Error calling "table.CreateData": %v`, err)
	}

	rowsCount := len(table.InsertRows)
	if rowsCount != 3 {
		t.Errorf(`Expected "table.CreateData" to create 3 rows, got %d rows.`, rowsCount)
	}
}

func TestToPsqlStatement(t *testing.T) {
	table := createFakeTable("fake_table")
	table.FillMetadata()

	table.InsertRows = append(table.InsertRows, []string{"Bill Bob", "2025-04-12 10:00:00 UTC"})
	table.InsertRows = append(table.InsertRows, []string{"Jim George", "2025-04-12 10:00:00 UTC"})

	actual := table.ToPsqlStatement()
	expected := "INSERT INTO fake_table (name,created_at) VALUES (Bill Bob,2025-04-12 10:00:00 UTC),(Jim George,2025-04-12 10:00:00 UTC);"

	if strings.Compare(actual, expected) != 0 {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, actual)
	}
}

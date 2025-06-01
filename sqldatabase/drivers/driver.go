package drivers

import (
	"database/sql"

	. "dummy/sqldatabase/column"
	. "dummy/sqldatabase/foreignkeyrelation"
	. "dummy/sqldatabase/table"
)

type SqlDatabaseDriver interface {
	Database() *sql.DB
	ForeignKeyRelations() (map[string][]ForeignKeyRelation, error)
	TableColumns(tableName string) ([]Column, error)
	InsertStatement(table *Table) string
}

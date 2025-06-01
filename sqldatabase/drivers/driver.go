package drivers

import (
	"database/sql"

	. "dummy/sqldatabase/column"
	. "dummy/sqldatabase/foreignkeyrelation"
)

type SqlDatabaseDriver interface {
	Database() *sql.DB
	ForeignKeyRelations() (map[string][]ForeignKeyRelation, error)
	TableColumns(tableName string) ([]Column, error)
}

package sqldatabase

import (
	"dummy/sqldatabase/drivers"

	. "dummy/sqldatabase/foreignkeyrelation"
	. "dummy/sqldatabase/table"
)

// SqlDatabase maintains the driver which is a handle to the underlying
// SqlDatabase connection with information about the database that we
// are connected to.
type SqlDatabase struct {
	Driver      drivers.SqlDatabaseDriver
	ForeignKeys map[string][]ForeignKeyRelation
	Tables      []Table
}

func New(driver drivers.SqlDatabaseDriver) (*SqlDatabase, error) {
	fks, err := driver.ForeignKeyRelations()
	if err != nil {
		return nil, err
	}

	return &SqlDatabase{
		Driver:      driver,
		ForeignKeys: fks,
		Tables:      make([]Table, 0),
	}, nil
}

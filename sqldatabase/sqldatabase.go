package sqldatabase

import (
	"dummy/sqldatabase/drivers"
	fkr "dummy/sqldatabase/foreignkeyrelation"
)

// SqlDatabase maintains the driver which is a handle to the underlying
// SqlDatabase connection with information about the database that we
// are connected to.
type SqlDatabase struct {
	Driver      drivers.SqlDatabaseDriver
	ForeignKeys map[string][]fkr.ForeignKeyRelation
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

func (sqlDb *SqlDatabase) Close() {
	sqlDb.Driver.Database().Close()
}

func (sqlDb *SqlDatabase) GetTableColumns(table string) ([]Column, error) {
	var columns []Column

	rows, err := sqlDb.Driver.Database().Query(`
		SELECT column_name, ordinal_position, column_default, is_nullable, data_type,
		character_maximum_length, character_octet_length, numeric_precision,
		numeric_precision_radix, numeric_scale, datetime_precision, udt_name,
		is_self_referencing, is_identity, identity_generation, identity_start,
		identity_increment, identity_maximum, identity_minimum, is_updatable
			FROM information_schema.columns
			WHERE table_name = $1
			ORDER BY ordinal_position`,
		table,
	)

	if err != nil {
		return make([]Column, 0), err
	}

	defer rows.Close()

	for rows.Next() {
		var c Column
		err := rows.Scan(&c.Name, &c.OrdinalPosition, &c.ColumnDefault, &c.IsNullable,
			&c.DataType, &c.CharacterMaximumLength, &c.CharacterOctetLength,
			&c.NumericPrecision, &c.NumericPrecisionRadix, &c.NumericScale,
			&c.DatetimePrecision, &c.UdtName, &c.IsSelfReferencing,
			&c.IsIdentity, &c.IdentityGeneration, &c.IdentityStart,
			&c.IdentityIncrement, &c.IdentityMaximum, &c.IdentityMinimum,
			&c.IsUpdateable)

		if err != nil {
			return make([]Column, 0), err
		}

		columns = append(columns, c)
	}

	return columns, nil
}

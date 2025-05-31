package sqldatabase

import (
	"database/sql"

	"dummy/sqldatabase/drivers"
)

// SqlDatabase maintains the driver which is a handle to the underlying
// SqlDatabase connection with information about the database that we
// are connected to.
type SqlDatabase struct {
	Driver      drivers.SqlDatabaseDriver
	ForeignKeys map[string][]ForeignKeyRelation
	Tables      []Table
}

type ForeignKeyRelation struct {
	TableSchema        string
	ConstraintName     string
	TableName          string
	ColumnName         string
	ForeignTableSchema string
	ForeignTableName   string
	ForeignColumnName  string
}

func New(driver drivers.SqlDatabaseDriver) (*SqlDatabase, error) {
	fks, err := foreignKeyRelations(driver.Database())
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

func foreignKeyRelations(db *sql.DB) (map[string][]ForeignKeyRelation, error) {
	fkMapping := make(map[string][]ForeignKeyRelation)
	fks, err := queryForeignKeyRelations(db)

	if err != nil {
		return make(map[string][]ForeignKeyRelation), err
	}

	for _, fk := range fks {
		table := fk.TableName
		val, ok := fkMapping[table]
		if !ok {
			val = make([]ForeignKeyRelation, 0)
		}

		val = append(val, fk)
		fkMapping[table] = val
	}

	return fkMapping, nil
}

func queryForeignKeyRelations(db *sql.DB) ([]ForeignKeyRelation, error) {
	var fks []ForeignKeyRelation

	rows, err := db.Query(`
		SELECT
			tc.table_schema, 
			tc.constraint_name, 
			tc.table_name, 
			kcu.column_name, 
			ccu.table_schema AS foreign_table_schema,
			ccu.table_name AS foreign_table_name,
			ccu.column_name AS foreign_column_name 
		FROM information_schema.table_constraints AS tc 
		JOIN information_schema.key_column_usage AS kcu
			ON tc.constraint_name = kcu.constraint_name
			AND tc.table_schema = kcu.table_schema
		JOIN information_schema.constraint_column_usage AS ccu
			ON ccu.constraint_name = tc.constraint_name
		WHERE tc.constraint_type = 'FOREIGN KEY'
		ORDER BY table_name`,
	)

	if err != nil {
		return make([]ForeignKeyRelation, 0), err
	}

	defer rows.Close()

	for rows.Next() {
		var fk ForeignKeyRelation
		err := rows.Scan(
			&fk.TableSchema,
			&fk.ConstraintName,
			&fk.TableName,
			&fk.ColumnName,
			&fk.ForeignTableSchema,
			&fk.ForeignTableName,
			&fk.ForeignColumnName,
		)

		if err != nil {
			return make([]ForeignKeyRelation, 0), err
		}

		fks = append(fks, fk)
	}

	return fks, nil
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

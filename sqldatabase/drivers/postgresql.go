package drivers

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"

	. "dummy/sqldatabase/column"
	. "dummy/sqldatabase/foreignkeyrelation"
	. "dummy/sqldatabase/table"
)

type PostgresqlDriver struct {
	database *sql.DB
}

func NewPostgresqlDriver(user, password, host, name string) (*PostgresqlDriver, error) {
	conn := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", user, password, host, name)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return &PostgresqlDriver{database: db}, nil
}

func (pd *PostgresqlDriver) Database() *sql.DB {
	return pd.database
}

func (pd *PostgresqlDriver) TableColumns(tableName string) ([]Column, error) {
	var columns []Column

	rows, err := pd.Database().Query(`
		SELECT column_name, ordinal_position, column_default, is_nullable, data_type,
		character_maximum_length, character_octet_length, numeric_precision,
		numeric_precision_radix, numeric_scale, datetime_precision, udt_name,
		is_self_referencing, is_identity, identity_generation, identity_start,
		identity_increment, identity_maximum, identity_minimum, is_updatable
			FROM information_schema.columns
			WHERE table_name = $1
			ORDER BY ordinal_position`,
		tableName,
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

func (pd *PostgresqlDriver) ForeignKeyRelations() (map[string][]ForeignKeyRelation, error) {
	db := pd.Database()
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

func (pd *PostgresqlDriver) InsertStatement(t *Table) string {
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
			if slices.Contains(t.Metadata.IdentityColumns, i) {
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

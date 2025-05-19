package postgresql

import (
	"database/sql"
	"fmt"
)

type PostgresqlDB struct {
	db *sql.DB
}

func New(user, password, host, name string) (*PostgresqlDB, error) {
	conn := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", user, password, host, name)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return &PostgresqlDB{
		db: db,
	}, nil
}

func (p *PostgresqlDB) Close() {
	p.db.Close()
}

func (p *PostgresqlDB) GetTableColumns(table string) ([]Column, error) {
	var columns []Column

	rows, err := p.db.Query(`
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

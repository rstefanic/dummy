package drivers

import (
	"database/sql"
	"fmt"

	fkr "dummy/sqldatabase/foreignkeyrelation"
)

type SqlDatabaseDriver interface {
	Database() *sql.DB
	ForeignKeyRelations() (map[string][]fkr.ForeignKeyRelation, error)
}

type PostgresDriver struct {
	database *sql.DB
}

func NewPostgresDriver(user, password, host, name string) (*PostgresDriver, error) {
	conn := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", user, password, host, name)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return &PostgresDriver{database: db}, nil
}

func (pd *PostgresDriver) Database() *sql.DB {
	return pd.database
}

func (pd *PostgresDriver) ForeignKeyRelations() (map[string][]fkr.ForeignKeyRelation, error) {
	db := pd.Database()
	fkMapping := make(map[string][]fkr.ForeignKeyRelation)
	fks, err := queryForeignKeyRelations(db)

	if err != nil {
		return make(map[string][]fkr.ForeignKeyRelation), err
	}

	for _, fk := range fks {
		table := fk.TableName
		val, ok := fkMapping[table]
		if !ok {
			val = make([]fkr.ForeignKeyRelation, 0)
		}

		val = append(val, fk)
		fkMapping[table] = val
	}

	return fkMapping, nil
}

func queryForeignKeyRelations(db *sql.DB) ([]fkr.ForeignKeyRelation, error) {
	var fks []fkr.ForeignKeyRelation

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
		return make([]fkr.ForeignKeyRelation, 0), err
	}

	defer rows.Close()

	for rows.Next() {
		var fk fkr.ForeignKeyRelation
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
			return make([]fkr.ForeignKeyRelation, 0), err
		}

		fks = append(fks, fk)
	}

	return fks, nil
}

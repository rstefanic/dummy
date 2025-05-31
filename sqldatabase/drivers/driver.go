package drivers

import (
	"database/sql"
	"fmt"
)

type SqlDatabaseDriver interface {
	Database() *sql.DB
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

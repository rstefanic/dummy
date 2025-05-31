package column

import "database/sql"

type Column struct {
	Name                   string
	OrdinalPosition        int
	ColumnDefault          sql.NullString
	IsNullable             string
	DataType               string
	CharacterMaximumLength sql.NullInt32
	CharacterOctetLength   sql.NullInt32
	NumericPrecision       sql.NullInt32
	NumericPrecisionRadix  sql.NullInt32
	NumericScale           sql.NullInt32
	DatetimePrecision      sql.NullInt16
	UdtName                string
	IsSelfReferencing      string
	IsIdentity             string
	IdentityGeneration     sql.NullString
	IdentityStart          sql.NullInt32
	IdentityIncrement      sql.NullInt32
	IdentityMaximum        sql.NullInt32
	IdentityMinimum        sql.NullInt32
	IsUpdateable           string
}

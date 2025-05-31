package foreignkeyrelation

type ForeignKeyRelation struct {
	TableSchema        string
	ConstraintName     string
	TableName          string
	ColumnName         string
	ForeignTableSchema string
	ForeignTableName   string
	ForeignColumnName  string
}

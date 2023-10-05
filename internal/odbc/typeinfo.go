package odbc

import "github.com/ninthclowd/unixodbc/internal/api"

type TypeInfo struct {
	DataType               api.SQLINTEGER  `col_name:"DATA_TYPE"`
	TypeName               string          `col_name:"TYPE_NAME"`
	ColumnSize             api.SQLULEN     `col_name:"COLUMN_SIZE"`
	LiteralPrefix          string          `col_name:"LITERAL_PREFIX"`
	LiteralSuffix          string          `col_name:"LITERAL_SUFFIX"`
	CreateParams           string          `col_name:"CREATE_PARAMS"`
	Nullable               api.SQLSMALLINT `col_name:"NULLABLE"`
	IsCaseSensitive        api.SQLUINTEGER `col_name:"CASE_SENSITIVE"`
	Searchable             api.SQLSMALLINT `col_name:"SEARCHABLE"`
	IsUnsigned             api.SQLUINTEGER `col_name:"UNSIGNED_ATTRIBUTE"`
	FixedPrecisionAndScale api.SQLUINTEGER `col_name:"FIXED_PREC_SCALE"`
	IsAutoIncrementing     api.SQLUINTEGER `col_name:"AUTO_UNIQUE_VALUE"`
	LocalTypeName          string          `col_name:"LOCAL_TYPE_NAME"`
	MinimumScale           api.SQLSMALLINT `col_name:"MINIMUM_SCALE"`
	MaximumScale           api.SQLSMALLINT `col_name:"MAXIMUM_SCALE"`
	SqlDataType            api.SQLINTEGER  `col_name:"SQL_DATA_TYPE"`
	DateTimeSubCode        api.SQLSMALLINT `col_name:"SQL_DATETIME_SUB"`
	NumPrecRadix           api.SQLINTEGER  `col_name:"NUM_PREC_RADIX"`
	IntervalPrecision      api.SQLSMALLINT `col_name:"INTERVAL_PRECISION"`
}

type TableType string

var (
	TypeTable      TableType = "TABLE"
	TypeView       TableType = "VIEW"
	TypeSysTable   TableType = "SYSTEM TABLE"
	TypeGlobalTemp TableType = "GLOBAL TEMPORARY"
	TypeLocalTemp  TableType = "LOCAL TEMPORARY"
	TypeAlias      TableType = "ALIAS"
	TypeSynonym    TableType = "SYNONYM"
)

type TableInfo struct {
	Catalog *string   `col_name:"TABLE_CAT"`
	Schema  *string   `col_name:"TABLE_SCHEM"`
	Name    string    `col_name:"TABLE_NAME"`
	Type    TableType `col_name:"TABLE_TYPE"`
	Remarks string    `col_name:"REMARKS"`
}

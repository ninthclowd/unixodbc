package odbc

import (
	"database/sql/driver"
	"github.com/ninthclowd/unixodbc/internal/api"
	"reflect"
)

func init() {
	registerColumnFactoryForType(newUTF8Column,
		api.SQL_CHAR,
		api.SQL_VARCHAR,
		api.SQL_LONGVARCHAR,
	)
}

func newUTF8Column(info *columnInfo, hnd handle) Column {
	return &columnUTF8{hnd, info}
}

type columnUTF8 struct {
	handle
	*columnInfo
}

func (c *columnUTF8) VariableLength() (length int64, ok bool) {
	return int64(c.columnSize), true
}

func (c *columnUTF8) ScanType() reflect.Type {
	return reflect.TypeOf((*string)(nil))
}

func (c *columnUTF8) Decimal() (precision int64, scale int64, ok bool) {
	return
}

func (c *columnUTF8) Value() (driver.Value, error) {
	value := make([]byte, c.columnSize)
	var valueLength api.SQLLEN
	if _, err := c.result(c.api().SQLGetData(api.SQLHSTMT(c.hnd()), c.columnNumber, api.SQL_C_CHAR, api.SQLPOINTER(&value[0]), api.SQLLEN(len(value)), &valueLength)); err != nil {
		return nil, err
	}
	if valueLength == api.SQL_NULL_DATA {
		return nil, nil
	}

	return string(value[:valueLength]), nil
}

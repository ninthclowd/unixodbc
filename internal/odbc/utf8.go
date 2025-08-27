package odbc

import (
	"database/sql/driver"
	"reflect"
	"unsafe"

	"github.com/ninthclowd/unixodbc/internal/api"
)

func init() {
	registerColumnFactoryForType(newUTF8Column,
		api.SQL_CHAR,
		api.SQL_VARCHAR,
		api.SQL_LONGVARCHAR,
	)
}

func newUTF8Column(info *columnInfo, hnd *handle) Column {
	return &columnUTF8{hnd, info}
}

type columnUTF8 struct {
	*handle
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

//go:nocheckptr
func (c *columnUTF8) Value() (driver.Value, error) {
	value := make([]uint8, c.columnSize+1)
	var valueLength api.SQLLEN
	maxLen := api.SQLLEN(len(value))
	if _, err := c.result(api.SQLGetData((*api.SQLHSTMT)(c.hnd()),
		c.columnNumber,
		api.SQL_C_CHAR,
		(*api.SQLPOINTER)(unsafe.Pointer(&value[0])),
		maxLen,
		&valueLength)); err != nil {
		return nil, err
	}
	if valueLength == api.SQL_NULL_DATA {
		return nil, nil
	}
	if valueLength > maxLen {
		valueLength = maxLen
	}
	out := string(value[:valueLength])
	value = nil //zero out for GC
	return out, nil
}

package odbc

import (
	"database/sql/driver"
	"github.com/ninthclowd/unixodbc/internal/api"
	"reflect"
	"unsafe"
)

func init() {
	registerColumnFactoryForType(newFloat32Column, api.SQL_REAL)
}

func newFloat32Column(info *columnInfo, hnd *handle) Column {
	return &columnFloat32{hnd, info}
}

type columnFloat32 struct {
	*handle
	*columnInfo
}

func (c *columnFloat32) VariableLength() (length int64, ok bool) {
	return
}

func (c *columnFloat32) ScanType() reflect.Type {
	return reflect.TypeOf((*float32)(nil))
}

func (c *columnFloat32) Decimal() (precision int64, scale int64, ok bool) {
	return int64(c.columnSize), int64(c.decimalDigits), true
}

//go:nocheckptr
func (c *columnFloat32) Value() (driver.Value, error) {
	var value api.SQLREAL
	var valueLength api.SQLLEN
	if _, err := c.result(api.SQLGetData((*api.SQLHSTMT)(c.hnd()),
		c.columnNumber,
		api.SQL_C_FLOAT,
		(*api.SQLPOINTER)(unsafe.Pointer(&value)),
		0,
		&valueLength)); err != nil {
		return nil, err
	}
	if valueLength == api.SQL_NULL_DATA {
		return nil, nil
	}
	return float32(value), nil
}

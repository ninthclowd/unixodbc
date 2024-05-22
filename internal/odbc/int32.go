package odbc

import (
	"database/sql/driver"
	"github.com/ninthclowd/unixodbc/internal/api"
	"reflect"
	"unsafe"
)

func init() {
	registerColumnFactoryForType(newInt32Column, api.SQL_INTEGER)
}

func newInt32Column(info *columnInfo, hnd *handle) Column {
	return &columnInt32{hnd, info}
}

type columnInt32 struct {
	*handle
	*columnInfo
}

func (c *columnInt32) VariableLength() (length int64, ok bool) {
	return
}

func (c *columnInt32) ScanType() reflect.Type {
	return reflect.TypeOf((*int32)(nil))
}

func (c *columnInt32) Decimal() (precision int64, scale int64, ok bool) {
	return
}

func (c *columnInt32) Value() (driver.Value, error) {
	var value api.SQLBIGINT
	var valueLength api.SQLLEN
	if _, err := c.result(api.SQLGetData((*api.SQLHSTMT)(c.hnd()),
		c.columnNumber,
		api.SQL_C_SLONG,
		(*api.SQLPOINTER)(unsafe.Pointer(&value)),
		0,
		&valueLength)); err != nil {
		return nil, err
	}
	if valueLength == api.SQL_NULL_DATA {
		return nil, nil
	}
	return value, nil
}

//go:nocheckptr
func (s *statement) bindInt32(index int, value int32) error {
	_, err := s.result(api.SQLBindParameter((*api.SQLHSTMT)(s.hnd()),
		api.SQLUSMALLINT(index+1),
		api.SQL_PARAM_INPUT,
		api.SQL_C_SLONG, api.SQL_INTEGER,
		0, 0,
		(*api.SQLPOINTER)(unsafe.Pointer(&value)),
		0, nil))
	return err
}

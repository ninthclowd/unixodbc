package odbc

import (
	"database/sql/driver"
	"github.com/ninthclowd/unixodbc/internal/api"
	"reflect"
	"unsafe"
)

func init() {
	registerColumnFactoryForType(newBoolColumn, api.SQL_BIT)
}

func newBoolColumn(info *columnInfo, hnd *handle) Column {
	return &columnBool{hnd, info}
}

type columnBool struct {
	*handle
	*columnInfo
}

func (c *columnBool) VariableLength() (length int64, ok bool) {
	return
}

func (c *columnBool) ScanType() reflect.Type {
	return reflect.TypeOf((*bool)(nil))
}

func (c *columnBool) Decimal() (precision int64, scale int64, ok bool) {
	return
}

func (c *columnBool) Value() (driver.Value, error) {
	var value api.SQLCHAR
	var valueLength api.SQLLEN
	if _, err := c.result(api.SQLGetData((*api.SQLHSTMT)(c.hnd()),
		c.columnNumber,
		api.SQL_C_BIT,
		(*api.SQLPOINTER)(unsafe.Pointer(&value)),
		0,
		&valueLength)); err != nil {
		return nil, err
	}
	if valueLength == api.SQL_NULL_DATA {
		return nil, nil
	}
	return value == 1, nil
}

//go:nocheckptr
func (s *statement) bindBool(index int, value bool) error {
	var data byte
	if value {
		data = 1
	}
	_, err := s.result(api.SQLBindParameter((*api.SQLHSTMT)(s.hnd()),
		api.SQLUSMALLINT(index+1),
		api.SQL_PARAM_INPUT,
		api.SQL_C_BIT,
		api.SQL_BIT,
		1,
		0,
		(*api.SQLPOINTER)(unsafe.Pointer(&data)),
		0,
		nil))
	return err
}

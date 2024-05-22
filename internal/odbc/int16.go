package odbc

import (
	"database/sql/driver"
	"github.com/ninthclowd/unixodbc/internal/api"
	"reflect"
	"unsafe"
)

func init() {
	registerColumnFactoryForType(newInt16Column, api.SQL_SMALLINT)
}

func newInt16Column(info *columnInfo, hnd *handle) Column {
	return &columnInt16{hnd, info}
}

type columnInt16 struct {
	*handle
	*columnInfo
}

func (c *columnInt16) VariableLength() (length int64, ok bool) {
	return
}

func (c *columnInt16) ScanType() reflect.Type {
	return reflect.TypeOf((*int16)(nil))
}

func (c *columnInt16) Decimal() (precision int64, scale int64, ok bool) {
	return
}

func (c *columnInt16) Value() (driver.Value, error) {
	var value api.SQLSMALLINT
	var valueLength api.SQLLEN
	if _, err := c.result(api.SQLGetData((*api.SQLHSTMT)(c.hnd()),
		c.columnNumber,
		api.SQL_C_SSHORT,
		(*api.SQLPOINTER)(unsafe.Pointer(&value)),
		0,
		&valueLength)); err != nil {
		return nil, err
	}
	if valueLength == api.SQL_NULL_DATA {
		return nil, nil
	}
	return int16(value), nil
}

//go:nocheckptr
func (s *Statement) bindInt16(index int, value int16) error {
	_, err := s.result(api.SQLBindParameter((*api.SQLHSTMT)(s.hnd()),
		api.SQLUSMALLINT(index+1),
		api.SQL_PARAM_INPUT,
		api.SQL_C_SSHORT,
		api.SQL_SMALLINT,
		0,
		0,
		(*api.SQLPOINTER)(unsafe.Pointer(&value)),
		0,
		nil))
	return err
}

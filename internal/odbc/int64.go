package odbc

import (
	"database/sql/driver"
	"github.com/ninthclowd/unixodbc/internal/api"
	"reflect"
	"unsafe"
)

func init() {
	registerColumnFactoryForType(newInt64Column, api.SQL_BIGINT)
}

func newInt64Column(info *columnInfo, hnd *handle) Column {
	return &columnInt64{hnd, info}
}

type columnInt64 struct {
	*handle
	*columnInfo
}

func (c *columnInt64) VariableLength() (length int64, ok bool) {
	return
}

func (c *columnInt64) ScanType() reflect.Type {
	return reflect.TypeOf((*int16)(nil))
}

func (c *columnInt64) Decimal() (precision int64, scale int64, ok bool) {
	return
}

func (c *columnInt64) Value() (driver.Value, error) {
	var value api.SQLBIGINT
	var valueLength api.SQLLEN
	if _, err := c.result(api.SQLGetData((*api.SQLHSTMT)(c.hnd()), c.columnNumber, api.SQL_C_SBIGINT,

		(*api.SQLPOINTER)(unsafe.Pointer(&value)),
		0, &valueLength)); err != nil {
		return nil, err
	}
	if valueLength == api.SQL_NULL_DATA {
		return nil, nil
	}
	return int64(value), nil
}

func (s *statement) bindInt64(index int, value int64) error {
	_, err := s.result(api.SQLBindParameter((*api.SQLHSTMT)(s.hnd()),
		api.SQLUSMALLINT(index+1),
		api.SQL_PARAM_INPUT,
		api.SQL_C_SBIGINT,
		api.SQL_BIGINT,
		0,
		0,
		(*api.SQLPOINTER)(unsafe.Pointer(&value)),
		0,
		nil))
	return err
}

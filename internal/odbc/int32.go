package odbc

import (
	"database/sql/driver"
	"github.com/ninthclowd/unixodbc/internal/api"
	"reflect"
)

func init() {
	registerColumnFactoryForType(newInt32Column, api.SQL_INTEGER)
}

func newInt32Column(info *columnInfo, hnd handle) Column {
	return &columnInt32{hnd, info}
}

type columnInt32 struct {
	handle
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
	var value api.SQLINTEGER
	var valueLength api.SQLLEN
	if _, err := c.result(c.api().SQLGetData(api.SQLHSTMT(c.hnd()), c.columnNumber, api.SQL_C_SLONG, api.SQLPOINTER(&value), 0, &valueLength)); err != nil {
		return nil, err
	}
	if valueLength == api.SQL_NULL_DATA {
		return nil, nil
	}
	return int32(value), nil
}

func (s *Statement) bindInt32(index int, value int32) error {
	_, err := s.result(s.api().SQLBindParameter((api.SQLHSTMT)(s.hnd()), api.SQLUSMALLINT(index+1), api.SQL_PARAM_INPUT,
		api.SQL_C_SLONG, api.SQL_INTEGER,
		0, 0,
		api.SQLPOINTER(&value),
		0, nil))
	return err
}

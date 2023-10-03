package odbc

import (
	"database/sql/driver"
	"github.com/ninthclowd/unixodbc/internal/api"
	"reflect"
)

func init() {
	registerColumnFactoryForType(newInt8Column, api.SQL_TINYINT)
}

func newInt8Column(info *columnInfo, hnd handle) Column {
	return &columnInt8{hnd, info}
}

type columnInt8 struct {
	handle
	*columnInfo
}

func (c *columnInt8) VariableLength() (length int64, ok bool) {
	return
}

func (c *columnInt8) ScanType() reflect.Type {
	return reflect.TypeOf((*int8)(nil))
}

func (c *columnInt8) Decimal() (precision int64, scale int64, ok bool) {
	return
}

func (c *columnInt8) Value() (driver.Value, error) {
	var value api.SQLSCHAR
	var valueLength api.SQLLEN
	if _, err := c.result(c.api().SQLGetData(api.SQLHSTMT(c.hnd()), c.columnNumber, api.SQL_C_STINYINT, api.SQLPOINTER(&value), 0, &valueLength)); err != nil {
		return nil, err
	}
	if valueLength == api.SQL_NULL_DATA {
		return nil, nil
	}
	return int8(value), nil
}

func (s *Statement) bindInt8(index int, value int8) error {
	_, err := s.result(s.api().SQLBindParameter((api.SQLHSTMT)(s.hnd()), api.SQLUSMALLINT(index+1), api.SQL_PARAM_INPUT,
		api.SQL_C_STINYINT, api.SQL_TINYINT,
		0, 0,
		api.SQLPOINTER(&value),
		0, nil))
	return err
}

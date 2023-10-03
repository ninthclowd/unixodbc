package odbc

import (
	"database/sql/driver"
	"github.com/ninthclowd/unixodbc/internal/api"
	"reflect"
)

func init() {
	registerColumnFactoryForType(newFloat64Column,
		api.SQL_DOUBLE,
		api.SQL_FLOAT,
		api.SQL_DECIMAL,
		api.SQL_NUMERIC,
	)
}

func newFloat64Column(info *columnInfo, hnd handle) Column {
	return &columnFloat64{hnd, info}
}

type columnFloat64 struct {
	handle
	*columnInfo
}

func (c *columnFloat64) VariableLength() (length int64, ok bool) {
	return
}

func (c *columnFloat64) ScanType() reflect.Type {
	return reflect.TypeOf((*float64)(nil))
}

func (c *columnFloat64) Decimal() (precision int64, scale int64, ok bool) {
	return int64(c.columnSize), int64(c.decimalDigits), true
}

func (c *columnFloat64) Value() (driver.Value, error) {
	var value api.SQLDOUBLE
	var valueLength api.SQLLEN
	if _, err := c.result(c.api().SQLGetData(api.SQLHSTMT(c.hnd()), c.columnNumber, api.SQL_C_DOUBLE, api.SQLPOINTER(&value), 0, &valueLength)); err != nil {
		return nil, err
	}
	if valueLength == api.SQL_NULL_DATA {
		return nil, nil
	}
	return float64(value), nil
}

func (s *Statement) bindFloat64(index int, value float64) error {
	_, err := s.result(s.api().SQLBindParameter((api.SQLHSTMT)(s.hnd()), api.SQLUSMALLINT(index+1), api.SQL_PARAM_INPUT,
		api.SQL_C_DOUBLE, api.SQL_DOUBLE,
		0, 0,
		api.SQLPOINTER(&value),
		0, nil))
	return err
}

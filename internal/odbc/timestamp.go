package odbc

import (
	"database/sql/driver"
	"github.com/ninthclowd/unixodbc/internal/api"
	"reflect"
	"time"
	"unsafe"
)

func init() {
	registerColumnFactoryForType(newTimestampColumn, api.SQL_TYPE_TIMESTAMP)
}

func newTimestampColumn(info *columnInfo, hnd *handle) Column {
	return &columnTimestamp{hnd, info}
}

type columnTimestamp struct {
	*handle
	*columnInfo
}

func (c *columnTimestamp) VariableLength() (length int64, ok bool) {
	return
}

func (c *columnTimestamp) ScanType() reflect.Type {
	return reflect.TypeOf((*time.Time)(nil))
}

func (c *columnTimestamp) Decimal() (precision int64, scale int64, ok bool) {
	return
}

func (c *columnTimestamp) Value() (driver.Value, error) {
	var value api.SQL_TIMESTAMP_STRUCT
	defer value.Free()
	var valueLength api.SQLLEN
	if _, err := c.result(api.SQLGetData((*api.SQLHSTMT)(c.hnd()),
		c.columnNumber,
		api.SQL_C_TIMESTAMP,
		(*api.SQLPOINTER)(unsafe.Pointer(&value)),
		0,
		&valueLength)); err != nil {
		return nil, err
	}
	if valueLength == api.SQL_NULL_DATA {
		return nil, nil
	}
	return time.Date(int(value.Year), time.Month(value.Month), int(value.Day), int(value.Hour), int(value.Minute), int(value.Second), int(value.Fraction), time.UTC), nil
}

func (s *Statement) bindTimestamp(index int, value *time.Time) error {
	ts := api.SQL_TIMESTAMP_STRUCT{
		Year:     api.SQLSMALLINT(value.Year()),
		Month:    api.SQLUSMALLINT(value.Month()),
		Day:      api.SQLUSMALLINT(value.Day()),
		Hour:     api.SQLUSMALLINT(value.Hour()),
		Minute:   api.SQLUSMALLINT(value.Minute()),
		Second:   api.SQLUSMALLINT(value.Second()),
		Fraction: api.SQLUINTEGER(value.Nanosecond()),
	}
	defer ts.Free()

	sz := unsafe.Sizeof(ts)
	_, err := s.result(api.SQLBindParameter((*api.SQLHSTMT)(s.hnd()),
		api.SQLUSMALLINT(index+1),
		api.SQL_PARAM_INPUT,
		api.SQL_C_TIMESTAMP,
		api.SQL_TYPE_TIMESTAMP,
		api.SQLULEN(sz),
		0,
		(*api.SQLPOINTER)(unsafe.Pointer(&ts)),
		0,
		nil))

	return err
}

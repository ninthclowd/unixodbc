package odbc

import (
	"database/sql/driver"
	"github.com/ninthclowd/unixodbc/internal/api"
	"reflect"
	"time"
)

func init() {
	registerColumnFactoryForType(newTimestampColumn, api.SQL_TYPE_TIMESTAMP)
}

func newTimestampColumn(info *columnInfo, hnd handle) Column {
	return &columnTimestamp{hnd, info}
}

type columnTimestamp struct {
	handle
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
	var valueLength api.SQLLEN
	if _, err := c.result(c.api().SQLGetData(api.SQLHSTMT(c.hnd()), c.columnNumber, api.SQL_C_TIMESTAMP, api.SQLPOINTER(&value), 0, nil)); err != nil {
		return nil, err
	}
	if valueLength == api.SQL_NULL_DATA {
		return nil, nil
	}
	return time.Date(int(value.Year), time.Month(value.Month), int(value.Day), int(value.Hour), int(value.Minute), int(value.Second), int(value.Fraction), time.UTC), nil
}

func (s *Statement) bindTimestamp(index int, value *time.Time) error {
	typeInfo, err := s.conn.TypeInfo(api.SQL_TYPE_TIMESTAMP)
	if err != nil {
		return err
	}
	ts := api.SQL_TIMESTAMP_STRUCT{
		Year:     api.SQLSMALLINT(value.Year()),
		Month:    api.SQLUSMALLINT(value.Month()),
		Day:      api.SQLUSMALLINT(value.Day()),
		Hour:     api.SQLUSMALLINT(value.Hour()),
		Minute:   api.SQLUSMALLINT(value.Minute()),
		Second:   api.SQLUSMALLINT(value.Second()),
		Fraction: api.SQLUINTEGER(value.Nanosecond()),
	}

	_, err = s.result(s.api().SQLBindParameter((api.SQLHSTMT)(s.hnd()), api.SQLUSMALLINT(index+1), api.SQL_PARAM_INPUT,
		api.SQL_C_TIMESTAMP, api.SQL_TYPE_TIMESTAMP,
		typeInfo.ColumnSize, 0, //TODO: validate this is correct for all server types
		api.SQLPOINTER(&ts),
		0, nil))
	return err
}

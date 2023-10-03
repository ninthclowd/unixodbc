package odbc

import (
	"database/sql/driver"
	"github.com/ninthclowd/unixodbc/internal/api"
	"reflect"
	"time"
)

func init() {
	registerColumnFactoryForType(newDateColumn, api.SQL_TYPE_DATE)
}

func newDateColumn(info *columnInfo, hnd handle) Column {
	return &columnDate{hnd, info}
}

type columnDate struct {
	handle
	*columnInfo
}

func (c *columnDate) VariableLength() (length int64, ok bool) {
	return
}

func (c *columnDate) ScanType() reflect.Type {
	return reflect.TypeOf((*time.Time)(nil))
}

func (c *columnDate) Decimal() (precision int64, scale int64, ok bool) {
	return
}

func (c *columnDate) Value() (driver.Value, error) {
	var value api.SQL_DATE_STRUCT
	var valueLength api.SQLLEN
	if _, err := c.result(c.api().SQLGetData(api.SQLHSTMT(c.hnd()), c.columnNumber, api.SQL_C_DATE, api.SQLPOINTER(&value), 0, nil)); err != nil {
		return nil, err
	}
	if valueLength == api.SQL_NULL_DATA {
		return nil, nil
	}
	return time.Date(int(value.Year), time.Month(value.Month), int(value.Day), 0, 0, 0, 0, time.UTC), nil
}

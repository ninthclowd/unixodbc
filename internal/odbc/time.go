package odbc

import (
	"database/sql/driver"
	"github.com/ninthclowd/unixodbc/internal/api"
	"reflect"
	"time"
	"unsafe"
)

func init() {
	registerColumnFactoryForType(newTimeColumn,
		api.SQL_TYPE_TIME,
		-154, //mssql time type
	)
}

func newTimeColumn(info *columnInfo, hnd *handle) Column {
	return &columnTime{hnd, info}
}

type columnTime struct {
	*handle
	*columnInfo
}

func (c *columnTime) VariableLength() (length int64, ok bool) {
	return
}

func (c *columnTime) ScanType() reflect.Type {
	return reflect.TypeOf((*time.Time)(nil))
}

func (c *columnTime) Decimal() (precision int64, scale int64, ok bool) {
	return
}

func (c *columnTime) Value() (driver.Value, error) {
	var value api.SQL_TIME_STRUCT
	defer value.Free()
	var valueLength api.SQLLEN
	if _, err := c.result(api.SQLGetData((*api.SQLHSTMT)(c.hnd()),
		c.columnNumber,
		api.SQL_C_TIME,
		(*api.SQLPOINTER)(unsafe.Pointer(&value)),
		0,
		&valueLength)); err != nil {
		return nil, err
	}
	if valueLength == api.SQL_NULL_DATA {
		return nil, nil
	}
	return time.Date(1, time.January, 1, int(value.Hour), int(value.Minute), int(value.Second), 0, time.UTC), nil
}

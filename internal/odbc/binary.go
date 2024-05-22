package odbc

import (
	"database/sql/driver"
	"github.com/ninthclowd/unixodbc/internal/api"
	"reflect"
	"unsafe"
)

const (
	defaultBinarySize = 10240
)

func init() {
	registerColumnFactoryForType(newBinaryColumn,
		api.SQL_BINARY,
		api.SQL_VARBINARY,
		api.SQL_LONGVARBINARY,
	)
}

func newBinaryColumn(info *columnInfo, hnd *handle) Column {
	return &columnBinary{hnd, info}
}

type columnBinary struct {
	*handle
	*columnInfo
}

func (c *columnBinary) VariableLength() (length int64, ok bool) {
	return int64(c.columnSize), true
}

func (c *columnBinary) ScanType() reflect.Type {
	return reflect.TypeOf((*[]byte)(nil))
}

func (c *columnBinary) Decimal() (precision int64, scale int64, ok bool) {
	return
}

//go:nocheckptr
func (c *columnBinary) Value() (driver.Value, error) {
	size := c.columnSize
	if size == 0 {
		size = defaultBinarySize
	}
	value := make([]byte, size)
	var valueLength api.SQLLEN
	if _, err := c.result(api.SQLGetData((*api.SQLHSTMT)(c.hnd()),
		c.columnNumber,
		api.SQL_C_BINARY,
		(*api.SQLPOINTER)(unsafe.Pointer(&value[0])),
		api.SQLLEN(len(value)),
		&valueLength)); err != nil {
		return nil, err
	}
	if valueLength == api.SQL_NULL_DATA {
		return nil, nil
	}
	return value[:valueLength], nil
}

//go:nocheckptr
func (s *Statement) bindBinary(index int, src []byte) error {
	execSize := api.SQLLEN(len(src))
	_, err := s.result(api.SQLBindParameter((*api.SQLHSTMT)(s.hnd()),
		api.SQLUSMALLINT(index+1),
		api.SQL_PARAM_INPUT,
		api.SQL_C_BINARY,
		api.SQL_BINARY,
		api.SQLULEN(len(src)),
		0,
		(*api.SQLPOINTER)(unsafe.Pointer(&src[0])),
		0,
		&execSize))
	return err
}

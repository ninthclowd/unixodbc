package odbc

import (
	"database/sql/driver"
	"reflect"
	"strings"
	"unsafe"

	"github.com/ninthclowd/unixodbc/internal/api"
)

func init() {
	registerColumnFactoryForType(newUTF8Column,
		api.SQL_CHAR,
		api.SQL_VARCHAR,
		api.SQL_LONGVARCHAR,
	)
}

func newUTF8Column(info *columnInfo, hnd *handle) Column {
	return &columnUTF8{hnd, info}
}

type columnUTF8 struct {
	*handle
	*columnInfo
}

func (c *columnUTF8) VariableLength() (length int64, ok bool) {
	return int64(c.columnSize), true
}

func (c *columnUTF8) ScanType() reflect.Type {
	return reflect.TypeOf((*string)(nil))
}

func (c *columnUTF8) Decimal() (precision int64, scale int64, ok bool) {
	return
}

//go:nocheckptr
func (c *columnUTF8) Value() (driver.Value, error) {
	buffer := make([]uint8, c.columnSize+1) // add 1 for null terminator
	bytesWritten := new(api.SQLLEN)
	defer func() {
		buffer = nil
		bytesWritten = nil
	}()
	maxWrite := api.SQLLEN(len(buffer))
	if _, err := c.result(api.SQLGetData((*api.SQLHSTMT)(c.hnd()),
		c.columnNumber,
		api.SQL_C_CHAR,
		(*api.SQLPOINTER)(unsafe.Pointer(&buffer[0])),
		maxWrite,
		bytesWritten)); err != nil {
		return nil, err
	}
	if *bytesWritten == api.SQL_NULL_DATA {
		return nil, nil
	}
	var builder strings.Builder
	builder.Grow(int(*bytesWritten))
	for i := 0; i < int(*bytesWritten); i++ {
		if buffer[i] == 0 {
			break
		}
		builder.WriteByte(buffer[i])
	}
	return builder.String(), nil
}

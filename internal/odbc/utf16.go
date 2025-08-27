package odbc

import "C"
import (
	"database/sql/driver"
	"encoding/binary"
	"reflect"
	"unicode/utf16"
	"unsafe"

	"github.com/ninthclowd/unixodbc/internal/api"
)

func init() {
	registerColumnFactoryForType(newUTF16Column,
		api.SQL_WCHAR,
		api.SQL_WVARCHAR,
		api.SQL_WLONGVARCHAR,
	)
}

func newUTF16Column(info *columnInfo, hnd *handle) Column {
	return &columnUTF16{hnd, info}
}

type columnUTF16 struct {
	*handle
	*columnInfo
}

func (c *columnUTF16) VariableLength() (length int64, ok bool) {
	return int64(c.columnSize), true
}

func (c *columnUTF16) ScanType() reflect.Type {
	return reflect.TypeOf((*string)(nil))
}

func (c *columnUTF16) Decimal() (precision int64, scale int64, ok bool) {
	return
}

//go:nocheckptr
func (c *columnUTF16) Value() (driver.Value, error) {
	utfLength := c.columnSize * 2
	value := make([]byte, utfLength+1)
	maxWrite := api.SQLLEN(len(value))
	var valueLength api.SQLLEN
	if _, err := c.result(api.SQLGetData((*api.SQLHSTMT)(c.hnd()),
		c.columnNumber,
		api.SQL_C_WCHAR,
		(*api.SQLPOINTER)(unsafe.Pointer(&value[0])),
		maxWrite,
		&valueLength)); err != nil {
		return nil, err
	}
	if valueLength == api.SQL_NULL_DATA || valueLength < 2 {
		return nil, nil
	}
	if valueLength > api.SQLLEN(utfLength) {
		valueLength = api.SQLLEN(utfLength)
	}

	utf := make([]uint16, valueLength/2)
	for i := 0; i < int(valueLength); i += 2 {
		utf[i/2] = binary.LittleEndian.Uint16(value[i : i+2])
	}
	return string(utf16.Decode(utf)), nil
}

//go:nocheckptr
func (s *statement) bindUTF16(index int, src string) error {
	sz := len(src)
	nts := make([]rune, sz+1)
	for i, r := range src {
		nts[i] = r
	}
	val := utf16.Encode(nts)
	_, err := s.result(api.SQLBindParameter((*api.SQLHSTMT)(s.hnd()),
		api.SQLUSMALLINT(index+1),
		api.SQL_PARAM_INPUT,
		api.SQL_C_WCHAR,
		api.SQLSMALLINT(api.SQL_WVARCHAR),
		api.SQLULEN(sz),
		0,
		(*api.SQLPOINTER)(unsafe.Pointer(&val[0])),
		0,
		nil))
	return err
}

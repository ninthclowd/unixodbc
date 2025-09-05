package odbc

import "C"
import (
	"database/sql/driver"
	"reflect"
	"strings"
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
	buffer := make([]uint16, c.columnSize+1)
	maxWrite := api.SQLLEN(len(buffer) * 2)
	bytesWritten := new(api.SQLLEN)
	defer func() {
		buffer = nil
		bytesWritten = nil
	}()
	if _, err := c.result(api.SQLGetData((*api.SQLHSTMT)(c.hnd()),
		c.columnNumber,
		api.SQL_C_WCHAR,
		(*api.SQLPOINTER)(unsafe.Pointer(&buffer[0])),
		maxWrite,
		bytesWritten)); err != nil {
		return nil, err
	}
	if *bytesWritten == api.SQL_NULL_DATA || *bytesWritten < 2 {
		return nil, nil
	}

	runesWritten := int(*bytesWritten / 2)

	// Use strings.Builder with pre-allocated capacity
	var builder strings.Builder
	builder.Grow(runesWritten) // Conservative estimate

	for i := 0; i < runesWritten; i++ {
		r1 := buffer[i]

		// Check for high surrogate
		if 0xD800 <= r1 && r1 <= 0xDBFF && i+1 < runesWritten {
			// Handle surrogate pair
			r2 := buffer[i+1]
			if 0xDC00 <= r2 && r2 <= 0xDFFF {
				// Valid surrogate pair
				r := 0x10000 + (rune(r1&0x3FF) << 10) + rune(r2&0x3FF)
				builder.WriteRune(r)
				i++ // Skip the low surrogate
				continue
			}
		}

		// Regular UTF16 code unit or invalid surrogate
		builder.WriteRune(rune(r1))
	}

	return builder.String(), nil
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

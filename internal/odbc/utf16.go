package odbc

import (
	"database/sql/driver"
	"encoding/binary"
	"github.com/ninthclowd/unixodbc/internal/api"
	"reflect"
	"unicode/utf16"
)

func init() {
	registerColumnFactoryForType(newUTF16Column,
		api.SQL_WCHAR,
		api.SQL_WVARCHAR,
		api.SQL_WLONGVARCHAR,
	)
}

func newUTF16Column(info *columnInfo, hnd handle) Column {
	return &columnUTF16{hnd, info}
}

type columnUTF16 struct {
	handle
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

func (c *columnUTF16) Value() (driver.Value, error) {
	utfLength := c.columnSize * 2
	value := make([]byte, utfLength+1)
	var valueLength api.SQLLEN
	if _, err := c.result(c.api().SQLGetData(api.SQLHSTMT(c.hnd()), c.columnNumber, api.SQL_C_WCHAR, api.SQLPOINTER(&value[0]), api.SQLLEN(len(value)), &valueLength)); err != nil {
		return nil, err
	}
	if valueLength == api.SQL_NULL_DATA {
		return nil, nil
	}

	str := utf16String(value[:valueLength])
	return str, nil
}

func (s *Statement) bindUTF16(index int, src string) error {
	size := len(src)
	nts := make([]rune, len(src)+1)
	for i, r := range src {
		nts[i] = r
	}
	val := utf16.Encode(nts)
	dType, dSize, err := s.conn.stringDataType(size)
	if err != nil {
		return err
	}
	execSize := api.SQLLEN(size * 2)
	_, err = s.result(s.api().SQLBindParameter((api.SQLHSTMT)(s.hnd()), api.SQLUSMALLINT(index+1), api.SQL_PARAM_INPUT,
		api.SQL_C_WCHAR, api.SQLSMALLINT(dType),
		dSize, 0,
		api.SQLPOINTER(&val[0]),
		0, &execSize))
	return err
}

func utf16String(b []byte) string {
	utf := make([]uint16, len(b)/2)
	for i := 0; i < len(b); i += 2 {
		utf[i/2] = binary.LittleEndian.Uint16(b[i:])
	}
	return string(utf16.Decode(utf))
}

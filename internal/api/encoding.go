package api

// #include <sql.h>
// #include <sqlext.h>
import "C"
import (
	"errors"
	"fmt"
	"reflect"
	"time"
	"unicode/utf16"
	"unsafe"
)

var (
	reflectTypeByte     = reflect.TypeOf((*byte)(nil))
	reflectTypeString   = reflect.TypeOf((*string)(nil))
	reflectTypeFloat32  = reflect.TypeOf((*float32)(nil))
	reflectTypeFloat64  = reflect.TypeOf((*float64)(nil))
	reflectTypeBool     = reflect.TypeOf((*bool)(nil))
	reflectTypeInt8     = reflect.TypeOf((*int8)(nil))
	reflectTypeInt16    = reflect.TypeOf((*int16)(nil))
	reflectTypeInt32    = reflect.TypeOf((*int32)(nil))
	reflectTypeInt64    = reflect.TypeOf((*int64)(nil))
	reflectTypeTime     = reflect.TypeOf((*time.Time)(nil))
	reflectTypeBytes    = reflect.TypeOf(([]byte)(nil))
	reflectTypeDuration = reflect.TypeOf((*time.Duration)(nil))
)

var registeredEncodings = map[SQLSMALLINT]Encoding{}

func registerEncoding(decoders ...Encoding) {
	for _, decoder := range decoders {
		for _, sqlType := range decoder.SQLTypes() {
			registeredEncodings[sqlType] = decoder
		}
	}
}

func init() {
	registerEncoding(
		&boolEncoding{},
		&utf8Encoding{},
		&utf16Encoding{},
		&int8Encoding{},
		&int16Encoding{},
		&int32Encoding{},
		&int64Encoding{},
		&float32Encoding{},
		&float64Encoding{},
		&binaryEncoding{},
		&timestampEncoding{},
		&dateEncoding{},
		&timeEncoding{},
	)

}

type Encoding interface {
	Encoder
	Decoder
	SQLTypes() []SQLSMALLINT
}

type boolEncoding struct{}

func (d *boolEncoding) Encode(p *Param, value interface{}) error {
	return errors.New("not implemented")
}

func (d *boolEncoding) ScanType() reflect.Type {
	return reflectTypeBool
}

func (d *boolEncoding) SQLTypes() []SQLSMALLINT {
	return []SQLSMALLINT{SQL_BIT}
}
func (d *boolEncoding) Decode(c *Column) (out interface{}, err error) {
	var value SQLCHAR
	code, err := c.stmt.Result(C.SQLGetData(c.stmt.hPtr,
		C.SQLUSMALLINT(c.pos),
		C.SQL_C_BIT,
		C.SQLPOINTER(&value),
		C.SQLLEN(0),
		nil))
	var b bool
	out = b
	if err == nil && code == SQL_NULL_DATA && value == 1 {
		out = true
	}
	return
}

type utf8Encoding struct{}

func (d *utf8Encoding) SQLTypes() []SQLSMALLINT {
	return []SQLSMALLINT{SQL_CHAR, SQL_LONGVARCHAR, SQL_VARCHAR}
}
func (d *utf8Encoding) ScanType() reflect.Type {
	return reflectTypeString
}
func (d *utf8Encoding) Decode(c *Column) (out interface{}, err error) {
	value := CHAR(make([]uint8, c.columnSize))
	var valueLength SQLLEN
	code, err := c.stmt.Result(C.SQLGetData(c.stmt.hPtr, C.SQLUSMALLINT(c.pos),
		C.SQL_C_CHAR,
		C.SQLPOINTER(value.Address()),
		C.SQLLEN(len(value)),
		(*C.SQLLEN)(&valueLength)))
	if err == nil && code != SQL_NULL_DATA {
		out = value[:valueLength-1].String()
	}
	return
}

func (d *utf8Encoding) Encode(p *Param, val interface{}) error {
	value := CHAR(*(val.(*string)))
	_, err := p.stmt.Result(C.SQLBindParameter(p.stmt.hPtr,
		C.SQLUSMALLINT(p.pos),
		C.SQL_PARAM_INPUT_OUTPUT,
		C.SQL_C_CHAR,
		C.SQLSMALLINT(p.dataType),
		C.SQLULEN(p.parameterSize),
		C.SQLSMALLINT(0),
		C.SQLPOINTER(value.Address()),
		C.SQLLEN(len(value)),
		nil))
	return err
}

type utf16Encoding struct{}

func (d *utf16Encoding) Encode(p *Param, val interface{}) error {
	value := WCHAR(utf16.Encode([]rune(*(val.(*string)))))
	_, err := p.stmt.Result(C.SQLBindParameter(p.stmt.hPtr,
		C.SQLUSMALLINT(p.pos),
		C.SQL_PARAM_INPUT_OUTPUT,
		C.SQL_C_WCHAR,
		C.SQLSMALLINT(p.dataType),
		C.SQLULEN(p.parameterSize),
		C.SQLSMALLINT(0),
		C.SQLPOINTER(value.Address()),
		C.SQLLEN(len(value)),
		nil))
	return err
}

func (d *utf16Encoding) SQLTypes() []SQLSMALLINT {
	return []SQLSMALLINT{SQL_WCHAR, SQL_WLONGVARCHAR, SQL_WVARCHAR}
}

func (d *utf16Encoding) ScanType() reflect.Type {
	return reflectTypeString
}
func (d *utf16Encoding) Decode(c *Column) (out interface{}, err error) {
	value := WCHAR(make([]uint16, c.columnSize))
	var valueLength SQLLEN
	code, err := c.stmt.Result(C.SQLGetData(c.stmt.hPtr, C.SQLUSMALLINT(c.pos),
		C.SQL_C_WCHAR,
		C.SQLPOINTER(unsafe.Pointer(value.Address())),
		C.SQLLEN(len(value)),
		(*C.SQLLEN)(unsafe.Pointer(&valueLength))))
	if err == nil && code != SQL_NULL_DATA {
		out = value[:valueLength/2].String()
	}
	return
}

type int8Encoding struct{}

func (d *int8Encoding) Encode(p *Param, value interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (d *int8Encoding) SQLTypes() []SQLSMALLINT {
	return []SQLSMALLINT{SQL_TINYINT}
}

func (d *int8Encoding) ScanType() reflect.Type {
	return reflectTypeInt8
}

func (d *int8Encoding) Decode(c *Column) (out interface{}, err error) {
	var value SQLSCHAR

	code, err := c.stmt.Result(C.SQLGetData(c.stmt.hPtr,
		C.SQLUSMALLINT(c.pos),
		C.SQL_C_STINYINT,
		C.SQLPOINTER(&value),
		C.SQLLEN(0),
		nil))

	if err == nil && code != SQL_NULL_DATA {
		out = int8(value)
	}
	return
}

type int16Encoding struct{}

func (d *int16Encoding) Encode(p *Param, value interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (d *int16Encoding) SQLTypes() []SQLSMALLINT {
	return []SQLSMALLINT{SQL_SMALLINT}
}

func (d *int16Encoding) ScanType() reflect.Type {
	return reflectTypeInt16
}

func (d *int16Encoding) Decode(c *Column) (out interface{}, err error) {

	var value SQLSMALLINT

	code, err := c.stmt.Result(C.SQLGetData(c.stmt.hPtr,
		C.SQLUSMALLINT(c.pos),
		C.SQL_C_SSHORT,
		C.SQLPOINTER(&value),
		C.SQLLEN(0),
		nil))

	if err == nil && code != SQL_NULL_DATA {
		out = int16(value)
	}
	return
}

type int32Encoding struct{}

func (d *int32Encoding) Encode(p *Param, value interface{}) error {
	var val = SQLINTEGER(*value.(*int32))
	_, err := p.stmt.Result(C.SQLBindParameter(p.stmt.hPtr,
		C.SQLUSMALLINT(p.pos),
		C.SQL_PARAM_INPUT_OUTPUT,
		C.SQL_C_SLONG,
		C.SQLSMALLINT(p.dataType),
		C.SQLULEN(p.parameterSize),
		C.SQLSMALLINT(p.decimalDigits),
		C.SQLPOINTER(&val),
		C.SQLLEN(0),
		nil))
	return err
}

func (d *int32Encoding) SQLTypes() []SQLSMALLINT {
	return []SQLSMALLINT{SQL_INTEGER}
}

func (d *int32Encoding) ScanType() reflect.Type {
	return reflectTypeInt32
}

func (d *int32Encoding) Decode(c *Column) (out interface{}, err error) {
	var val SQLINTEGER
	code, err := c.stmt.Result(C.SQLGetData(
		c.stmt.hPtr,
		C.SQLUSMALLINT(c.pos),
		C.SQL_C_SLONG,
		C.SQLPOINTER(&val),
		C.SQLLEN(0),
		nil,
	))
	if err == nil && code != SQL_NULL_DATA {
		out = int32(val)
	}
	return
}

type int64Encoding struct{}

func (d *int64Encoding) Encode(p *Param, value interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (d *int64Encoding) SQLTypes() []SQLSMALLINT {
	return []SQLSMALLINT{SQL_BIGINT, SQL_NUMERIC}
}

func (d *int64Encoding) ScanType() reflect.Type {
	return reflectTypeInt64
}

func (d *int64Encoding) Decode(c *Column) (out interface{}, err error) {

	var value SQLBIGINT
	code, err := c.stmt.Result(C.SQLGetData(c.stmt.hPtr,
		C.SQLUSMALLINT(c.pos),
		C.SQL_C_SBIGINT,
		C.SQLPOINTER(&value),
		C.SQLLEN(0),
		nil))
	if err == nil && code != SQL_NULL_DATA {
		out = int64(value)
	}
	return
}

type float32Encoding struct{}

func (d *float32Encoding) Encode(p *Param, value interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (d *float32Encoding) SQLTypes() []SQLSMALLINT {
	return []SQLSMALLINT{SQL_REAL}
}

func (d *float32Encoding) ScanType() reflect.Type {
	return reflectTypeFloat32
}

func (d *float32Encoding) Decode(c *Column) (out interface{}, err error) {

	var value SQLREAL
	code, err := c.stmt.Result(C.SQLGetData(c.stmt.hPtr,
		C.SQLUSMALLINT(c.pos),
		C.SQL_C_FLOAT,
		C.SQLPOINTER(&value),
		C.SQLLEN(0),
		nil))
	if err == nil && code != SQL_NULL_DATA {
		out = float32(value)
	}
	return
}

type float64Encoding struct{}

func (d *float64Encoding) Encode(p *Param, value interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (d *float64Encoding) SQLTypes() []SQLSMALLINT {
	return []SQLSMALLINT{SQL_DOUBLE, SQL_FLOAT, SQL_DECIMAL}
}

func (d *float64Encoding) ScanType() reflect.Type {
	return reflectTypeFloat64
}

func (d *float64Encoding) Decode(c *Column) (out interface{}, err error) {
	var value SQLFLOAT
	code, err := c.stmt.Result(C.SQLGetData(c.stmt.hPtr,
		C.SQLUSMALLINT(c.pos),
		C.SQL_C_DOUBLE,
		C.SQLPOINTER(&value),
		C.SQLLEN(0),
		nil))
	if err == nil && code != SQL_NULL_DATA {
		out = float64(value)
	}
	return
}

type binaryEncoding struct{}

func (d *binaryEncoding) Encode(p *Param, value interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (d *binaryEncoding) SQLTypes() []SQLSMALLINT {
	return []SQLSMALLINT{SQL_BINARY, SQL_VARBINARY, SQL_LONGVARBINARY}
}

func (d *binaryEncoding) ScanType() reflect.Type {
	return reflectTypeBytes
}

func (d *binaryEncoding) Decode(c *Column) (out interface{}, err error) {
	value := CHAR(make([]byte, c.columnSize))
	var valueLength SQLLEN
	code, err := c.stmt.Result(C.SQLGetData(c.stmt.hPtr, C.SQLUSMALLINT(c.pos),
		C.SQL_C_CHAR,
		C.SQLPOINTER(value.Address()),
		C.SQLLEN(len(value)),
		(*C.SQLLEN)(&valueLength)))
	if err == nil && code != SQL_NULL_DATA {
		out = value[:valueLength]
	}
	return
}

type timestampEncoding struct{}

func (d *timestampEncoding) Encode(p *Param, value interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (d *timestampEncoding) SQLTypes() []SQLSMALLINT {
	return []SQLSMALLINT{SQL_TYPE_TIMESTAMP}
}

func (d *timestampEncoding) ScanType() reflect.Type {
	return reflectTypeTime
}

func (d *timestampEncoding) Decode(c *Column) (out interface{}, err error) {
	var value SQL_TIMESTAMP_STRUCT
	defer value.Free() //TODO
	code, err := c.stmt.Result(C.SQLGetData(c.stmt.hPtr, C.SQLUSMALLINT(c.pos),
		C.SQL_C_TYPE_TIMESTAMP,
		C.SQLPOINTER(unsafe.Pointer(&value)),
		C.SQLLEN(0),
		nil))
	if err == nil && code != SQL_NULL_DATA {
		out = time.Date(int(value.Year), time.Month(value.Month), int(value.Day), int(value.Hour), int(value.Minute), int(value.Second), int(value.Fraction), time.UTC)
	}
	return
}

type dateEncoding struct{}

func (d *dateEncoding) Encode(p *Param, value interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (d *dateEncoding) SQLTypes() []SQLSMALLINT {
	return []SQLSMALLINT{SQL_TYPE_DATE}
}

func (d *dateEncoding) ScanType() reflect.Type {
	return reflectTypeTime
}

func (d *dateEncoding) Decode(c *Column) (out interface{}, err error) {
	var value SQL_DATE_STRUCT
	defer value.Free() //TODO
	code, err := c.stmt.Result(C.SQLGetData(c.stmt.hPtr, C.SQLUSMALLINT(c.pos),
		C.SQL_C_TYPE_TIMESTAMP,
		C.SQLPOINTER(unsafe.Pointer(&value)),
		C.SQLLEN(0),
		nil))
	if err == nil && code != SQL_NULL_DATA {
		out = time.Date(int(value.Year), time.Month(value.Month), int(value.Day), 0, 0, 0, 0, time.UTC)
	}
	return
}

type timeEncoding struct{}

func (d *timeEncoding) Encode(p *Param, value interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (d *timeEncoding) SQLTypes() []SQLSMALLINT {
	return []SQLSMALLINT{SQL_TYPE_TIME}
}

func (d *timeEncoding) ScanType() reflect.Type {
	return reflectTypeDuration
}

func (d *timeEncoding) Decode(c *Column) (out interface{}, err error) {
	var value SQL_TIME_STRUCT
	defer value.Free() //TODO
	code, err := c.stmt.Result(C.SQLGetData(c.stmt.hPtr, C.SQLUSMALLINT(c.pos),
		C.SQL_C_TYPE_TIMESTAMP,
		C.SQLPOINTER(unsafe.Pointer(&value)),
		C.SQLLEN(0),
		nil))
	if err == nil && code != SQL_NULL_DATA {
		out, _ = time.ParseDuration(fmt.Sprintf("%dh%dm%ds", int(value.Hour), int(value.Minute), int(value.Second)))
	}
	return
}

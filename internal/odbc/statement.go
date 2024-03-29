package odbc

import (
	"context"
	"errors"
	"fmt"
	"github.com/ninthclowd/unixodbc/internal/api"

	"math"
	"reflect"
	"time"
	"unicode/utf16"
)

var (
	ErrRecordSetOpen = errors.New("recordset is still open")
)

type CursorSensitivity uint64

const (
	CursorSensitive   = CursorSensitivity(api.SQL_SENSITIVE)
	CursorInsensitive = CursorSensitivity(api.SQL_INSENSITIVE)
)

type Concurrency uint64

const (
	ConcurrencyLock = Concurrency(api.SQL_CONCUR_LOCK)
)

type Statement struct {
	handle
	conn *Connection
	rs   *RecordSet
}

func (s *Statement) SetCursorSensitivity(sensitivity CursorSensitivity) error {
	_, err := s.result(s.api().SQLSetStmtAttrConst(api.SQLHSTMT(s.hnd()), api.SQL_ATTR_CURSOR_SENSITIVITY, uint64(sensitivity)))
	return err
}

func (s *Statement) SetConcurrency(concurrency Concurrency) error {
	_, err := s.result(s.api().SQLSetStmtAttrConst(api.SQLHSTMT(s.hnd()), api.SQL_ATTR_CONCURRENCY, uint64(concurrency)))
	return err
}

func (s *Statement) NumParams() (int, error) {
	var paramCount api.SQLSMALLINT
	if _, err := s.result(s.api().SQLNumParams((api.SQLHSTMT)(s.hnd()), &paramCount)); err != nil {
		return 0, fmt.Errorf("getting Parameter count: %w", err)
	}
	return int(paramCount), nil
}

func (s *Statement) ResetParams() error {
	_, err := s.result(s.api().SQLFreeStmt((api.SQLHSTMT)(s.hnd()), api.SQL_RESET_PARAMS))
	return err
}

func (s *Statement) Close() error {
	if _, err := s.result(s.api().SQLFreeStmt((api.SQLHSTMT)(s.hnd()), api.SQL_CLOSE)); err != nil {
		return fmt.Errorf("freeing Statement: %w", err)
	}
	return s.free()
}

func (s *Statement) ExecDirect(ctx context.Context, sql string) error {
	done := cancelHandleOnContext(ctx, s)

	statementBytes := utf16.Encode([]rune(sql))
	_, err := s.result(s.api().SQLExecDirect(api.SQLHSTMT(s.hnd()), statementBytes, api.SQLINTEGER(len(statementBytes))))
	done()
	if err == nil {
		err = ctx.Err()
	}
	return err
}

func (s *Statement) Execute(ctx context.Context) error {
	done := cancelHandleOnContext(ctx, s)
	_, err := s.result(s.api().SQLExecute(api.SQLHSTMT(s.hnd())))
	done()
	if err == nil {
		err = ctx.Err()
	}
	return err
}

func (s *Statement) Prepare(ctx context.Context, sql string) error {
	done := cancelHandleOnContext(ctx, s)

	statementBytes := utf16.Encode([]rune(sql))
	_, err := s.result(s.api().SQLPrepare(api.SQLHSTMT(s.hnd()), statementBytes, api.SQLINTEGER(len(statementBytes))))
	done()
	if err == nil {
		err = ctx.Err()
	}
	return err
}

func (s *Statement) RecordSet() (*RecordSet, error) {
	if s.rs != nil {
		return nil, ErrRecordSetOpen
	}
	col, err := columnsForStatement(s.handle, newColumnLoader(s.handle))
	if err != nil {
		return nil, err
	}

	s.rs = &RecordSet{stmt: s, columns: col}
	return s.rs, nil
}

func (s *Statement) closeCursor() error {
	if _, err := s.result(s.api().SQLCloseCursor((api.SQLHSTMT)(s.hnd()))); err != nil {
		return fmt.Errorf("closing cursor: %w", err)
	}
	s.rs = nil
	return nil
}

func (s *Statement) fetch() (more bool, err error) {
	if code, err := s.result(s.api().SQLFetch((api.SQLHSTMT)(s.hnd()))); err != nil {
		return false, err
	} else if code == api.SQL_NO_DATA {
		return false, nil
	}
	return true, nil
}

func (s *Statement) BindParams(params ...interface{}) error {
	for i, param := range params {
		if err := s.bindParam(i, param); err != nil {
			return err
		}
	}
	return nil
}

func (s *Statement) bindParam(index int, value interface{}) error {
	value = compressInt(value) //TODO, benchmark don't switch on non int types
	switch value.(type) {
	case nil:
		return s.bindNil(index)
	case bool:
		return s.bindBool(index, value.(bool))
	case string:
		return s.bindUTF16(index, value.(string))
	case int8:
		return s.bindInt8(index, value.(int8))
	case int16:
		return s.bindInt16(index, value.(int16))
	case int32:
		return s.bindInt32(index, value.(int32))
	case int64:
		return s.bindInt64(index, value.(int64))
	case float64:
		return s.bindFloat64(index, value.(float64))
	case []byte:
		return s.bindBinary(index, value.([]byte))
	case time.Time:
		v := value.(time.Time)
		return s.bindTimestamp(index, &v)
	default:
		return fmt.Errorf("unable to bind parameter of type %s", reflect.TypeOf(value).String())
	}
}

func (s *Statement) bindNil(index int) error {
	strLenOrIndPtr := api.SQLLEN(api.SQL_NULL_DATA)
	_, err := s.result(s.api().SQLBindParameter((api.SQLHSTMT)(s.hnd()), api.SQLUSMALLINT(index+1), api.SQL_PARAM_INPUT,
		api.SQL_C_CHAR, api.SQL_CHAR,
		1, 0,
		nil,
		0, &strLenOrIndPtr))
	return err
}

type columnsDetails struct {
	names   []string
	byName  map[string]Column
	byIndex []Column
}

func compressInt(val interface{}) interface{} {
	switch val.(type) {
	case int8:
		return val
	case int16:
		if v := val.(int16); v <= math.MaxInt8 {
			return int8(v)
		}
		return val
	case int32:
		if v := val.(int32); v <= math.MaxInt8 {
			return int8(v)
		} else if v < math.MaxInt16 {
			return int16(v)
		}
		return val
	case int64:
		if v := val.(int64); v <= math.MaxInt8 {
			return int8(v)
		} else if v <= math.MaxInt16 {
			return int16(v)
		} else if v <= math.MaxInt32 {
			return int32(v)
		}
		return val
	default:
		return val
	}
}

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

//go:generate mockgen -source=statement.go -package mocks -destination ../mocks/statement.go
type Statement interface {
	SetCursorSensitivity(sensitivity CursorSensitivity) error
	SetConcurrency(concurrency Concurrency) error
	NumParams() (int, error)
	ResetParams() error
	Close() error
	ExecDirect(ctx context.Context, sql string) error
	Execute(ctx context.Context) error
	Prepare(ctx context.Context, sql string) error
	BindParams(params ...interface{}) error
	RecordSet() (RecordSet, error)
}

var _ Statement = (*statement)(nil)

type statement struct {
	*handle
	conn *connection
	rs   *recordSet
}

func (s *statement) SetCursorSensitivity(sensitivity CursorSensitivity) error {
	_, err := s.result(api.SQLSetStmtAttr((*api.SQLHSTMT)(s.hnd()),
		api.SQL_ATTR_CURSOR_SENSITIVITY,
		api.Const(uint64(sensitivity)),
		api.SQL_IS_UINTEGER))
	return err
}

func (s *statement) SetConcurrency(concurrency Concurrency) error {
	_, err := s.result(api.SQLSetStmtAttr((*api.SQLHSTMT)(s.hnd()),
		api.SQL_ATTR_CONCURRENCY,
		api.Const(uint64(concurrency)),
		api.SQL_IS_UINTEGER))
	return err
}

func (s *statement) NumParams() (int, error) {
	var paramCount api.SQLSMALLINT
	if _, err := s.result(api.SQLNumParams((*api.SQLHSTMT)(s.hnd()), &paramCount)); err != nil {
		return 0, fmt.Errorf("getting Parameter count: %w", err)
	}
	return int(paramCount), nil
}

func (s *statement) ResetParams() error {
	_, err := s.result(api.SQLFreeStmt((*api.SQLHSTMT)(s.hnd()), api.SQL_RESET_PARAMS))
	return err
}

func (s *statement) Close() error {
	if _, err := s.result(api.SQLFreeStmt((*api.SQLHSTMT)(s.hnd()), api.SQL_CLOSE)); err != nil {
		return fmt.Errorf("freeing statement: %w", err)
	}
	return s.free()
}

func (s *statement) ExecDirect(ctx context.Context, sql string) error {
	done := cancelHandleOnContext(ctx, s.handle)

	statementBytes := utf16.Encode([]rune(sql))
	_, err := s.result(api.SQLExecDirectW((*api.SQLHSTMT)(s.hnd()),
		(*api.SQLWCHAR)(&statementBytes[0]),
		api.SQLINTEGER(len(statementBytes))))
	done()
	if err == nil {
		err = ctx.Err()
	}
	return err
}

func (s *statement) Execute(ctx context.Context) error {
	done := cancelHandleOnContext(ctx, s.handle)
	_, err := s.result(api.SQLExecute((*api.SQLHSTMT)(s.hnd())))
	done()
	if err == nil {
		err = ctx.Err()
	}
	return err
}

func (s *statement) Prepare(ctx context.Context, sql string) error {
	done := cancelHandleOnContext(ctx, s.handle)

	statementBytes := utf16.Encode([]rune(sql))
	_, err := s.result(api.SQLPrepareW((*api.SQLHSTMT)(s.hnd()),
		(*api.SQLWCHAR)(&statementBytes[0]),
		api.SQLINTEGER(len(statementBytes))))
	done()
	if err == nil {
		err = ctx.Err()
	}
	return err
}

func (s *statement) RecordSet() (RecordSet, error) {
	if s.rs != nil {
		return nil, ErrRecordSetOpen
	}
	col, err := columnsForStatement(s.handle, newColumnLoader(s.handle))
	if err != nil {
		return nil, err
	}

	s.rs = &recordSet{stmt: s, columns: col}
	return s.rs, nil
}

func (s *statement) closeCursor() error {
	if _, err := s.result(api.SQLCloseCursor((*api.SQLHSTMT)(s.hnd()))); err != nil {
		return fmt.Errorf("closing cursor: %w", err)
	}
	s.rs = nil
	return nil
}

func (s *statement) fetch() (more bool, err error) {
	if code, err := s.result(api.SQLFetch((*api.SQLHSTMT)(s.hnd()))); err != nil {
		return false, err
	} else if code == api.SQL_NO_DATA {
		return false, nil
	}
	return true, nil
}

func (s *statement) BindParams(params ...interface{}) error {
	for i, param := range params {
		if err := s.bindParam(i, param); err != nil {
			return err
		}
	}
	return nil
}

func (s *statement) bindParam(index int, value interface{}) error {
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

func (s *statement) bindNil(index int) error {
	strLenOrIndPtr := api.SQLLEN(api.SQL_NULL_DATA)
	_, err := s.result(api.SQLBindParameter((*api.SQLHSTMT)(s.hnd()),
		api.SQLUSMALLINT(index+1),
		api.SQL_PARAM_INPUT,
		api.SQL_C_CHAR,
		api.SQL_CHAR,
		1,
		0,
		nil,
		0,
		&strLenOrIndPtr))
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

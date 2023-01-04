package api

// #include <sql.h>
// #include <sqlext.h>
// #include <stdint.h>
import "C"
import (
	"fmt"
	"sync/atomic"
	"unicode/utf16"
	"unsafe"
)

type Statement struct {
	*Handle
	conn       *Connection
	params     Params
	columns    Columns
	queryCount atomic.Int64
}

func (stmt *Statement) Params() (Params, error) {
	if stmt.params != nil {
		return stmt.params, nil
	}
	num, err := stmt.sqlNumParams()
	if err != nil {
		return nil, err
	}

	params := make(Params, num)

	for i := 0; i < num; i++ {
		params[i], err = stmt.sqlDescribeParam(i)
		if err != nil {
			return nil, err

		}
	}
	stmt.params = params
	return stmt.params, nil
}

func (stmt *Statement) Cancel() error {
	_, err := stmt.Result(C.SQLCancel(stmt.hPtr))
	return err
}

func (stmt *Statement) Columns() (Columns, error) {
	if stmt.columns != nil {
		return stmt.columns, nil
	}
	num, err := stmt.sqlNumResultCols()
	if err != nil {
		return nil, err
	}

	columns := make(Columns, num)

	for i := 0; i < num; i++ {
		columns[i], err = stmt.sqlDescribeCol(i)
		if err != nil {
			return nil, err

		}
	}
	stmt.columns = columns
	return stmt.columns, nil
}

func (stmt *Statement) Free() error {
	if _, err := stmt.Result(C.SQLFreeStmt(stmt.hPtr, C.SQL_CLOSE)); err != nil {
		return fmt.Errorf("freeing statement: %w", err)
	}
	return nil
}

func (stmt *Statement) sqlDescribeParam(idx int) (*Param, error) {
	parameterNumber := SQLUSMALLINT(idx + 1)
	var dataType, decimalDigits, nullable SQLSMALLINT
	var parameterSize SQLULEN

	if _, err := stmt.Result(C.SQLDescribeParam(stmt.hPtr, C.SQLUSMALLINT(parameterNumber),
		(*C.SQLSMALLINT)(&dataType), (*C.SQLULEN)(&parameterSize),
		(*C.SQLSMALLINT)(&decimalDigits), (*C.SQLSMALLINT)(&nullable))); err != nil {
		return nil, fmt.Errorf("describing parameter: %w", err)
	}

	n := false
	if nullable == C.SQL_NULLABLE {
		n = true
	}

	encoder, found := registeredEncodings[dataType]
	if !found {
		return nil, fmt.Errorf("encoder not registered for SQL data type %d", dataType)
	}
	return &Param{
		encoder:       encoder,
		stmt:          stmt,
		pos:           int(parameterNumber),
		dataType:      dataType,
		parameterSize: parameterSize,
		decimalDigits: decimalDigits,
		nullable:      n,
	}, nil

}

func (stmt *Statement) SQLGetTypeInfo(dataType SQLSMALLINT) error {
	stmt.queryCount.Add(1)
	_, err := stmt.Result(C.SQLGetTypeInfo(stmt.hPtr, C.SQLSMALLINT(dataType)))
	return err
}

func (stmt *Statement) sqlNumParams() (int, error) {
	num := new(SQLSMALLINT)
	if _, err := stmt.Result(C.SQLNumParams(stmt.hPtr, (*C.SQLSMALLINT)(num))); err != nil {
		return 0, fmt.Errorf("getting parameter count: %w", err)
	}
	return (int)(*num), nil
}

func (stmt *Statement) SQLPrepare(query string) (err error) {
	stmt.params = nil
	stmt.columns = nil
	str := WCHAR(utf16.Encode([]rune(query)))
	if _, err = stmt.Result(C.SQLPrepareW(stmt.hPtr,
		(*C.SQLWCHAR)(unsafe.Pointer(str.Address())),
		C.SQLINTEGER(len(str)))); err != nil {
		return fmt.Errorf("preparing statement: %w", err)
	}
	stmt.params, err = stmt.Params()
	return
}

func (stmt *Statement) SQLExecute() (err error) {
	stmt.queryCount.Add(1)
	stmt.columns = nil
	if _, err := stmt.Result(C.SQLExecute(stmt.hPtr)); err != nil {
		return fmt.Errorf("executing statement: %w", err)
	}

	stmt.columns, err = stmt.Columns()
	//todo SQL_SUCCESS, SQL_SUCCESS_WITH_INFO, SQL_NEED_DATA, SQL_STILL_EXECUTING, SQL_ERROR, SQL_NO_DATA, SQL_INVALID_HANDLE, or SQL_PARAM_DATA_AVAILABLE.
	return
}

func (stmt *Statement) SQLCloseCursor() error {
	if _, err := stmt.Result(C.SQLCloseCursor(stmt.hPtr)); err != nil {
		return fmt.Errorf("closing cursor: %w", err)
	}
	return nil
}

func (stmt *Statement) SQLFetch() (bool, error) {
	if code, err := stmt.Result(C.SQLFetch(stmt.hPtr)); err != nil {
		return false, fmt.Errorf("fetching: %w", err)
	} else if code == SQL_NO_DATA {
		return false, nil
	}
	return true, nil
}

func (stmt *Statement) sqlNumResultCols() (int, error) {
	var columnCount C.SQLSMALLINT
	if _, err := stmt.Result(C.SQLNumResultCols(stmt.hPtr, (*C.SQLSMALLINT)(&columnCount))); err != nil {
		return 0, fmt.Errorf("getting column count: %w", err)
	}
	return int(columnCount), nil
}

func (stmt *Statement) sqlDescribeCol(idx int) (*Column, error) {
	columnNumber := SQLUSMALLINT(idx + 1)
	columnName := CHAR(make([]uint8, 100))

	var nameLength, dataType, decimalDigits, nullInfo SQLSMALLINT
	var columnSize SQLULEN
	_, err := stmt.Result(C.SQLDescribeCol(stmt.hPtr,
		C.SQLUSMALLINT(columnNumber),
		(*C.SQLCHAR)(columnName.Address()),
		C.SQLSMALLINT(len(columnName)),
		(*C.SQLSMALLINT)(&nameLength),
		(*C.SQLSMALLINT)(&dataType),
		(*C.SQLULEN)(&columnSize),
		(*C.SQLSMALLINT)(&decimalDigits),
		(*C.SQLSMALLINT)(&nullInfo)))
	if err != nil {
		return nil, fmt.Errorf("describing column: %w", err)
	}

	var nullable, nullableKnown bool
	switch nullInfo {
	case C.SQL_NULLABLE:
		nullable = true
		nullableKnown = true
	case C.SQL_NO_NULLS:
		nullable = false
		nullableKnown = true
	}

	decoder, found := registeredEncodings[dataType]
	if !found {
		return nil, fmt.Errorf("no decoder registered for SQLDataType %d", dataType)
	}

	return &Column{
		decoder:       decoder,
		pos:           idx + 1,
		stmt:          stmt,
		name:          columnName[:nameLength].String(),
		dataType:      dataType,
		decimalDigits: int64(decimalDigits),
		nullable:      nullable,
		nullableKnown: nullableKnown,
		columnSize:    int64(columnSize),
	}, nil
}

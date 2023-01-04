package api

// #include <sql.h>
// #include <sqlext.h>
// #include <stdint.h>
import "C"
import (
	"errors"
	"fmt"
	"sync"
	"unicode/utf16"
)

type TypeInfo struct {
	TypeName string
}

type Connection struct {
	env *Environment
	mux sync.Mutex
	*Handle
	closed   bool
	typeInfo map[SQLSMALLINT]*TypeInfo
}

func (conn *Connection) TypeInfo(dataType SQLSMALLINT) (*TypeInfo, error) {
	conn.mux.Lock()
	defer conn.mux.Unlock()

	if ti, ok := conn.typeInfo[dataType]; ok {
		return ti, nil
	}

	stmt, err := conn.Statement()
	if err != nil {
		return nil, err
	}
	defer func() {
		stmt.SQLCloseCursor()
		stmt.Free()
	}()
	if err = stmt.SQLGetTypeInfo(dataType); err != nil {
		return nil, err
	}

	columns, err := stmt.Columns()
	if err != nil {
		return nil, err
	}
	columnForField := columns.Map()

	if more, err := stmt.SQLFetch(); err != nil {
		return nil, err
	} else if !more {
		return nil, errors.New("unable to get TypeInfo")
	}
	ti := new(TypeInfo)
	conn.typeInfo[dataType] = ti

	if val, err := columnForField["TYPE_NAME"].Decode(); err != nil {
		return nil, err
	} else {
		ti.TypeName = val.(string)
	}

	return ti, nil
}

func (conn *Connection) SQLDisconnect() error {
	if conn.closed {
		return ErrClosed
	}
	if _, err := conn.Result(C.SQLDisconnect(conn.hPtr)); err != nil {
		return fmt.Errorf("disconnecting: %w", err)
	}
	conn.closed = true
	return nil
}

func (conn *Connection) DriverConnect(connStr string) error {
	str := WCHAR(utf16.Encode([]rune(connStr)))
	if _, err := conn.Result(C.SQLDriverConnectW(conn.hPtr, nil, (*C.SQLWCHAR)(str.Address()), C.SQLSMALLINT(len(str)), nil, 0, nil, C.SQL_DRIVER_NOPROMPT)); err != nil {
		return fmt.Errorf("connecting: %w", err)
	}
	return nil
}

func (conn *Connection) Statement() (*Statement, error) {
	if conn.closed {
		return nil, ErrClosed
	}
	hnd := new(SQLHANDLE)
	code, err := conn.Result(C.SQLAllocHandle(C.SQL_HANDLE_STMT, conn.hPtr, (*C.SQLHANDLE)(hnd)))
	if err != nil {
		if code == SQL_ERROR {
			defer conn.Free()
		}
		return nil, fmt.Errorf("alloc statement handle: %w", err)
	}

	stmt := &Statement{
		conn: conn,
		Handle: &Handle{
			hType: C.SQL_HANDLE_STMT,
			hPtr:  C.SQLHANDLE(*hnd),
		},
	}

	return stmt, nil
}

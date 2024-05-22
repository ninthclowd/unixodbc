package odbc

import "C"
import (
	"errors"
	"fmt"
	"github.com/ninthclowd/unixodbc/internal/api"
	"unsafe"
)

var (
	ErrConnectionDead            = errors.New("connection dead")
	ErrUnsupportedIsolationLevel = errors.New("isolation level is not supported")
)

type IsolationLevel uint64

const (
	LevelReadCommitted   = IsolationLevel(api.SQL_TRANSACTION_READ_COMMITTED)
	LevelReadUncommitted = IsolationLevel(api.SQL_TRANSACTION_READ_UNCOMMITTED)
	LevelRepeatableRead  = IsolationLevel(api.SQL_TRANSACTION_REPEATABLE_READ)
	LevelSerializable    = IsolationLevel(api.SQL_TRANSACTION_SERIALIZABLE)
)

type ReadOnlyMode uint64

const (
	ModeReadOnly  = ReadOnlyMode(api.SQL_MODE_READ_ONLY)
	ModeReadWrite = ReadOnlyMode(api.SQL_MODE_READ_WRITE)
	ModeDefault   = ReadOnlyMode(api.SQL_MODE_DEFAULT)
)

type Connection struct {
	*handle
	env *Environment
}

func (c *Connection) Ping() error {
	var dead api.SQLBIGINT
	_, err := c.result(api.SQLGetConnectAttr((*api.SQLHDBC)(c.hnd()),
		api.SQL_ATTR_CONNECTION_DEAD,
		(*api.SQLPOINTER)(unsafe.Pointer(&dead)),
		0,
		nil))
	if err != nil {
		return err
	}
	if (int64)(dead) == api.SQL_CD_TRUE {
		return ErrConnectionDead
	}
	return nil
}

func (c *Connection) Close() error {
	if _, err := c.result(api.SQLDisconnect((*api.SQLHDBC)(c.hnd()))); err != nil {
		return fmt.Errorf("disconnecting: %w", err)
	}
	return c.free()
}

func (c *Connection) Statement() (*Statement, error) {
	hnd, err := c.child(api.SQL_HANDLE_STMT)
	if err != nil {
		return nil, fmt.Errorf("unable to alloc new statement: %w", err)
	}

	stmt := &Statement{
		handle: hnd,
		conn:   c,
	}

	if err := stmt.SetConcurrency(ConcurrencyLock); err != nil {
		return nil, fmt.Errorf("setting concurrency lock: %w", err)
	}

	return stmt, nil
}

func (c *Connection) SetAutoCommit(autoCommit bool) error {
	val := api.SQL_AUTOCOMMIT_ON
	if !autoCommit {
		val = api.SQL_AUTOCOMMIT_OFF
	}
	_, err := c.result(api.SQLSetConnectAttr((*api.SQLHDBC)(c.hnd()),
		api.SQL_ATTR_AUTOCOMMIT,
		api.Const(val),
		api.SQL_IS_UINTEGER))
	return err
}

func (c *Connection) SetReadOnlyMode(readOnly ReadOnlyMode) error {
	_, err := c.result(api.SQLSetConnectAttr((*api.SQLHDBC)(c.hnd()),
		api.SQL_ATTR_ACCESS_MODE,
		api.Const(uint64(readOnly)),
		api.SQL_IS_UINTEGER))
	return err
}

func (c *Connection) SetIsolationLevel(level IsolationLevel) error {
	_, err := c.result(api.SQLSetConnectAttr((*api.SQLHDBC)(c.hnd()),
		api.SQL_ATTR_TXN_ISOLATION,
		api.Const(uint64(level)),
		api.SQL_IS_UINTEGER))
	return err
}

func (c *Connection) Commit() error {
	_, err := c.result(api.SQLEndTran(api.SQL_HANDLE_DBC, c.hnd(), api.SQL_COMMIT))
	return err
}

func (c *Connection) Rollback() error {
	_, err := c.result(api.SQLEndTran(api.SQL_HANDLE_DBC, c.hnd(), api.SQL_ROLLBACK))
	return err
}

package odbc

import "C"
import (
	"context"
	"fmt"
	"unicode/utf16"

	"github.com/ninthclowd/unixodbc/internal/api"
)

type PoolOption uint64

const (
	PoolOff            = PoolOption(api.SQL_CP_OFF)
	PoolPerDriver      = PoolOption(api.SQL_CP_ONE_PER_DRIVER)
	PoolPerEnvironment = PoolOption(api.SQL_CP_ONE_PER_HENV)
)

type Version uint64

const (
	Version380 = Version(api.SQL_OV_ODBC3_80)
	Version3   = Version(api.SQL_OV_ODBC3)
	Version2   = Version(api.SQL_OV_ODBC2)
)

//go:generate mockgen -source=environment.go -package mocks -destination ../mocks/environment.go
type Environment interface {
	SetVersion(version Version) error
	SetPoolOption(option PoolOption) error
	Connect(ctx context.Context, connStr string) (Connection, error)
	Close() error
}

func NewEnvironment() (Environment, error) {
	if err := setProcessConnectionPooling(PoolOff); err != nil {
		return nil, err
	}

	hnd, err := newEnvHandle()
	if err != nil {
		return nil, err
	}

	e := &environment{handle: hnd}

	return e, nil
}

func setProcessConnectionPooling(option PoolOption) error {
	if code := api.SQLSetEnvAttr(
		nil,
		api.SQL_ATTR_CONNECTION_POOLING,
		api.Const(uint64(option)),
		api.SQL_IS_UINTEGER,
	); code != api.SQL_SUCCESS {
		return fmt.Errorf("setting process connection pooling: received code %d", code)
	}
	return nil
}

var _ Environment = (*environment)(nil)

type environment struct {
	*handle
}

func (e *environment) SetVersion(version Version) error {

	_, err := e.result(api.SQLSetEnvAttr((*api.SQLHENV)(e.hnd()),
		api.SQL_ATTR_ODBC_VERSION,
		api.Const(uint64(version)),
		api.SQL_IS_UINTEGER))
	return err
}

func (e *environment) SetPoolOption(option PoolOption) error {
	_, err := e.result(api.SQLSetEnvAttr((*api.SQLHENV)(e.hnd()),
		api.SQL_ATTR_CONNECTION_POOLING,
		api.Const(uint64(option)),
		api.SQL_IS_UINTEGER))
	return err
}

func (e *environment) Connect(ctx context.Context, connStr string) (Connection, error) {
	hnd, err := e.child(api.SQL_HANDLE_DBC)
	if err != nil {
		return nil, fmt.Errorf("unable to alloc new connection: %w", err)
	}

	done := cancelHandleOnContext(ctx, hnd)

	connStrBytes := utf16.Encode([]rune(connStr))

	_, err = hnd.result(api.SQLDriverConnectW(
		(*api.SQLHDBC)(hnd.hnd()),
		nil,
		(*api.SQLWCHAR)(&connStrBytes[0]),
		api.SQLSMALLINT(len(connStrBytes)),
		nil,
		0,
		nil,
		api.SQL_DRIVER_NOPROMPT))

	done()

	if err == nil {
		err = ctx.Err()
	}
	if err != nil {
		_ = hnd.free()
		return nil, err
	}

	return &connection{handle: hnd, env: e}, nil

}

func (e *environment) Close() error {
	return e.free()
}

package odbc

import "C"
import (
	"context"
	"fmt"
	"github.com/ninthclowd/unixodbc/internal/api"
	"unicode/utf16"
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

func NewEnvironment() (*Environment, error) {
	hnd, err := newEnvHandle()
	if err != nil {
		return nil, err
	}

	e := &Environment{handle: hnd}

	return e, nil
}

type Environment struct {
	*handle
}

func (e *Environment) SetVersion(version Version) error {

	_, err := e.result(api.SQLSetEnvAttr((*api.SQLHENV)(e.hnd()),
		api.SQL_ATTR_ODBC_VERSION,
		api.Const(uint64(version)),
		api.SQL_IS_UINTEGER))
	return err
}

func (e *Environment) SetPoolOption(option PoolOption) error {
	_, err := e.result(api.SQLSetEnvAttr((*api.SQLHENV)(e.hnd()),
		api.SQL_ATTR_CONNECTION_POOLING,
		api.Const(uint64(option)),
		api.SQL_IS_UINTEGER))
	return err
}

func (e *Environment) Connect(ctx context.Context, connStr string) (*Connection, error) {
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

	return &Connection{handle: hnd, env: e}, nil

}

func (e *Environment) Close() error {
	return e.free()
}

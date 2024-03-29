package odbc

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
	Version3_80 = Version(api.SQL_OV_ODBC3_80)
	Version3    = Version(api.SQL_OV_ODBC3)
	Version2    = Version(api.SQL_OV_ODBC2)
)

func NewEnvironment(config *Config) (*Environment, error) {
	var capi odbcAPI
	if config != nil && config.api != nil {
		capi = config.api
	} else {
		capi = new(api.API)
	}
	hnd, err := newEnvHandle(capi)
	if err != nil {
		return nil, err
	}

	e := &Environment{handle: hnd}

	return e, nil
}

type Environment struct {
	handle
}

func (e *Environment) SetVersion(version Version) error {
	_, err := e.result(e.api().SQLSetEnvAttrConst((api.SQLHENV)(e.hnd()), api.SQL_ATTR_ODBC_VERSION, uint64(version)))
	return err
}

func (e *Environment) SetPoolOption(option PoolOption) error {
	_, err := e.result(e.api().SQLSetEnvAttrConst((api.SQLHENV)(e.hnd()), api.SQL_ATTR_CONNECTION_POOLING, uint64(option)))
	return err
}

// SetTraceFile enables unixodbc trace output to the specified file, or disables tracing if the filePath is empty
func (e *Environment) SetTraceFile(filePath string) error {
	val := api.SQL_OPT_TRACE_OFF

	if filePath != "" {
		val = api.SQL_OPT_TRACE_ON
		connStrBytes := []byte(filePath)
		_, err := e.result(e.api().SQLSetEnvAttrStr((api.SQLHENV)(e.hnd()), api.SQL_ATTR_TRACEFILE, api.SQLPOINTER(&connStrBytes), api.SQLINTEGER(len(connStrBytes))))
		if err != nil {
			return err
		}
	}

	_, err := e.result(e.api().SQLSetEnvAttrConst((api.SQLHENV)(e.hnd()), api.SQL_ATTR_TRACE, val))
	return err
}

func (e *Environment) Connect(ctx context.Context, connStr string) (*Connection, error) {
	hnd, err := e.child(api.SQL_HANDLE_DBC)
	if err != nil {
		return nil, fmt.Errorf("unable to alloc new connection: %w", err)
	}

	done := cancelHandleOnContext(ctx, hnd)

	connStrBytes := utf16.Encode([]rune(connStr))

	_, err = hnd.result(hnd.api().SQLDriverConnectW(
		(api.SQLHDBC)(hnd.hnd()),
		nil,
		connStrBytes,
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

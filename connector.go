package unixodbc

import (
	"context"
	"database/sql/driver"
	"errors"
	"github.com/ninthclowd/unixodbc/internal/cache"
	"github.com/ninthclowd/unixodbc/internal/odbc"
	"io"
)

// StaticConnStr converts a static connection string into ConnectionStringFactory usable by Connector
type StaticConnStr string

func (s StaticConnStr) ConnectionString() (string, error) {
	return string(s), nil
}

type ConnectionStringFactory interface {
	ConnectionString() (string, error)
}

var _ driver.Connector = (*Connector)(nil)
var _ io.Closer = (*Connector)(nil)

// Connector can be used with sql.OpenDB to allow more control of the unixodbc driver
type Connector struct {
	//ConnectionString is a factory that generates connection strings for each new connection that is opened.
	//Use StaticConnStr if you have a static connection string that does not need to change with each new connection.
	//Ex: If you are connecting using a system DSN called "myDatabase", this could be:
	//	StaticConnStr("DSN=myDatabase")
	ConnectionString   ConnectionStringFactory
	StatementCacheSize int

	UseODBCCursor bool   //SQL_ATTR_ODBC_CURSORS TODO
	PacketSize    uint32 //SQL_ATTR_PACKET_SIZE TODO
	TraceFile     string //SQL_ATTR_TRACEFILE TODO

	odbcEnvironment *odbc.Environment
}

// Close implements io.Closer
func (c *Connector) Close() error {
	return c.odbcEnvironment.Close()
}

func (c *Connector) initEnvironment(ctx context.Context) (err error) {
	if c.odbcEnvironment != nil {
		return nil
	}
	var env *odbc.Environment

	ctx, trace := Tracer.NewTask(ctx, "Connection::initEnvironment")
	defer trace.End()

	Tracer.WithRegion(ctx, "initializing ODBC environment", func() {
		env, err = odbc.NewEnvironment(nil)
	})
	if err != nil {
		return
	}

	Tracer.WithRegion(ctx, "setting version", func() {
		err = env.SetVersion(odbc.Version3_80)
	})
	if err != nil {
		return
	}

	//do not enable connection pooling at the driver level since go sql will be managing a connection pool
	Tracer.WithRegion(ctx, "setting pool option", func() {
		err = env.SetPoolOption(odbc.PoolOff)
	})
	if err != nil {
		return
	}

	//if err = env.SetTraceFile(c.TraceFile); err != nil {
	//	return
	//}

	c.odbcEnvironment = env
	return
}

// Connect implements driver.Connector
func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	ctx, trace := Tracer.NewTask(ctx, "Connect")
	defer trace.End()

	var err error
	if err = c.initEnvironment(ctx); err != nil {
		return nil, err
	}

	if c.ConnectionString == nil {
		return nil, errors.New("ConnectionString is required")
	}

	var connStr string
	Tracer.WithRegion(ctx, "generating connection string", func() {
		connStr, err = c.ConnectionString.ConnectionString()
	})
	if err != nil {
		return nil, err
	}

	conn := &Connection{
		connector: c,
		cachedStatements: cache.NewLRU[PreparedStatement](c.StatementCacheSize, func(key string, value *PreparedStatement) error {
			return value.odbcStatement.Close()
		}),
	}

	Tracer.WithRegion(ctx, "connecting", func() {
		conn.odbcConnection, err = c.odbcEnvironment.Connect(ctx, connStr)
	})

	if err != nil {
		return nil, err
	}
	if err = conn.odbcConnection.SetAutoCommit(true); err != nil {
		return nil, err
	}

	return conn, nil

}

// Driver implements driver.Connector
func (c *Connector) Driver() driver.Driver {
	return driverInstance
}

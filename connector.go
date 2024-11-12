package unixodbc

import (
	"context"
	"database/sql/driver"
	"errors"
	"github.com/ninthclowd/unixodbc/internal/cache"
	"github.com/ninthclowd/unixodbc/internal/odbc"
	"io"
	"sync"
)

// StaticConnStr converts a static connection string into ConnectionStringFactory usable by Connector
type StaticConnStr string

func (s StaticConnStr) ConnectionString(ctx context.Context) (string, error) {
	return string(s), nil
}

// ConnectionStringFactory can be implemented to provide dynamic connection strings for each new connection to the
// database, allowing for connections to systems that require token based authentication
type ConnectionStringFactory interface {
	ConnectionString(ctx context.Context) (string, error)
}

var _ driver.Connector = (*Connector)(nil)
var _ io.Closer = (*Connector)(nil)

// Connector can be used with sql.OpenDB to allow more control of the unixodbc driver
type Connector struct {
	// ConnectionString is a factory that generates connection strings for each new connection that is opened.
	//Use StaticConnStr if you have a static connection string that does not need to change with each new connection.
	//Ex: If you are connecting using a system DSN called "myDatabase", this could be:
	//	StaticConnStr("DSN=myDatabase")
	ConnectionString ConnectionStringFactory
	// StatementCacheSize is the number of prepared statements to cache for each connection.  The driver will cache
	// statements up to this limit and purge them using the least recently used algorithm.  0 will disable prepared
	// statement caching.
	StatementCacheSize int

	odbcEnvironment odbc.Environment

	initialized bool
	mux         sync.Mutex
}

// Close implements io.Closer
func (c *Connector) Close() error {
	return c.odbcEnvironment.Close()
}

// initialize the environment if it has not already been done
func (c *Connector) initialize(ctx context.Context) (err error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if c.initialized {
		return nil
	}

	ctx, trace := Tracer.NewTask(ctx, "connection::initialize")
	defer trace.End()

	if c.odbcEnvironment == nil {
		Tracer.WithRegion(ctx, "initializing ODBC environment", func() {
			c.odbcEnvironment, err = odbc.NewEnvironment()
		})
		if err != nil {
			return
		}
	}

	Tracer.WithRegion(ctx, "setting version", func() {
		err = c.odbcEnvironment.SetVersion(odbc.Version380)
	})
	if err != nil {
		return
	}

	//do not enable connection pooling at the driver level since go sql will be managing a connection pool
	Tracer.WithRegion(ctx, "setting pool option", func() {
		err = c.odbcEnvironment.SetPoolOption(odbc.PoolOff)
	})
	if err != nil {
		return
	}

	c.initialized = true
	return
}

// Connect implements driver.Connector
func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	ctx, trace := Tracer.NewTask(ctx, "Connect")
	defer trace.End()

	var err error
	if err = c.initialize(ctx); err != nil {
		return nil, err
	}

	if c.ConnectionString == nil {
		return nil, errors.New("ConnectionString is required")
	}

	var connStr string
	Tracer.WithRegion(ctx, "generating connection string", func() {
		connStr, err = c.ConnectionString.ConnectionString(ctx)
	})
	if err != nil {
		return nil, err
	}

	conn := &Connection{
		connector:          c,
		cachedStatements:   cache.NewLRU[PreparedStatement](c.StatementCacheSize, onCachePurged),
		uncachedStatements: make(map[*PreparedStatement]bool),
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

func onCachePurged(key string, value *PreparedStatement) error {
	return value.odbcStatement.Close()
}

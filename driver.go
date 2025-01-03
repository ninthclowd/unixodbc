package unixodbc

import (
	"database/sql"
	"database/sql/driver"
	"github.com/ninthclowd/unixodbc/internal/odbc"
)

const DefaultCacheSize = 0

var driverInstance = &Driver{}

func init() {
	sql.Register("unixodbc", driverInstance)
}

// OpenHandles reports the number of handles open to unixodbc. Useful in testing to ensure all handles are closed correctly.
func OpenHandles() int64 {
	return odbc.OpenHandles()
}

var _ driver.DriverContext = (*Driver)(nil)

type Driver struct{}

// OpenConnector implements driver.DriverContext
func (d *Driver) OpenConnector(connStr string) (driver.Connector, error) {
	return &Connector{
		ConnectionString:   StaticConnStr(connStr),
		StatementCacheSize: DefaultCacheSize,
	}, nil
}

// Open should never be called because driver.DriverContext is implemented
func (d *Driver) Open(connStr string) (driver.Conn, error) {
	panic("unexpected call to Open() from driver")
}

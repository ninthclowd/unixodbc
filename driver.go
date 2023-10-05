package unixodbc

import (
	"database/sql"
	"database/sql/driver"
)

var driverInstance = &Driver{}

func init() {
	sql.Register("unixodbc", driverInstance)
}

var _ driver.DriverContext = (*Driver)(nil)

type Driver struct{}

// OpenConnector implements driver.DriverContext
func (d *Driver) OpenConnector(connStr string) (driver.Connector, error) {
	return &Connector{
		ConnectionString:   StaticConnStr(connStr),
		StatementCacheSize: 0,
	}, nil
}

// Open should never be called because driver.DriverContext is implemented
func (d *Driver) Open(connStr string) (driver.Conn, error) {
	panic("unexpected call to Open() from driver")
}
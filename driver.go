package unixodbc

import (
	"context"
	"database/sql"
	"database/sql/driver"
)

var driverInstance = &odbcDriver{}

func init() {
	sql.Register("unixodbc", driverInstance)
}

var _ driver.DriverContext = (*odbcDriver)(nil)

type odbcDriver struct {
}

// OpenConnector implements driver.DriverContext
func (d *odbcDriver) OpenConnector(connStr string) (driver.Connector, error) {
	return Connector(WithConnectionString(connStr)), nil
}

// Open implements driver.Driver
func (d *odbcDriver) Open(connStr string) (driver.Conn, error) {
	if c, err := d.OpenConnector(connStr); err != nil {
		return nil, err
	} else {
		return c.Connect(context.Background())
	}
}

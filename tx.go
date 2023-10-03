package unixodbc

import "database/sql/driver"

var _ driver.Tx = (*TX)(nil)

type TX struct {
	conn *Connection
}

func (t *TX) Commit() error {
	if err := t.conn.odbcConnection.Commit(); err != nil {
		return err
	}
	return t.conn.endTx()
}

func (t *TX) Rollback() error {
	if err := t.conn.odbcConnection.Rollback(); err != nil {
		return err
	}
	return t.conn.endTx()
}

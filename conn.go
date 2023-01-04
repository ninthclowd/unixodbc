package unixodbc

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/ninthclowd/unixodbc/internal/api"
)

var (
	_ driver.QueryerContext     = (*conn)(nil)
	_ driver.ExecerContext      = (*conn)(nil)
	_ driver.ConnPrepareContext = (*conn)(nil)
	_ driver.ConnBeginTx        = (*conn)(nil)
	_ driver.Pinger             = (*conn)(nil)
	_ driver.Conn               = (*conn)(nil)
)

type conn struct {
	psCache   psCache
	connector *connector
	hnd       *api.Connection
	invalid   bool
}

// Prepare implements driver.Conn
func (c *conn) Prepare(query string) (driver.Stmt, error) {
	return nil, ErrDeprecated
}

// Close implements driver.Conn
func (c *conn) Close() error {
	if err := c.hnd.SQLDisconnect(); err != nil {
		return err
	}
	if err := c.hnd.Free(); err != nil {
		return err
	}
	c.invalid = true
	return nil
}

// Begin implements driver.Conn
func (c *conn) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{
		Isolation: driver.IsolationLevel(sql.LevelDefault),
		ReadOnly:  false,
	})
}

// BeginTx implements driver.ConnBeginTx
func (c *conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	//TODO implement me
	panic("implement me")
}

func (c *conn) prepareStatement(query string) (*stmt, error) {
	return c.psCache.Get(query, func() (*stmt, error) {
		odbcStmt, err := c.hnd.Statement()
		if err != nil {
			return nil, err
		}
		if err = odbcStmt.SQLPrepare(query); err != nil {
			_ = odbcStmt.Free()
			return nil, err
		}

		params, err := odbcStmt.Params()
		if err != nil {
			return nil, err
		}
		return &stmt{
			hnd:      odbcStmt,
			conn:     c,
			numInput: len(params),
		}, nil
	})
}

// PrepareContext implements driver.ConnPrepareContext
func (c *conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	type result struct {
		stmt *stmt
		err  error
	}

	ch := make(chan *result)

	go func() {
		r := new(result)
		r.stmt, r.err = c.prepareStatement(query)
		ch <- r
	}()

	select {
	case res := <-ch:
		return res.stmt, res.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// ExecContext implements driver.ExecerContext
func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return nil, driver.ErrSkip
}

// QueryContext implements driver.QueryerContext
func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return nil, driver.ErrSkip
}

// Ping implements driver.Pinger
func (c *conn) Ping(ctx context.Context) error {
	if c.invalid {
		return driver.ErrBadConn
	}
	return nil
}

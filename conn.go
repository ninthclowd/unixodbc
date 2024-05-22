package unixodbc

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/ninthclowd/unixodbc/internal/cache"
	"github.com/ninthclowd/unixodbc/internal/odbc"
	"time"
)

var (
	_ driver.QueryerContext     = (*Connection)(nil)
	_ driver.ExecerContext      = (*Connection)(nil)
	_ driver.ConnPrepareContext = (*Connection)(nil)
	_ driver.ConnBeginTx        = (*Connection)(nil)
	_ driver.Pinger             = (*Connection)(nil)
	_ driver.Conn               = (*Connection)(nil)
	_ driver.NamedValueChecker  = (*Connection)(nil)
	_ driver.Validator          = (*Connection)(nil)
)

var toODBCIsoLvl = map[sql.IsolationLevel]odbc.IsolationLevel{
	sql.LevelRepeatableRead:  odbc.LevelRepeatableRead,
	sql.LevelReadUncommitted: odbc.LevelReadUncommitted,
	sql.LevelSerializable:    odbc.LevelSerializable,
	sql.LevelReadCommitted:   odbc.LevelReadCommitted,
}

type Connection struct {
	connector        *Connector
	odbcConnection   *odbc.Connection
	openTX           *TX
	cachedStatements *cache.LRU[PreparedStatement]
}

// IsValid implements driver.Validator
func (c *Connection) IsValid() bool {
	if err := c.Ping(context.Background()); err != nil {
		return false
	}
	//TODO return false on cancelled queries?
	return true
}

// CheckNamedValue implements driver.NamedValueChecker
func (c *Connection) CheckNamedValue(value *driver.NamedValue) error {
	switch value.Value.(type) {
	case float64, int8, int16, int32, int64, string, bool, nil, time.Time, []byte:
		return nil
	default:
		return driver.ErrRemoveArgument
	}
}

// Close implements driver.Conn
func (c *Connection) Close() error {
	if err := c.cachedStatements.Purge(); err != nil {
		return err
	}
	if err := c.odbcConnection.Close(); err != nil {
		return err
	}
	c.odbcConnection = nil
	return nil
}

// Begin will never be called because driver.ConnBeginTx is implemented
func (c *Connection) Begin() (driver.Tx, error) {
	panic("unexpected call to Begin() from driver")
}

// BeginTx implements driver.ConnBeginTx
func (c *Connection) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	ctx, trace := Tracer.NewTask(ctx, "BeginTx")
	defer trace.End()
	var err error

	if sqlIsoLvl := sql.IsolationLevel(opts.Isolation); sqlIsoLvl != sql.LevelDefault {

		odbcIsoLvl, found := toODBCIsoLvl[sqlIsoLvl]
		if !found {
			return nil, fmt.Errorf("isolation level %d is not supported", opts.Isolation)
		}

		Tracer.WithRegion(ctx, "setIsolationLevel", func() {
			err = c.odbcConnection.SetIsolationLevel(odbcIsoLvl)
		})
		if err != nil {
			return nil, err
		}
	}

	if opts.ReadOnly {
		Tracer.WithRegion(ctx, "setReadOnly", func() {
			err = c.odbcConnection.SetReadOnlyMode(odbc.ModeReadOnly)
		})
		if err != nil {
			return nil, err
		}
	}

	Tracer.WithRegion(ctx, "setAutoCommit", func() {
		err = c.odbcConnection.SetAutoCommit(false)
	})
	if err != nil {
		return nil, err
	}

	c.openTX = &TX{conn: c}
	return c.openTX, nil
}

func (c *Connection) endTx() error {
	c.openTX = nil
	return c.odbcConnection.SetAutoCommit(true)
}

// Prepare will never be called because driver.ConnPrepareContext is implemented
func (c *Connection) Prepare(query string) (driver.Stmt, error) {
	panic("unexpected call to Prepare() from driver")
}

// PrepareContext implements driver.ConnPrepareContext
func (c *Connection) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	ctx, trace := Tracer.NewTask(ctx, "Connection::PrepareContext")
	defer trace.End()
	Tracer.Logf(ctx, "query", query)

	var stmt *PreparedStatement
	var err error

	Tracer.WithRegion(ctx, "Cache lookup", func() {
		stmt = c.cachedStatements.Get(query, true)
	})

	if stmt != nil {
		Tracer.WithRegion(ctx, "Reset parameters", func() {
			err = stmt.odbcStatement.ResetParams()
		})
		if err != nil {
			return nil, err
		}
		return stmt, nil
	}

	var st *odbc.Statement
	Tracer.WithRegion(ctx, "Create statement", func() {
		st, err = c.odbcConnection.Statement()
	})
	if err != nil {
		return nil, err
	}
	Tracer.WithRegion(ctx, "Prepare statement", func() {
		err = st.Prepare(ctx, query)
	})
	if err != nil {
		Tracer.WithRegion(ctx, "Close statement", func() {
			_ = st.Close()
		})
		return nil, err
	}

	var numParam int

	Tracer.WithRegion(ctx, "Read parameter count", func() {
		numParam, err = st.NumParams()
	})

	if err != nil {
		Tracer.WithRegion(ctx, "Close statement", func() {
			_ = st.Close()
		})
		return nil, err
	}
	return &PreparedStatement{odbcStatement: st, conn: c, numInput: numParam, query: query}, nil
}

// ExecContext implements driver.ExecerContext
func (c *Connection) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	ctx, trace := Tracer.NewTask(ctx, "Connection::ExecContext")
	defer trace.End()
	Tracer.Logf(ctx, "query", query)
	var st *odbc.Statement
	var err error

	Tracer.WithRegion(ctx, "Create statement", func() {
		st, err = c.odbcConnection.Statement()
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		Tracer.WithRegion(ctx, "Close statement", func() {
			_ = st.Close()
		})
	}()
	Tracer.WithRegion(ctx, "Bind parameters", func() {
		err = st.BindParams(toValues(args)...)
	})
	if err != nil {
		return nil, err
	}
	Tracer.WithRegion(ctx, "Execute statement", func() {
		err = st.ExecDirect(ctx, query)
	})
	if err != nil {
		return nil, err
	}
	return &result{lastInsertId: 0, rowsAffected: 0}, nil //TODO populate result
}

// QueryContext implements driver.QueryerContext
func (c *Connection) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	ctx, trace := Tracer.NewTask(ctx, "Connection::QueryContext")
	defer trace.End()
	Tracer.Logf(ctx, "query", query)

	var st *odbc.Statement
	var err error

	Tracer.WithRegion(ctx, "Create statement", func() {
		st, err = c.odbcConnection.Statement()
	})
	if err != nil {
		return nil, err
	}

	Tracer.WithRegion(ctx, "Bind parameters", func() {
		err = st.BindParams(toValues(args)...)
	})
	if err != nil {
		_ = st.Close()
		return nil, err
	}

	Tracer.WithRegion(ctx, "Executing statement", func() {
		err = st.ExecDirect(ctx, query)
	})
	if err != nil {
		_ = st.Close()
		return nil, err
	}
	var rs *odbc.RecordSet
	Tracer.WithRegion(ctx, "Getting recordset", func() {
		rs, err = st.RecordSet()
	})
	if err != nil {
		_ = st.Close()
		return nil, err
	}

	return &Rows{odbcRecordset: rs, closeStmtOnRSClose: st, ctx: ctx}, nil
}

// Ping implements driver.Pinger
func (c *Connection) Ping(ctx context.Context) error {
	ctx, trace := Tracer.NewTask(ctx, "Connection::Ping")
	defer trace.End()
	if c.odbcConnection == nil {
		return driver.ErrBadConn
	}
	switch err := c.odbcConnection.Ping(); {
	case errors.Is(err, odbc.ErrConnectionDead):
		return driver.ErrBadConn
	default:
		return err
	}
}

func toValues(args []driver.NamedValue) (values []interface{}) {
	values = make([]interface{}, len(args))
	for _, arg := range args {
		values[arg.Ordinal-1] = arg.Value
	}
	return
}

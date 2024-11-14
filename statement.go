package unixodbc

import (
	"context"
	"database/sql/driver"
	"github.com/ninthclowd/unixodbc/internal/odbc"
)

var (
	_ driver.StmtExecContext  = (*PreparedStatement)(nil)
	_ driver.StmtQueryContext = (*PreparedStatement)(nil)
	_ driver.Stmt             = (*PreparedStatement)(nil)
)

type PreparedStatement struct {
	odbcStatement odbc.Statement
	query         string
	conn          *Connection
	numInput      int
}

func (s *PreparedStatement) closeWithError(err error) error {
	if s.odbcStatement != nil {
		closeErr := s.odbcStatement.Close()
		s.odbcStatement = nil
		if closeErr != nil {
			return closeErr
		}
	}
	return err
}

// Close implements driver.Stmt
func (s *PreparedStatement) Close() error {
	delete(s.conn.uncachedStatements, s)
	if s.odbcStatement == nil {
		return nil
	}
	//move the statement to the LRU, closing the statement if no room in cache
	return s.conn.cachedStatements.Put(s.query, s)
}

// NumInput implements driver.Stmt
func (s *PreparedStatement) NumInput() int {
	return s.numInput
}

// Exec will never be called because driver.StmtExecContext is implemented
func (s *PreparedStatement) Exec(args []driver.Value) (driver.Result, error) {
	panic("unexpected call to Exec() from driver")
}

// ExecContext implements driver.StmtExecContext
func (s *PreparedStatement) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	if s.odbcStatement == nil {
		return nil, odbc.ErrInvalidHandle
	}
	ctx, trace := Tracer.NewTask(ctx, "statement::ExecContext")
	defer trace.End()
	var err error
	Tracer.WithRegion(ctx, "BindParams", func() {
		err = s.odbcStatement.BindParams(toValues(args)...)
	})
	if err != nil {
		return nil, s.closeWithError(err)
	}
	Tracer.WithRegion(ctx, "Execute", func() {
		err = s.odbcStatement.Execute(ctx)
	})
	if err != nil {
		return nil, s.closeWithError(err)
	}
	return &result{lastInsertId: 0, rowsAffected: 0}, nil //TODO stats for exec
}

// Query will never be called because driver.StmtQueryContext is implemented
func (s *PreparedStatement) Query(args []driver.Value) (driver.Rows, error) {
	panic("unexpected call to Query() from driver")
}

// QueryContext implements driver.StmtQueryContext
func (s *PreparedStatement) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	if s.odbcStatement == nil {
		return nil, odbc.ErrInvalidHandle
	}
	ctx, trace := Tracer.NewTask(ctx, "statement::QueryContext")
	defer trace.End()
	var err error
	Tracer.WithRegion(ctx, "BindParams", func() {
		err = s.odbcStatement.BindParams(toValues(args)...)
	})
	if err != nil {
		return nil, s.closeWithError(err)
	}
	Tracer.WithRegion(ctx, "Execute", func() {
		err = s.odbcStatement.Execute(ctx)
	})
	if err != nil {
		return nil, s.closeWithError(err)
	}

	var rs odbc.RecordSet
	Tracer.WithRegion(ctx, "RecordSet", func() {
		rs, err = s.odbcStatement.RecordSet()
	})
	if err != nil {
		return nil, s.closeWithError(err)
	}
	return &Rows{odbcRecordset: rs, ctx: ctx}, nil
}

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
	rows          *Rows
	numInput      int
}

// Close implements driver.Stmt
func (s *PreparedStatement) Close() error {
	errs := make(MultipleErrors)
	if s.rows != nil { //TODO is this check needed? will the driver always close rows before the statement
		errs["closing recordset"] = s.rows.Close()
		s.rows = nil
	}
	errs["cache"] = s.conn.cachedStatements.Put(s.query, s)
	return errs.Error()
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
	ctx, trace := Tracer.NewTask(ctx, "statement::ExecContext")
	defer trace.End()
	var err error
	Tracer.WithRegion(ctx, "BindParams", func() {
		err = s.odbcStatement.BindParams(toValues(args)...)
	})
	if err != nil {
		return nil, err
	}
	Tracer.WithRegion(ctx, "Execute", func() {
		err = s.odbcStatement.Execute(ctx)
	})
	if err != nil {
		return nil, err
	}
	return &result{lastInsertId: 0, rowsAffected: 0}, nil //TODO stats for exec
}

// Query will never be called because driver.StmtQueryContext is implemented
func (s *PreparedStatement) Query(args []driver.Value) (driver.Rows, error) {
	panic("unexpected call to Query() from driver")
}

// QueryContext implements driver.StmtQueryContext
func (s *PreparedStatement) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	ctx, trace := Tracer.NewTask(ctx, "statement::QueryContext")
	defer trace.End()
	var err error
	Tracer.WithRegion(ctx, "BindParams", func() {
		err = s.odbcStatement.BindParams(toValues(args)...)
	})
	if err != nil {
		return nil, err
	}
	Tracer.WithRegion(ctx, "Execute", func() {
		err = s.odbcStatement.Execute(ctx)
	})
	if err != nil {
		return nil, err
	}

	var rs odbc.RecordSet
	Tracer.WithRegion(ctx, "recordSet", func() {
		rs, err = s.odbcStatement.RecordSet()
	})
	if err != nil {
		return nil, err
	}
	return &Rows{odbcRecordset: rs, ctx: ctx}, nil
}

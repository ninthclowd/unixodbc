package unixodbc

import (
	"context"
	"database/sql/driver"
	"errors"
	"github.com/ninthclowd/unixodbc/internal/api"
	"sync"
)

var (
	_                   driver.StmtExecContext  = (*stmt)(nil)
	_                   driver.StmtQueryContext = (*stmt)(nil)
	_                   driver.Stmt             = (*stmt)(nil)
	ErrTimeoutStatement                         = errors.New("timeout waiting for statement")
	ErrDeprecated                               = errors.New("method is deprecated. use context aware methods")
)

type stmt struct {
	hnd      *api.Statement
	conn     *conn
	closed   bool
	numInput int
	mux      sync.Mutex

	resultSetLock sync.Mutex
}

// Close implements driver.Stmt
func (s *stmt) Close() error {
	s.mux.Lock()
	defer s.mux.Unlock()
	if err := s.hnd.Free(); err != nil {
		return err
	}
	s.closed = true
	return nil
}

// NumInput implements driver.Stmt
func (s *stmt) NumInput() int {
	return s.numInput
}

// Exec implements driver.Stmt
func (s *stmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, ErrDeprecated
}

// Query implements driver.Stmt
func (s *stmt) Query(args []driver.Value) (driver.Rows, error) {
	return nil, ErrDeprecated
}
func (s *stmt) lockForResultSet(ctx context.Context) error {
	ch := make(chan bool)
	go func() {
		s.resultSetLock.Lock()
		if ctx.Err() != nil {
			s.resultSetLock.Unlock()
		} else {
			close(ch)
		}
	}()
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *stmt) execContext(ctx context.Context) error {
	ch := make(chan error)
	go func() {
		err := s.hnd.SQLExecute()
		if ctx.Err() == nil {
			ch <- err
		}
	}()
	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		_ = s.hnd.Cancel()
		return ctx.Err()
	}
}

// QueryContext implements driver.StmtQueryContext
func (s *stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	if err := s.lockForResultSet(ctx); err != nil {
		return nil, err
	}
	if err := s.bindParams(ctx, args); err != nil {
		s.resultSetLock.Unlock()
		return nil, err
	}

	if err := s.execContext(ctx); err != nil {
		s.resultSetLock.Unlock()
		return nil, err
	}

	cols, err := s.hnd.Columns()
	if err != nil {
		s.resultSetLock.Unlock()
		return nil, err
	}

	return &rows{
		stmt:    s,
		columns: cols,
	}, nil
}

func (s *stmt) bindParams(ctx context.Context, args []driver.NamedValue) error {
	//TODO context
	params, err := s.hnd.Params()
	if err != nil {
		return err
	}
	for i, arg := range args {
		if err = params[i].Bind(arg.Value); err != nil {
			return err
		}
	}
	return nil

}

// ExecContext implements driver.StmtExecContext
func (s *stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	if err := s.lockForResultSet(ctx); err != nil {
		return nil, err
	}
	defer s.resultSetLock.Unlock()
	if err := s.bindParams(ctx, args); err != nil {
		return nil, err
	}
	if err := s.execContext(ctx); err != nil {
		return nil, err
	}

	return &result{
		lastInsertId: 0,
		rowsAffected: 0, //TODO
	}, nil
}

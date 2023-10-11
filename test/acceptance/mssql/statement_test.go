package acceptance_test

import (
	"context"
	"database/sql"
	_ "github.com/ninthclowd/unixodbc"
	"testing"
)

func newTestStatement(t *testing.T) (stmt *sql.Stmt, ctx context.Context, finish func()) {
	var conn *sql.Conn
	var done func()
	_, conn, ctx, done = newTestConnection(t)

	stmt, err := conn.PrepareContext(ctx, "SELECT 1 WHERE ? = 52.42")
	if err != nil {
		t.Fatalf("expected no error from PrepareContext, got: %s", err.Error())
	}
	finish = func() {
		stmt.Close()
		done()
	}
	return
}

func TestStatement_QueryContext(t *testing.T) {
	stmt, ctx, finish := newTestStatement(t)
	defer finish()

	rows, err := stmt.QueryContext(ctx, 52.42)
	if err != nil {
		t.Fatalf("expected no error, got: %s", err.Error())
	}
	if rows.Next() != true {

		t.Errorf("expected Next() to return true")
	}
}

func TestStatement_ExecContext(t *testing.T) {
	stmt, ctx, finish := newTestStatement(t)
	defer finish()

	_, err := stmt.ExecContext(ctx, 52.42)
	if err != nil {
		t.Fatalf("expected no error, got: %s", err.Error())
	}

}

package acceptance_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/ninthclowd/unixodbc"
	"os"
	"testing"
	"time"
)

func newTestConnection(t *testing.T) (db *sql.DB, conn *sql.Conn, ctx context.Context, finish func()) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		t.Skip("DB_HOST not set")
	}
	db, err := sql.Open("unixodbc", fmt.Sprintf("DSN=MariaDB;SERVER=%s", host))
	if err != nil {
		t.Fatalf("unable to open database: %s", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	conn, err = db.Conn(ctx)
	if err != nil {
		t.Fatalf(fmt.Sprintf("unable to connect to database: %s", err.Error()))
	}

	finish = func() {
		cancel()
		conn.Close()
		db.Close()
	}

	return
}

func TestConnection_Ping(t *testing.T) {
	_, conn, ctx, finish := newTestConnection(t)
	defer finish()

	if err := conn.PingContext(ctx); err != nil {
		t.Fatalf("expected no error, got: %s", err.Error())
	}
}

func TestConnection_Transaction_Commit(t *testing.T) {
	_, conn, ctx, finish := newTestConnection(t)
	defer finish()

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
		ReadOnly:  true,
	})

	if err != nil {
		t.Fatalf("expected no error from BeginTx, got: %s", err.Error())
	}

	_, err = tx.ExecContext(ctx, "SELECT 1 WHERE 1=1")
	if err != nil {
		t.Fatalf("expected no error from ExecContext, got: %s", err.Error())
	}
	stmt1, err := tx.PrepareContext(ctx, "SELECT 2 WHERE 1=1")
	if err != nil {
		t.Fatalf("expected no error from PrepareContext, got: %s", err.Error())
	}
	_, err = stmt1.ExecContext(ctx)
	stmt2, err := tx.PrepareContext(ctx, "SELECT 3 WHERE 1=1")
	if err != nil {
		t.Fatalf("expected no error from PrepareContext, got: %s", err.Error())
	}
	_, err = stmt2.ExecContext(ctx)
	if err != nil {
		t.Fatalf("expected no error from ExecContext, got: %s", err.Error())
	}
	err = tx.Commit()
	if err != nil {
		t.Fatalf("expected no error from Commit, got: %s", err.Error())
	}
}

func TestConnection_Transaction_Rollback(t *testing.T) {
	_, conn, ctx, finish := newTestConnection(t)
	defer finish()

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		t.Fatalf("expected no error from BeginTx, got: %s", err.Error())
	}

	_, err = tx.ExecContext(ctx, "SELECT 1 WHERE 1=1")
	if err != nil {
		t.Fatalf("expected no error from ExecContext, got: %s", err.Error())
	}
	stmt1, err := tx.PrepareContext(ctx, "SELECT 2 WHERE 1=1")
	if err != nil {
		t.Fatalf("expected no error from PrepareContext, got: %s", err.Error())
	}
	_, err = stmt1.ExecContext(ctx)
	stmt2, err := tx.PrepareContext(ctx, "SELECT 3 WHERE 1=1")
	if err != nil {
		t.Fatalf("expected no error from PrepareContext, got: %s", err.Error())
	}
	_, err = stmt2.ExecContext(ctx)
	if err != nil {
		t.Fatalf("expected no error from ExecContext, got: %s", err.Error())
	}
	err = tx.Rollback()
	if err != nil {
		t.Fatalf("expected no error from Rollback, got: %s", err.Error())
	}
}

func TestConnection_ExecContext_Cancel(t *testing.T) {
	_, conn, ctx, finish := newTestConnection(t)
	defer finish()

	start := time.Now()
	_, err := conn.ExecContext(ctx, "SELECT SLEEP(1)")
	elapsed := time.Since(start)
	if elapsed.Seconds() < 1 {
		t.Error("sleep validation query returned before 1 second")
	}

	execCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	start = time.Now()
	_, err = conn.ExecContext(execCtx, "SELECT SLEEP(5)")
	elapsed = time.Since(start)
	if elapsed.Seconds() > 1 {
		t.Fatalf("query did not return immediately when cancelled")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected error to be a context error. got: %s", err.Error())
	}

	_, err = conn.ExecContext(ctx, "SELECT 1")
	if err != nil {
		t.Errorf("expected no error from subsequent exec, got: %s", err.Error())
	}
}

func TestConnection_QueryContext_Cancel(t *testing.T) {
	_, conn, ctx, finish := newTestConnection(t)
	defer finish()

	start := time.Now()
	_, err := conn.ExecContext(ctx, "SELECT SLEEP(1)")
	elapsed := time.Since(start)
	if elapsed.Seconds() < 1 {
		t.Error("sleep validation query returned before 1 second")
	}

	execCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	start = time.Now()
	_, err = conn.QueryContext(execCtx, "SELECT SLEEP(2)")
	elapsed = time.Since(start)
	if elapsed.Seconds() > 1 {
		t.Fatalf("query did not return immediately when cancelled")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected error to be a context error.  got: %s", err.Error())
	}

	_, err = conn.QueryContext(ctx, "SELECT 1")
	if err != nil {
		t.Errorf("expected no error from subsequent query, got: %s", err.Error())
	}
}

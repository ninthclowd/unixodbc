package acceptance_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ninthclowd/unixodbc"
	"os"
	"strings"
	"testing"
	"time"
)

func newTestConnection(t *testing.T) (db *sql.DB, conn *sql.Conn, ctx context.Context, finish func()) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		t.Skip("DB_HOST not set")
	}
	db, err := sql.Open("unixodbc", fmt.Sprintf("Driver=MariaDB Unicode;Server=%s;UID=root;PWD=test", host))
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
		_ = conn.Close()
		_ = db.Close()
		if count := unixodbc.OpenHandles(); count > 0 {
			t.Fatalf("%d open ODBC handles", count)
		}
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

func TestValidateSleep(t *testing.T) {
	_, conn, ctx, finish := newTestConnection(t)
	defer finish()

	start := time.Now()
	_, err := conn.ExecContext(ctx, "SELECT SLEEP(1)")
	elapsed := time.Since(start)
	if elapsed.Seconds() < 1 {
		t.Error("sleep validation query returned before 1 second")
	}
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestConnection_ExecContext_Cancel(t *testing.T) {
	_, conn, ctx, finish := newTestConnection(t)
	defer finish()

	execCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := conn.ExecContext(execCtx, "SELECT SLEEP(2)")
	elapsed := time.Since(start)
	if elapsed.Seconds() > 1 {
		t.Errorf("query did not return immediately when cancelled")
	}

	if err == nil || !strings.Contains(err.Error(), "[HY018:1317]") {
		t.Errorf("expected a cancellation error, got: %v", err)
	}

	_, err = conn.ExecContext(ctx, "SELECT 1")
	if err != nil {
		t.Errorf("expected no error from subsequent exec, got: %s", err.Error())
	}
}

func TestConnection_QueryContext_Cancel(t *testing.T) {
	_, conn, ctx, finish := newTestConnection(t)
	defer finish()

	execCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := conn.QueryContext(execCtx, "SELECT SLEEP(2)")
	elapsed := time.Since(start)
	if elapsed.Seconds() > 1 {
		t.Errorf("query did not return immediately when cancelled")
	}
	if err == nil || !strings.Contains(err.Error(), "[HY018:1317]") {
		t.Errorf("expected a cancellation error, got: %v", err)
	}

	_, err = conn.QueryContext(ctx, "SELECT 1")
	if err != nil {
		t.Errorf("expected no error from subsequent query, got: %s", err.Error())
	}
}

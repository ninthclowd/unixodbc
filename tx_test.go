package unixodbc

import (
	"context"
	"database/sql"
	"github.com/ninthclowd/unixodbc/internal/odbc"
	"testing"
)

func TestTX_Commit(t *testing.T) {
	ctrl, conn, mockConn := testDBConnection(t, 0)
	defer ctrl.Finish()

	ctx := context.Background()

	mockConn.EXPECT().SetIsolationLevel(odbc.LevelReadCommitted).Return(nil).Times(1)
	mockConn.EXPECT().SetReadOnlyMode(odbc.ModeReadOnly).Return(nil).Times(1)
	mockConn.EXPECT().SetAutoCommit(false).Return(nil).Times(1)

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		t.Fatalf("expected no error preparing statement, got %v", err)
	}

	mockConn.EXPECT().Commit().Return(nil).Times(1)
	mockConn.EXPECT().SetAutoCommit(true).Return(nil).Times(1)

	gotErr := tx.Commit()
	if gotErr != nil {
		t.Fatalf("expected no error, got %v", gotErr)
	}

}

func TestTX_Rollback(t *testing.T) {
	ctrl, conn, mockConn := testDBConnection(t, 0)
	defer ctrl.Finish()

	ctx := context.Background()

	mockConn.EXPECT().SetIsolationLevel(odbc.LevelReadCommitted).Return(nil).Times(1)
	mockConn.EXPECT().SetReadOnlyMode(odbc.ModeReadOnly).Return(nil).Times(1)
	mockConn.EXPECT().SetAutoCommit(false).Return(nil).Times(1)

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		t.Fatalf("expected no error preparing statement, got %v", err)
	}

	mockConn.EXPECT().SetAutoCommit(true).Return(nil).Times(1)
	mockConn.EXPECT().Rollback().Return(nil).Times(1)
	gotErr := tx.Rollback()
	if gotErr != nil {
		t.Fatalf("expected no error, got %v", gotErr)
	}

}

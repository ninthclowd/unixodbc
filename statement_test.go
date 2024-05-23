package unixodbc

import (
	"context"
	"database/sql/driver"
	"github.com/golang/mock/gomock"
	"github.com/ninthclowd/unixodbc/internal/cache"
	"github.com/ninthclowd/unixodbc/internal/mocks"
	"reflect"
	"testing"
)

func testStatement(t *testing.T) (ctrl *gomock.Controller, stmt *PreparedStatement, mockStmt *mocks.MockStatement) {
	ctrl = gomock.NewController(t)
	mockStmt = mocks.NewMockStatement(ctrl)
	mockConn := mocks.NewMockConnection(ctrl)
	conn := &Connection{
		odbcConnection:   mockConn,
		cachedStatements: cache.NewLRU[PreparedStatement](1, onCachePurged),
	}
	stmt = &PreparedStatement{
		odbcStatement: mockStmt,
		query:         "SELECT * FROM foo WHERE bar = ?",
		conn:          conn,
		numInput:      1,
	}
	return
}

func TestPreparedStatement_QueryContext(t *testing.T) {
	ctrl, stmt, mockStmt := testStatement(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockRS := mocks.NewMockRecordSet(ctrl)
	mockStmt.EXPECT().NumParams().Return(1, nil).AnyTimes()
	mockStmt.EXPECT().BindParams("value").Return(nil).AnyTimes()
	mockStmt.EXPECT().Execute(gomock.Any()).Return(nil)
	mockStmt.EXPECT().RecordSet().Return(mockRS, nil)

	rs, err := stmt.QueryContext(ctx, []driver.NamedValue{{
		Name:    "",
		Ordinal: 1,
		Value:   "value",
	}})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	rows, ok := rs.(*Rows)
	if !ok {
		t.Fatalf("expected Rows, got %v", reflect.TypeOf(rs))
	}
	if rows.odbcRecordset != mockRS {
		t.Errorf("expected recordset to be set on Rows, got %v", rows.odbcRecordset)
	}
	if gotNumInput := stmt.NumInput(); gotNumInput != 1 {
		t.Errorf("expected 1 input, got %v", gotNumInput)
	}
	if err := stmt.Close(); err != nil {
		t.Errorf("expected no error closing rows, got %v", rows.Close())
	}

}

func TestPreparedStatement_ExecContext(t *testing.T) {
	ctrl, stmt, mockStmt := testStatement(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockStmt.EXPECT().NumParams().Return(1, nil).AnyTimes()
	mockStmt.EXPECT().BindParams("value").Return(nil).AnyTimes()
	mockStmt.EXPECT().Execute(gomock.Any()).Return(nil)

	res, err := stmt.ExecContext(ctx, []driver.NamedValue{{
		Name:    "",
		Ordinal: 1,
		Value:   "value",
	}})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if lastInsertId, _ := res.LastInsertId(); lastInsertId != 0 {
		t.Errorf("expected LastInsertId to be 0, got %v", lastInsertId)
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected != 0 {
		t.Errorf("expected RowsAffected to be 0, got %v", rowsAffected)
	}

}

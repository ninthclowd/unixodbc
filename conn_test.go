package unixodbc

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/ninthclowd/unixodbc/internal/cache"
	"github.com/ninthclowd/unixodbc/internal/mocks"
	"github.com/ninthclowd/unixodbc/internal/odbc"
	"testing"
	"time"
)

func testDBConnection(t *testing.T, cacheSize int) (ctrl *gomock.Controller, conn *sql.Conn, mockConn *mocks.MockConnection) {

	ctrl = gomock.NewController(t)
	mockEnv := mocks.NewMockEnvironment(ctrl)

	mockEnv.EXPECT().SetVersion(odbc.Version380).Return(nil).Times(1)
	mockEnv.EXPECT().SetPoolOption(odbc.PoolOff).Return(nil).Times(1)

	connString := "connString"

	db := sql.OpenDB(&Connector{
		ConnectionString:   StaticConnStr(connString),
		StatementCacheSize: cacheSize,
		odbcEnvironment:    mockEnv,
	})

	ctx := context.Background()

	mockConn = mocks.NewMockConnection(ctrl)
	mockConn.EXPECT().SetAutoCommit(true).Return(nil).Times(1)

	mockEnv.EXPECT().Connect(gomock.Any(), connString).Return(mockConn, nil)
	var err error
	conn, err = db.Conn(ctx)
	if err != nil {
		t.Fatalf("unable to open connection: %v", err)
	}

	return
}

func testConnection(t *testing.T) (ctrl *gomock.Controller, conn *Connection, mockConn *mocks.MockConnection) {
	ctrl = gomock.NewController(t)

	mockConn = mocks.NewMockConnection(ctrl)
	conn = &Connection{
		odbcConnection:     mockConn,
		cachedStatements:   cache.NewLRU[PreparedStatement](1, onCachePurged),
		uncachedStatements: map[*PreparedStatement]bool{},
	}
	return
}

func TestConnection_IsValid(t *testing.T) {
	ctrl, conn, mockConn := testConnection(t)
	defer ctrl.Finish()

	mockConn.EXPECT().Ping().Return(nil).Times(1)

	if conn.IsValid() != true {
		t.Fatal("expected IsValid to be true if the Ping is successful")
	}

	mockConn.EXPECT().Ping().Return(odbc.ErrConnectionDead).Times(1)

	if conn.IsValid() != false {
		t.Fatal("expected IsValid to be false if the Ping fails")
	}

}

func TestConnection_Ping(t *testing.T) {

	type Test struct {
		Description   string
		HasConnection bool
		PingError     error
		WantError     error
	}
	tests := []Test{
		{
			Description:   "it should return ErrBadConn if Ping returns ErrConnectionDead",
			WantError:     driver.ErrBadConn,
			PingError:     odbc.ErrConnectionDead,
			HasConnection: true,
		},
		{
			Description:   "it should return ErrBadConn if no connection",
			WantError:     driver.ErrBadConn,
			HasConnection: false,
		},

		{
			Description:   "it should return no error if ping is successful",
			WantError:     nil,
			PingError:     nil,
			HasConnection: true,
		},
	}
	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			ctrl, conn, mockConn := testConnection(t)
			defer ctrl.Finish()

			mockConn.EXPECT().Ping().Return(test.PingError).AnyTimes()

			if !test.HasConnection {
				conn.odbcConnection = nil
			}

			gotErr := conn.Ping(context.Background())

			if !errors.Is(gotErr, test.WantError) {
				t.Fatalf("expected error %v, got %v", test.WantError, gotErr)
			}
		})
	}
}

func TestConnection_CheckNamedValue(t *testing.T) {

	type Test struct {
		Description  string
		Value        any
		ShouldReject bool
	}
	tests := []Test{
		{
			Description:  "int8",
			Value:        int8(1),
			ShouldReject: false,
		},
		{
			Description:  "int16",
			Value:        int16(1),
			ShouldReject: false,
		},
		{
			Description:  "int32",
			Value:        int32(1),
			ShouldReject: false,
		},
		{
			Description:  "int64",
			Value:        int64(1),
			ShouldReject: false,
		},
		{
			Description:  "string",
			Value:        "string",
			ShouldReject: false,
		},
		{
			Description:  "boolean",
			Value:        true,
			ShouldReject: false,
		},
		{
			Description:  "nil",
			Value:        nil,
			ShouldReject: false,
		},
		{
			Description:  "time",
			Value:        time.Now(),
			ShouldReject: false,
		},
		{
			Description:  "[]byte",
			Value:        []byte{1, 2, 3},
			ShouldReject: false,
		},
		{
			Description:  "struct",
			Value:        struct{}{},
			ShouldReject: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			ctrl, conn, _ := testConnection(t)
			defer ctrl.Finish()

			gotErr := conn.CheckNamedValue(&driver.NamedValue{
				Name:    "Col",
				Ordinal: 0,
				Value:   test.Value,
			})
			if test.ShouldReject {
				if !errors.Is(gotErr, driver.ErrRemoveArgument) {
					t.Fatal("expected ErrRemoveArgument")
				}
			} else if gotErr != nil {
				t.Fatalf("expected no error but got %v", gotErr)

			}

		})

	}
}

func TestConnection_Close(t *testing.T) {

	ctrl, conn, mockConn := testConnection(t)
	defer ctrl.Finish()

	mockConn.EXPECT().Close().Return(nil).Times(1)

	err := conn.Close()
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if conn.odbcConnection != nil {
		t.Error("expected connection to be closed")
	}

}

func TestConnection_BeginTx(t *testing.T) {
	ctrl, conn, mockConn := testConnection(t)
	defer ctrl.Finish()

	mockConn.EXPECT().SetIsolationLevel(odbc.LevelReadCommitted).Return(nil).Times(1)
	mockConn.EXPECT().SetReadOnlyMode(odbc.ModeReadOnly).Return(nil).Times(1)
	mockConn.EXPECT().SetAutoCommit(false).Return(nil).Times(1)
	tx, err := conn.BeginTx(context.Background(), driver.TxOptions{
		Isolation: driver.IsolationLevel(sql.LevelReadCommitted),
		ReadOnly:  true,
	})
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if tx == nil {
		t.Errorf("expected a transaction but got nil")
	}

	mockConn.EXPECT().Commit().Return(nil).Times(1)
	mockConn.EXPECT().SetAutoCommit(true).Return(nil).Times(1)

	err = tx.Commit()
	if err != nil {
		t.Errorf("expected no error from commit but got %v", err)
	}
}

func TestConnection_PrepareContext(t *testing.T) {
	ctrl, conn, mockConn := testConnection(t)
	defer ctrl.Finish()

	q := "SELECT * FROM foo WHERE bar = ?"

	ctx := context.Background()
	mockStmt1 := mocks.NewMockStatement(ctrl)

	mockStmt1.EXPECT().Prepare(gomock.Any(), q).Return(nil).Times(1)
	mockStmt1.EXPECT().NumParams().Return(1, nil).Times(1)
	mockConn.EXPECT().Statement().Return(mockStmt1, nil).Times(1)

	stmt1, err := conn.PrepareContext(ctx, q)
	if err != nil {
		t.Fatalf("expected no error from prepareContext but got %v", err)
	}
	ps, ok := stmt1.(*PreparedStatement)
	if !ok {
		t.Fatalf("expected a statement to be returnedbut got %v", err)
	}
	if ps.odbcStatement != mockStmt1 {
		t.Errorf("expected statement to be %v but got %v", mockStmt1, ps.odbcStatement)

	}
	if gotNumInput := ps.NumInput(); gotNumInput != 1 {
		t.Errorf("expected num input to be %v but got %v", 1, gotNumInput)
	}

	//since prepared statement caching is on, the statement shouldn't be closed
	stmt1.Close()

	//parameters should be reset when the statement is reused
	mockStmt1.EXPECT().ResetParams().Times(1)

	stmt2, err := conn.PrepareContext(ctx, q)
	if err != nil {
		t.Fatalf("expected no error from prepareContext but got %v", err)
	}
	if stmt1 != stmt2 {
		t.Fatalf("expected the second statement to be the same but got %v", stmt2)
	}

	//since prepared statement caching is on, the statement shouldn't be closed
	stmt2.Close()

	q3 := "SELECT * FROM all"

	mockStmt2 := mocks.NewMockStatement(ctrl)

	mockStmt2.EXPECT().Prepare(gomock.Any(), q3).Return(nil).Times(1)
	mockStmt2.EXPECT().NumParams().Return(1, nil).Times(1)
	mockConn.EXPECT().Statement().Return(mockStmt2, nil).Times(1)

	//first statement should be evicted from cache when the next statement is created
	mockStmt1.EXPECT().Close().Return(nil).Times(1)

	stmt3, err := conn.PrepareContext(ctx, q3)
	if err != nil {
		t.Fatalf("expected no error from prepareContext but got %v", err)
	}
	if ps, ok := stmt3.(*PreparedStatement); !ok {
		t.Fatalf("expected a statement to be returnedbut got %v", err)
	} else if ps.odbcStatement != mockStmt2 {
		t.Errorf("expected statement to be %v but got %v", mockStmt2, ps.odbcStatement)

	}

	stmt3.Close()

	//first statement should be evicted from cache when the next statement is created
	mockStmt2.EXPECT().Close().Return(nil).Times(1)

	mockConn.EXPECT().Close().Return(nil).Times(1)

	conn.Close()

}

func TestConnection_ExecContext(t *testing.T) {
	ctrl, conn, mockConn := testConnection(t)
	defer ctrl.Finish()

	q := "SELECT * FROM foo WHERE bar = ?"

	ctx := context.Background()
	mockStmt := mocks.NewMockStatement(ctrl)

	mockConn.EXPECT().Statement().Return(mockStmt, nil).Times(1)

	mockStmt.EXPECT().BindParams("Value").Return(nil).Times(1)
	mockStmt.EXPECT().Close().Return(nil).Times(1)
	mockStmt.EXPECT().ExecDirect(gomock.Any(), q).Return(nil).Times(1)

	_, err := conn.ExecContext(ctx, q, []driver.NamedValue{
		{Name: "", Ordinal: 1, Value: "Value"},
	})
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

}

func TestConnection_QueryContext(t *testing.T) {
	ctrl, conn, mockConn := testConnection(t)
	defer ctrl.Finish()

	q := "SELECT * FROM foo WHERE bar = ?"

	ctx := context.Background()
	mockStmt := mocks.NewMockStatement(ctrl)

	mockConn.EXPECT().Statement().Return(mockStmt, nil).Times(1)

	mockStmt.EXPECT().BindParams("Value").Return(nil).Times(1)
	mockStmt.EXPECT().ExecDirect(gomock.Any(), q).Return(nil).Times(1)

	mockRS := mocks.NewMockRecordSet(ctrl)

	mockStmt.EXPECT().RecordSet().Return(mockRS, nil).Times(1)

	gotRows, err := conn.QueryContext(ctx, q, []driver.NamedValue{
		{Name: "", Ordinal: 1, Value: "Value"},
	})
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
	r, ok := gotRows.(*Rows)
	if !ok {
		t.Fatalf("unexpected rows result. got %v", gotRows)
	}
	if r.closeStmtOnRSClose != mockStmt {
		t.Errorf("stmt not populated on rows. got %v", gotRows)
	}
	if r.odbcRecordset != mockRS {
		t.Errorf("recordset not populated on Rows. got %v", gotRows)
	}

}

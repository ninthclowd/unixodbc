package odbc

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/ninthclowd/unixodbc/internal/api"
	"strings"
	"testing"
	"time"
	"unicode/utf16"
)

func newTestStatement(t *testing.T) (stmt *Statement, ctrl *gomock.Controller, mockAPI *MockAPI, mockHnd *MockHandle) {
	ctrl = gomock.NewController(t)
	mockAPI = NewMockAPI(ctrl)
	mockHnd = newTestHandle(ctrl, mockAPI)
	stmt = &Statement{
		handle: mockHnd,
		conn:   nil,
		rs:     nil,
	}
	return
}

func TestStatement_Close(t *testing.T) {
	stmt, ctrl, mockAPI, mockHnd := newTestStatement(t)
	defer ctrl.Finish()

	mockAPI.EXPECT().SQLFreeStmt(api.SQLHSTMT(mockHnd.hnd()), api.SQLUSMALLINT(api.SQL_CLOSE)).Return(api.SQLRETURN(api.SQL_SUCCESS))
	mockHnd.EXPECT().free().Return(nil)

	if gotErr := stmt.Close(); gotErr != nil {
		t.Fatalf("expected no error but got: %v", gotErr)
	}
}

func TestStatement_Close_Err(t *testing.T) {
	stmt, ctrl, mockAPI, _ := newTestStatement(t)
	defer ctrl.Finish()

	mockAPI.EXPECT().SQLFreeStmt(gomock.Any(), gomock.Any()).Return(api.SQLRETURN(api.SQL_ERROR))

	if gotErr := stmt.Close(); gotErr == nil {
		t.Fatal("expected an error but received none")
	}
}

func TestStatement_ResetParams(t *testing.T) {
	stmt, ctrl, mockAPI, mockHnd := newTestStatement(t)
	defer ctrl.Finish()

	mockAPI.EXPECT().SQLFreeStmt(api.SQLHSTMT(mockHnd.hnd()), api.SQLUSMALLINT(api.SQL_RESET_PARAMS)).Return(api.SQLRETURN(api.SQL_SUCCESS))

	if gotErr := stmt.ResetParams(); gotErr != nil {
		t.Fatalf("expected no error but got: %v", gotErr)
	}
}

func TestStatement_ResetParams_Err(t *testing.T) {
	stmt, ctrl, mockAPI, _ := newTestStatement(t)
	defer ctrl.Finish()

	mockAPI.EXPECT().SQLFreeStmt(gomock.Any(), gomock.Any()).Return(api.SQLRETURN(api.SQL_ERROR))

	if gotErr := stmt.ResetParams(); gotErr == nil {
		t.Fatal("expected an error but received none")
	}
}

func TestStatement_NumParams(t *testing.T) {
	stmt, ctrl, mockAPI, mockHnd := newTestStatement(t)
	defer ctrl.Finish()

	wantNum := 2

	mockAPI.EXPECT().SQLNumParams(api.SQLHSTMT(mockHnd.hnd()), gomock.Any()).DoAndReturn(func(hnd api.SQLHSTMT, paramCountPtr *api.SQLSMALLINT) api.SQLRETURN {
		*paramCountPtr = api.SQLSMALLINT(wantNum)
		return api.SQL_SUCCESS
	})

	if gotNum, gotErr := stmt.NumParams(); gotErr != nil {
		t.Fatalf("expected no error but got: %v", gotErr)
	} else if gotNum != wantNum {
		t.Fatalf("expected params to be %d but got %d", wantNum, gotNum)
	}
}

func TestStatement_NumParams_Err(t *testing.T) {
	stmt, ctrl, mockAPI, _ := newTestStatement(t)
	defer ctrl.Finish()

	mockAPI.EXPECT().SQLNumParams(gomock.Any(), gomock.Any()).Return(api.SQLRETURN(api.SQL_ERROR))
	if _, gotErr := stmt.NumParams(); gotErr == nil {
		t.Fatal("expected an error but received none")
	}
}

func TestStatement_Execute(t *testing.T) {
	stmt, ctrl, mockAPI, mockHnd := newTestStatement(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	mockAPI.EXPECT().SQLExecute(api.SQLHSTMT(mockHnd.hnd())).Return(api.SQLRETURN(api.SQL_SUCCESS))

	if gotErr := stmt.Execute(ctx); gotErr != nil {
		t.Fatalf("expected no error but got: %v", gotErr)
	}
}

func TestStatement_Execute_Cancel(t *testing.T) {
	stmt, ctrl, mockAPI, mockHnd := newTestStatement(t)
	defer ctrl.Finish()

	mockAPI.EXPECT().SQLExecute(api.SQLHSTMT(mockHnd.hnd())).DoAndReturn(func(hnd api.SQLHSTMT) api.SQLRETURN {
		time.Sleep(10 * time.Millisecond)
		return api.SQL_SUCCESS
	})

	mockHnd.EXPECT().cancel().Return(nil).Times(1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	start := time.Now()
	if gotErr := stmt.Execute(ctx); !strings.Contains(gotErr.Error(), context.DeadlineExceeded.Error()) {
		t.Fatalf("expected context cancellation but received: %v", gotErr)
	}
	if elapsed := time.Since(start); elapsed.Milliseconds() > 8 {
		t.Fatalf("cancelled context did not return immediately")
	}
}

func TestStatement_ExecDirect(t *testing.T) {
	stmt, ctrl, mockAPI, mockHnd := newTestStatement(t)
	defer ctrl.Finish()

	sql := "SELECT 1"
	unicodeSQL := utf16.Encode([]rune(sql))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	mockAPI.EXPECT().SQLExecDirect(api.SQLHSTMT(mockHnd.hnd()), unicodeSQL, api.SQLINTEGER(len(sql))).Return(api.SQLRETURN(api.SQL_SUCCESS))

	if gotErr := stmt.ExecDirect(ctx, sql); gotErr != nil {
		t.Fatalf("expected no error but got: %v", gotErr)
	}
}

func TestStatement_ExecDirect_Cancel(t *testing.T) {
	stmt, ctrl, mockAPI, mockHnd := newTestStatement(t)
	defer ctrl.Finish()

	sql := "SELECT 1"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	mockAPI.EXPECT().SQLExecDirect(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(hnd api.SQLHSTMT, sql []uint16, l api.SQLINTEGER) api.SQLRETURN {
		time.Sleep(10 * time.Millisecond)
		return api.SQL_SUCCESS
	})

	mockHnd.EXPECT().cancel().Return(nil).Times(1)

	start := time.Now()
	if gotErr := stmt.ExecDirect(ctx, sql); !strings.Contains(gotErr.Error(), context.DeadlineExceeded.Error()) {
		t.Fatalf("expected context cancellation but received: %v", gotErr)
	}
	if elapsed := time.Since(start); elapsed.Milliseconds() > 8 {
		t.Fatalf("cancelled context did not return immediately")
	}
}

func TestStatement_Prepare(t *testing.T) {
	stmt, ctrl, mockAPI, mockHnd := newTestStatement(t)
	defer ctrl.Finish()

	sql := "SELECT 1"
	unicodeSQL := utf16.Encode([]rune(sql))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	mockAPI.EXPECT().SQLPrepare(api.SQLHSTMT(mockHnd.hnd()), unicodeSQL, api.SQLINTEGER(len(sql))).Return(api.SQLRETURN(api.SQL_SUCCESS))

	if gotErr := stmt.Prepare(ctx, sql); gotErr != nil {
		t.Fatalf("expected no error but got: %v", gotErr)
	}
}

func TestStatement_Prepare_Cancel(t *testing.T) {
	stmt, ctrl, mockAPI, mockHnd := newTestStatement(t)
	defer ctrl.Finish()

	sql := "SELECT 1"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	mockAPI.EXPECT().SQLPrepare(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(hnd api.SQLHSTMT, sql []uint16, l api.SQLINTEGER) api.SQLRETURN {
		time.Sleep(10 * time.Millisecond)
		return api.SQL_SUCCESS
	})

	mockHnd.EXPECT().cancel().Return(nil)

	start := time.Now()
	if gotErr := stmt.Prepare(ctx, sql); !strings.Contains(gotErr.Error(), context.DeadlineExceeded.Error()) {
		t.Fatalf("expected context cancellation but received: %v", gotErr)
	}
	if elapsed := time.Since(start); elapsed.Milliseconds() > 8 {
		t.Fatalf("cancelled context did not return immediately")
	}
}

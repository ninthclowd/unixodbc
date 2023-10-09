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

func newTestEnvironment(t *testing.T) (env *Environment, ctrl *gomock.Controller, mockAPI *MockAPI, mockHnd *MockHandle) {
	ctrl = gomock.NewController(t)
	mockAPI = NewMockAPI(ctrl)
	mockHnd = newTestHandle(ctrl, mockAPI)
	env = &Environment{handle: mockHnd}
	return
}

func TestEnvironment_Connect_Cancel(t *testing.T) {
	env, ctrl, mockAPI, mockHnd := newTestEnvironment(t)
	defer ctrl.Finish()

	mockChild := newTestHandle(ctrl, mockAPI)

	mockHnd.
		EXPECT().
		child(api.SQLSMALLINT(api.SQL_HANDLE_DBC)).
		Return(mockChild, nil)

	wantConnStr := "connection string"

	mockAPI.
		EXPECT().
		SQLDriverConnectW(
			(api.SQLHDBC)(mockChild.hnd()),
			api.SQLHWND(nil),
			utf16.Encode([]rune(wantConnStr)),
			api.SQLSMALLINT(len(wantConnStr)),
			nil,
			api.SQLSMALLINT(0),
			nil,
			api.SQLUSMALLINT(api.SQL_DRIVER_NOPROMPT),
		).
		DoAndReturn(func(connectionHandle api.SQLHDBC, windowHandle api.SQLHWND, inConnectionString []uint16, stringLength1 api.SQLSMALLINT, outConnectionString *[]uint16, bufferLength api.SQLSMALLINT, stringLength2Ptr *api.SQLSMALLINT, driverCompletion api.SQLUSMALLINT) api.SQLRETURN {
			time.Sleep(10 * time.Millisecond)
			return api.SQL_SUCCESS
		})

	mockChild.
		EXPECT().
		cancel().
		Return(nil).
		Times(1)

	mockChild.
		EXPECT().
		free().
		Return(nil).
		Times(1)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	start := time.Now()
	if _, gotErr := env.Connect(ctx, wantConnStr); !strings.Contains(gotErr.Error(), context.DeadlineExceeded.Error()) {
		t.Fatalf("expected context cancellation but received: %v", gotErr)
	}
	if elapsed := time.Since(start); elapsed.Milliseconds() > 8 {
		t.Fatalf("cancelled context did not return immediately")
	}
}

func TestEnvironment_Connect(t *testing.T) {
	env, ctrl, mockAPI, mockHnd := newTestEnvironment(t)
	defer ctrl.Finish()

	mockChild := newTestHandle(ctrl, mockAPI)

	mockHnd.
		EXPECT().
		child(api.SQLSMALLINT(api.SQL_HANDLE_DBC)).
		Return(mockChild, nil)

	wantConnStr := "connection string"

	mockAPI.
		EXPECT().
		SQLDriverConnectW(
			(api.SQLHDBC)(mockChild.hnd()),
			api.SQLHWND(nil),
			utf16.Encode([]rune(wantConnStr)),
			api.SQLSMALLINT(len(wantConnStr)),
			nil,
			api.SQLSMALLINT(0),
			nil,
			api.SQLUSMALLINT(api.SQL_DRIVER_NOPROMPT),
		).
		Return(api.SQLRETURN(api.SQL_SUCCESS))

	if gotConnection, gotErr := env.Connect(context.Background(), wantConnStr); gotErr != nil {
		t.Errorf("expected no error, got: %s", gotErr.Error())
	} else if gotConnection == nil {
		t.Errorf("connection was nil")
	}
}

func TestEnvironment_Close(t *testing.T) {
	env, ctrl, _, mockHnd := newTestEnvironment(t)
	defer ctrl.Finish()

	mockHnd.EXPECT().free().Times(1)

	if gotErr := env.Close(); gotErr != nil {
		t.Errorf("expected no error, got: %s", gotErr)
	}

}

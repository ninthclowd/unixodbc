package odbc

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/ninthclowd/unixodbc/internal/api"
)

func newTestHandle(ctrl *gomock.Controller, mockAPI *MockAPI) (mockHnd *MockHandle) {
	var hnd = api.SQLHANDLE(new(string))
	mockHnd = NewMockHandle(ctrl)
	mockHnd.EXPECT().api().Return(mockAPI).AnyTimes()
	mockHnd.EXPECT().result(api.SQLRETURN(api.SQL_SUCCESS)).Return(api.SQLRETURN(api.SQL_SUCCESS), nil).AnyTimes()
	mockHnd.EXPECT().result(api.SQLRETURN(api.SQL_ERROR)).Return(api.SQLRETURN(api.SQL_ERROR), errors.New("mock error")).AnyTimes()
	mockHnd.EXPECT().hnd().Return(hnd).AnyTimes()
	return
}

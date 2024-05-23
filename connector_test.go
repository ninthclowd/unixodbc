package unixodbc

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/ninthclowd/unixodbc/internal/mocks"
	"github.com/ninthclowd/unixodbc/internal/odbc"
	"testing"
)

func TestConnector_Connect(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEnv := mocks.NewMockEnvironment(ctrl)

	connString := "connString"
	connector := &Connector{
		ConnectionString:   StaticConnStr(connString),
		StatementCacheSize: 5,
		odbcEnvironment:    mockEnv,
	}

	ctx := context.Background()

	mockEnv.EXPECT().SetVersion(odbc.Version380).Return(nil).Times(1)
	mockEnv.EXPECT().SetPoolOption(odbc.PoolOff).Return(nil).Times(1)

	mockConn1 := mocks.NewMockConnection(ctrl)
	mockConn1.EXPECT().SetAutoCommit(true).Return(nil).Times(1)
	mockEnv.EXPECT().Connect(gomock.Any(), connString).Return(mockConn1, nil).Times(1)

	gotConn, err := connector.Connect(ctx)
	if err != nil {
		t.Fatalf("expected no error connecting, got %v", err)
	}
	c, ok := gotConn.(*Connection)
	if !ok {
		t.Fatalf("connection was unexpected, got %v", err)
	}

	if c.odbcConnection != mockConn1 {
		t.Errorf("connection reference was unexpected, got %v", c.odbcConnection)
	}
	if capacity := c.cachedStatements.Capacity(); capacity != 5 {
		t.Errorf("capacity not set on cache unexpected, got %v", capacity)
	}

	if connector.Driver() != driverInstance {
		t.Errorf("driver instance not set on connector, got %v", driverInstance)
	}

	mockConn2 := mocks.NewMockConnection(ctrl)
	mockConn2.EXPECT().SetAutoCommit(true).Return(nil).Times(1)
	mockEnv.EXPECT().Connect(gomock.Any(), connString).Return(mockConn2, nil).Times(1)

	gotConn2, err := connector.Connect(ctx)
	if err != nil {
		t.Fatalf("expected no error connecting a second time, got %v", err)
	}
	c, ok = gotConn2.(*Connection)
	if !ok {
		t.Fatalf("second connection was unexpected, got %v", err)
	}

	if c.odbcConnection != mockConn2 {
		t.Errorf("second connection reference was unexpected, got %v", c.odbcConnection)
	}

}

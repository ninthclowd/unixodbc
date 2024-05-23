package unixodbc

import (
	"database/sql"
	"github.com/ninthclowd/unixodbc/internal/odbc"
	"testing"
)

func TestOpenHandles(t *testing.T) {
	want := odbc.OpenHandles()
	got := OpenHandles()
	if got != want {
		t.Fatalf("OpenHandles()=%v, want %v", got, want)
	}
}

func TestSQL_Open(t *testing.T) {
	db, err := sql.Open("unixodbc", "MyDSN")
	if err != nil {
		t.Fatalf("expected no error opening, got %v", err.Error())
	}
	if db.Driver() != driverInstance {
		t.Fatalf("incorrect driver, got %v", db.Driver())
	}
}

func TestDriver_OpenConnector(t *testing.T) {
	db, err := driverInstance.OpenConnector("MyDSN")
	if err != nil {
		t.Fatalf("expected no error opening, got %v", err.Error())
	}
	if db.Driver() != driverInstance {
		t.Fatalf("incorrect driver, got %v", db.Driver())
	}
	connector, ok := db.(*Connector)
	if !ok {
		t.Fatalf("incorrect connector, got %v", db)
	}
	if connStr, _ := connector.ConnectionString.ConnectionString(); connStr != "MyDSN" {
		t.Errorf("incorrect connection string, got %v", connStr)
	}
	if connector.StatementCacheSize != DefaultCacheSize {
		t.Errorf("incorrect default cache size, got %v", connector.StatementCacheSize)
	}

}

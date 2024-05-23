package unixodbc

import "testing"

func TestResult(t *testing.T) {
	r := result{
		lastInsertId: 0,
		rowsAffected: 0,
	}
	if got, _ := r.LastInsertId(); got != r.lastInsertId {
		t.Errorf("r.LastInsertId() = %v, want %v", got, r.lastInsertId)
	}
	if got, _ := r.RowsAffected(); got != r.rowsAffected {
		t.Errorf("r.RowsAffected() = %v, want %v", got, r.rowsAffected)
	}
}

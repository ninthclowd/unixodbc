package unixodbc

import "database/sql"

var _ sql.Result = (*result)(nil)

type result struct {
	lastInsertId int64
	rowsAffected int64
}

// LastInsertId implements sql.Result
func (r *result) LastInsertId() (int64, error) {
	return r.lastInsertId, nil
}

// RowsAffected implements sql.Result
func (r *result) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

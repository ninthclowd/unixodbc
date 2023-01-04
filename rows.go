package unixodbc

import "C"
import (
	"database/sql/driver"
	"github.com/ninthclowd/unixodbc/internal/api"
	"io"
	"reflect"
)

var _ driver.Rows = (*rows)(nil)
var _ driver.RowsColumnTypeDatabaseTypeName = (*rows)(nil)
var _ driver.RowsColumnTypeLength = (*rows)(nil)
var _ driver.RowsColumnTypeNullable = (*rows)(nil)
var _ driver.RowsColumnTypeScanType = (*rows)(nil)

//var _ driver.RowsColumnTypePrecisionScale = (*rows)(nil)

type rows struct {
	stmt    *stmt
	columns api.Columns
}

// ColumnTypeScanType implements driver.RowsColumnTypeScanType
func (r *rows) ColumnTypeScanType(index int) reflect.Type {
	return r.columns[index].ScanType()
}

// ColumnTypeNullable implements driver.RowsColumnTypeNullable
func (r *rows) ColumnTypeNullable(index int) (nullable, ok bool) {
	return r.columns[index].Nullable()
}

// ColumnTypeLength implements driver.RowsColumnTypeLength
func (r *rows) ColumnTypeLength(index int) (length int64, ok bool) {
	return r.columns[index].TypeLength(), r.columns[index].VariableLengthType()
}

// ColumnTypeDatabaseTypeName implements driver.RowsColumnTypeDatabaseTypeName
func (r *rows) ColumnTypeDatabaseTypeName(index int) string {
	if typeInfo, err := r.columns[index].TypeInfo(); err == nil {
		return typeInfo.TypeName
	}
	return ""
}

// Columns implements driver.Rows
func (r *rows) Columns() []string {
	return r.columns.Columns()
}

// Close implements driver.Rows
func (r *rows) Close() error {
	defer r.stmt.resultSetLock.Unlock()
	return r.stmt.hnd.SQLCloseCursor()
}

// Next implements driver.Rows
func (r *rows) Next(dest []driver.Value) (err error) {
	more, err := r.stmt.hnd.SQLFetch()

	if err != nil {
		return
	}

	if !more {
		return io.EOF
	}

	for i, _ := range dest {
		if dest[i], err = r.columns[i].Decode(); err != nil {
			return err
		}
	}
	return
}

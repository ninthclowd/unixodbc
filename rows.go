package unixodbc

import "C"
import (
	"context"
	"database/sql/driver"
	"github.com/ninthclowd/unixodbc/internal/odbc"
	"io"
	"reflect"
)

var _ driver.Rows = (*Rows)(nil)

// var _ driver.RowsColumnTypeDatabaseTypeName = (*Rows)(nil)
var _ driver.RowsColumnTypeLength = (*Rows)(nil)
var _ driver.RowsColumnTypeNullable = (*Rows)(nil)
var _ driver.RowsColumnTypeScanType = (*Rows)(nil)
var _ driver.RowsColumnTypePrecisionScale = (*Rows)(nil)

type Rows struct {
	ctx                context.Context
	odbcRecordset      *odbc.RecordSet
	closeStmtOnRSClose *odbc.Statement
}

// ColumnTypePrecisionScale implements driver.RowsColumnTypePrecisionScale
func (r *Rows) ColumnTypePrecisionScale(index int) (precision, scale int64, ok bool) {
	return r.odbcRecordset.Column(index).Decimal()

}

//// ColumnTypeDatabaseTypeName implements RowsColumnTypeDatabaseTypeName
//func (r *Rows) ColumnTypeDatabaseTypeName(index int) string {
//	col := r.odbcRecordset.Column(index)
//	return col.TypeName
//}

// ColumnTypeScanType implements driver.RowsColumnTypeScanType
func (r *Rows) ColumnTypeScanType(index int) reflect.Type {
	return r.odbcRecordset.Column(index).ScanType()
}

// ColumnTypeNullable implements driver.RowsColumnTypeNullable
func (r *Rows) ColumnTypeNullable(index int) (nullable, ok bool) {
	return r.odbcRecordset.Column(index).Nullable()
}

// ColumnTypeLength implements driver.RowsColumnTypeLength
func (r *Rows) ColumnTypeLength(index int) (length int64, ok bool) {
	return r.odbcRecordset.Column(index).VariableLength()
}

// Columns implements driver.Rows
func (r *Rows) Columns() []string {
	return r.odbcRecordset.ColumnNames()
}

// Close implements driver.Rows
func (r *Rows) Close() error {
	errs := make(MultipleErrors)
	Tracer.WithRegion(r.ctx, "Rows::Close", func() {
		errs["closing recordset"] = r.odbcRecordset.Close()
	})
	r.odbcRecordset = nil
	if r.closeStmtOnRSClose != nil {
		errs["closing statement"] = r.closeStmtOnRSClose.Close()
	}
	return errs.Error()
}

// Next implements driver.Rows
func (r *Rows) Next(dest []driver.Value) error {
	var more bool
	var err error
	Tracer.WithRegion(r.ctx, "Fetching row", func() {
		more, err = r.odbcRecordset.Fetch()
	})
	if err != nil {
		return err
	}

	if !more {
		return io.EOF
	}

	errs := make(MultipleErrors)
	for i := range dest {
		col := r.odbcRecordset.Column(i)
		Tracer.WithRegion(r.ctx, "Scanning column "+col.Name(), func() {
			dest[i], errs[col.Name()] = col.Value()
		})
	}
	return errs.Error()
}

//TODO: is it possible to paginate with SQLExtendedFetch and the go sql driver

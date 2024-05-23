package unixodbc

import (
	"context"
	"database/sql/driver"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/ninthclowd/unixodbc/internal/mocks"
	"io"
	"reflect"
	"testing"
)

func testRows(t *testing.T) (ctrl *gomock.Controller, rows *Rows, mockRS *mocks.MockRecordSet, mockStmt *mocks.MockStatement, ctx context.Context) {
	ctrl = gomock.NewController(t)

	ctx = context.Background()
	mockRS = mocks.NewMockRecordSet(ctrl)
	mockStmt = mocks.NewMockStatement(ctrl)
	rows = &Rows{
		ctx:                ctx,
		odbcRecordset:      mockRS,
		closeStmtOnRSClose: mockStmt,
	}
	return
}

func TestRows_ColumnTypePrecisionScale(t *testing.T) {
	ctrl, rows, mockRS, _, _ := testRows(t)
	defer ctrl.Finish()

	wantPrecision := int64(1)
	wantOk := true
	wantScale := int64(10)

	mockCol := mocks.NewMockColumn(ctrl)
	mockCol.EXPECT().Decimal().Return(wantPrecision, wantScale, wantOk).Times(1)
	mockRS.EXPECT().Column(1).Return(mockCol)
	gotPrecision, gotScale, gotOk := rows.ColumnTypePrecisionScale(1)
	if gotPrecision != wantPrecision {
		t.Errorf("got precision %v, want %v", gotPrecision, wantPrecision)
	}
	if gotOk != wantOk {
		t.Errorf("got ok %v, want %v", gotOk, wantOk)
	}
	if gotScale != wantScale {
		t.Errorf("got scale %v, want %v", gotScale, wantScale)
	}

}

func TestRows_ColumnTypeScanType(t *testing.T) {
	ctrl, rows, mockRS, _, _ := testRows(t)
	defer ctrl.Finish()

	wantScanType := reflect.TypeOf("string")

	mockCol := mocks.NewMockColumn(ctrl)
	mockCol.EXPECT().ScanType().Return(wantScanType).Times(1)
	mockRS.EXPECT().Column(1).Return(mockCol)
	gotScanType := rows.ColumnTypeScanType(1)
	if gotScanType != wantScanType {
		t.Errorf("got scanType %v, want %v", gotScanType, wantScanType)
	}
}

func TestRows_ColumnTypeNullable(t *testing.T) {
	ctrl, rows, mockRS, _, _ := testRows(t)
	defer ctrl.Finish()

	wantNullable := true
	wantNullableOk := true

	mockCol := mocks.NewMockColumn(ctrl)
	mockCol.EXPECT().Nullable().Return(wantNullable, wantNullableOk).Times(1)
	mockRS.EXPECT().Column(1).Return(mockCol)
	gotNullable, gotNullableOk := rows.ColumnTypeNullable(1)
	if gotNullable != wantNullable {
		t.Errorf("got nullable %v, want %v", gotNullable, wantNullable)
	}
	if gotNullableOk != wantNullableOk {
		t.Errorf("got nullable ok %v, want %v", gotNullableOk, wantNullableOk)
	}
}

func TestRows_ColumnTypeLength(t *testing.T) {
	ctrl, rows, mockRS, _, _ := testRows(t)
	defer ctrl.Finish()

	wantLength := int64(10)
	wantLengthOk := true

	mockCol := mocks.NewMockColumn(ctrl)
	mockCol.EXPECT().VariableLength().Return(wantLength, wantLengthOk).Times(1)
	mockRS.EXPECT().Column(1).Return(mockCol)
	gotLength, gotLengthOk := rows.ColumnTypeLength(1)
	if gotLength != wantLength {
		t.Errorf("got length %v, want %v", gotLength, wantLength)
	}
	if gotLengthOk != wantLengthOk {
		t.Errorf("got length ok %v, want %v", gotLengthOk, wantLengthOk)
	}
}

func TestRows_Columns(t *testing.T) {
	ctrl, rows, mockRS, _, _ := testRows(t)
	defer ctrl.Finish()

	wantColumns := []string{"one", "two"}

	mockRS.EXPECT().ColumnNames().Return(wantColumns)
	gotColumns := rows.Columns()
	if !reflect.DeepEqual(gotColumns, wantColumns) {
		t.Errorf("got columns %v, want %v", gotColumns, wantColumns)
	}
}

func TestRows_Close(t *testing.T) {
	ctrl, rows, mockRS, mockStmt, _ := testRows(t)
	defer ctrl.Finish()

	mockRS.EXPECT().Close().Return(nil).Times(1)
	mockStmt.EXPECT().Close().Return(nil).Times(1)
	gotErr := rows.Close()
	if gotErr != nil {
		t.Errorf("expected no error, got %v", gotErr)
	}
}

func TestRows_Next(t *testing.T) {
	ctrl, rows, mockRS, _, _ := testRows(t)
	defer ctrl.Finish()

	cols := []string{"col1", "col2"}
	wantVal1 := int64(10)
	wantVal2 := "stringVal"

	mockCol1 := mocks.NewMockColumn(ctrl)
	mockCol1.EXPECT().Name().Return(cols[0]).AnyTimes()
	mockCol1.EXPECT().Value().Return(wantVal1, nil).AnyTimes()
	mockRS.EXPECT().Column(0).Return(mockCol1).AnyTimes()

	mockCol2 := mocks.NewMockColumn(ctrl)
	mockCol2.EXPECT().Name().Return(cols[1]).AnyTimes()
	mockCol2.EXPECT().Value().Return(wantVal2, nil).AnyTimes()
	mockRS.EXPECT().Column(1).Return(mockCol2).AnyTimes()

	mockRS.EXPECT().Fetch().Return(true, nil).Times(1)

	gotRow1 := make([]driver.Value, 2)
	gotRow1Err := rows.Next(gotRow1)
	if gotRow1Err != nil {
		t.Fatalf("expected no error, got %v", gotRow1Err)
	}
	if gotRow1[0] != wantVal1 {
		t.Errorf("got value 1 %v, want %v", gotRow1[0], wantVal1)
	}
	if gotRow1[1] != wantVal2 {
		t.Errorf("got value 2 %v, want %v", gotRow1[1], wantVal2)
	}

	mockRS.EXPECT().Fetch().Return(false, nil).Times(1)

	gotRow2Err := rows.Next([]driver.Value{})
	if !errors.Is(gotRow2Err, io.EOF) {
		t.Fatalf("expected EOF on second fetch got %v", gotRow2Err)
	}

}

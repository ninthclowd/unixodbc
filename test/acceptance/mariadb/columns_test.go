package acceptance_test

import (
	_ "github.com/ninthclowd/unixodbc"
	"reflect"
	"testing"
	"time"
)

func TestColumnBinding(t *testing.T) {

	tests := []struct {
		description          string
		sql                  string
		wantColumnName       string
		wantTypeName         string
		wantDatabaseTypeName string
		wantLengthOk         bool
		wantLength           int64
		wantNullableOk       bool
		wantNullable         bool
		wantDecimalOk        bool
		wantPrecision        int64
		wantScale            int64
		scanInto             interface{}
		assertScannedValue   func(t *testing.T, got interface{})
	}{
		{
			description:          "it should decode a DATE value",
			sql:                  `SELECT CAST(TIMESTAMP('2009-05-18') AS DATE) AS "COL_DATE"`,
			wantColumnName:       "COL_DATE",
			wantTypeName:         "COL_DATE",
			wantDatabaseTypeName: "",
			wantLengthOk:         false,
			wantNullableOk:       true,
			wantNullable:         true,
			wantDecimalOk:        false,
			scanInto:             new(time.Time),
			assertScannedValue: func(t *testing.T, got interface{}) {
				gotValue := *got.(*time.Time)
				if !gotValue.Equal(time.Date(2009, 05, 18, 0, 0, 0, 0, time.UTC)) {
					t.Errorf("unexpected scanned value. got: %s", gotValue.String())
				}
			},
		},
		{
			description:          "it should decode a DATETIME value",
			sql:                  `SELECT CAST(TIMESTAMP('2009-05-18 12:03:05') AS DATETIME) AS "COL_DATETIME"`,
			wantColumnName:       "COL_DATETIME",
			wantTypeName:         "COL_DATETIME",
			wantDatabaseTypeName: "",
			wantNullableOk:       true,
			wantNullable:         true,
			wantDecimalOk:        false,
			scanInto:             new(time.Time),
			assertScannedValue: func(t *testing.T, got interface{}) {
				gotValue := *got.(*time.Time)
				if !gotValue.Equal(time.Date(2009, 05, 18, 12, 3, 5, 0, time.UTC)) {
					t.Errorf("unexpected scanned value. got: %s", gotValue.String())
				}
			},
		},
		{
			description:          "it should decode a TIME value",
			sql:                  `SELECT CAST('01:23:45' AS TIME) AS "COL_TIME"`,
			wantColumnName:       "COL_TIME",
			wantTypeName:         "COL_TIME",
			wantDatabaseTypeName: "",
			wantNullableOk:       true,
			wantNullable:         true,
			wantDecimalOk:        false,
			scanInto:             new(time.Time),
			assertScannedValue: func(t *testing.T, got interface{}) {
				gotValue := *got.(*time.Time)
				if !gotValue.Equal(time.Date(1, 1, 1, 1, 23, 45, 0, time.UTC)) {
					t.Errorf("unexpected scanned value. got: %s", gotValue.String())
				}
			},
		},
		{
			description:          "it should decode a CHAR value",
			sql:                  `SELECT CAST('FOO' AS CHAR) AS "COL_CHAR"`,
			wantColumnName:       "COL_CHAR",
			wantTypeName:         "COL_CHAR",
			wantDatabaseTypeName: "",
			wantNullableOk:       true,
			wantNullable:         true,
			wantLengthOk:         true,
			wantLength:           3,
			wantDecimalOk:        false,
			scanInto:             new(string),
			assertScannedValue: func(t *testing.T, got interface{}) {
				gotValue := *got.(*string)
				if gotValue != "FOO" {
					t.Errorf("unexpected scanned value. got: %s", gotValue)
				}
			},
		},
		{
			description:          "it should decode a SIGNED value",
			sql:                  `SELECT CAST(35 as SIGNED) AS "COL_SIGNED"`,
			wantColumnName:       "COL_SIGNED",
			wantTypeName:         "COL_SIGNED",
			wantDatabaseTypeName: "",
			wantNullableOk:       true,
			wantNullable:         false,
			wantLengthOk:         false,
			wantDecimalOk:        false,
			scanInto:             new(int32),
			assertScannedValue: func(t *testing.T, got interface{}) {
				gotValue := *got.(*int32)
				if gotValue != 35 {
					t.Errorf("unexpected scanned value. got: %d", gotValue)
				}
			},
		},
		{
			description:          "it should decode a BOOLEAN value",
			sql:                  `SELECT TRUE AS "COL_BOOLEAN"`,
			wantColumnName:       "COL_BOOLEAN",
			wantTypeName:         "COL_BOOLEAN",
			wantDatabaseTypeName: "",
			wantNullableOk:       true,
			wantNullable:         false,
			wantLengthOk:         false,
			wantDecimalOk:        false,
			scanInto:             new(int32),
			assertScannedValue: func(t *testing.T, got interface{}) {
				gotValue := *got.(*int32)
				if gotValue != 1 {
					t.Errorf("unexpected scanned value. got: %d", gotValue)
				}
			},
		},
		{
			description:          "it should decode a BIT value",
			sql:                  `SELECT b'1' AS "COL_BIT"`,
			wantColumnName:       "COL_BIT",
			wantTypeName:         "COL_BIT",
			wantDatabaseTypeName: "",
			wantNullableOk:       true,
			wantNullable:         false,
			wantLengthOk:         true,
			wantLength:           1,
			wantDecimalOk:        false,
			scanInto:             new([]uint8),
			assertScannedValue: func(t *testing.T, got interface{}) {
				gotValue := *got.(*[]uint8)
				if gotValue[0] != 1 {
					t.Errorf("unexpected scanned value. got: %d", gotValue[0])
				}
			},
		},
		{
			description:          "it should decode an INTEGER value",
			sql:                  `SELECT CAST(35513 as INTEGER) AS "COL_INT"`,
			wantColumnName:       "COL_INT",
			wantTypeName:         "COL_INT",
			wantDatabaseTypeName: "",
			wantNullableOk:       true,
			wantNullable:         false,
			wantLengthOk:         false,
			wantDecimalOk:        false,
			scanInto:             new(int32),
			assertScannedValue: func(t *testing.T, got interface{}) {
				gotValue := *got.(*int32)
				if gotValue != 35513 {
					t.Errorf("unexpected scanned value. got: %d", gotValue)
				}
			},
		},
		{
			description:          "it should decode an UNSIGNED value",
			sql:                  `SELECT CAST(35 as UNSIGNED) AS "COL_UNSIGNED"`,
			wantColumnName:       "COL_UNSIGNED",
			wantTypeName:         "COL_UNSIGNED",
			wantDatabaseTypeName: "",
			wantNullableOk:       true,
			wantNullable:         false,
			wantLengthOk:         false,
			wantDecimalOk:        false,
			scanInto:             new(int32),
			assertScannedValue: func(t *testing.T, got interface{}) {
				gotValue := *got.(*int32)
				if gotValue != 35 {
					t.Errorf("unexpected scanned value. got: %d", gotValue)
				}
			},
		},
		{
			description:          "it should decode a DECIMAL value",
			sql:                  `SELECT CAST(12.34 as DECIMAL(10,2)) AS "COL_DEC"`,
			wantColumnName:       "COL_DEC",
			wantTypeName:         "COL_DEC",
			wantDatabaseTypeName: "",
			wantNullableOk:       true,
			wantNullable:         false,
			wantLengthOk:         false,
			wantDecimalOk:        true,
			wantPrecision:        10,
			wantScale:            2,
			scanInto:             new(float64),
			assertScannedValue: func(t *testing.T, got interface{}) {
				gotValue := *got.(*float64)
				if gotValue != 12.34 {
					t.Errorf("unexpected scanned value. got: %f", gotValue)
				}
			},
		},
		{
			description:          "it should decode a DOUBLE value",
			sql:                  `SELECT CAST(543515121.33 as DOUBLE) AS "COL_DOUBLE"`,
			wantColumnName:       "COL_DOUBLE",
			wantTypeName:         "COL_DOUBLE",
			wantDatabaseTypeName: "",
			wantNullableOk:       true,
			wantNullable:         true,
			wantLengthOk:         false,
			wantDecimalOk:        true,
			wantPrecision:        15,
			wantScale:            0,
			scanInto:             new(float64),
			assertScannedValue: func(t *testing.T, got interface{}) {
				gotValue := *got.(*float64)
				if gotValue != 543515121.33 {
					t.Errorf("unexpected scanned value. got: %f", gotValue)
				}
			},
		},
		{
			description:          "it should decode a FLOAT value",
			sql:                  `SELECT CAST(1234.15 as FLOAT) AS "COL_FLOAT"`,
			wantColumnName:       "COL_FLOAT",
			wantTypeName:         "COL_FLOAT",
			wantDatabaseTypeName: "",
			wantNullableOk:       true,
			wantNullable:         true,
			wantLengthOk:         false,
			wantDecimalOk:        true,
			wantPrecision:        7,
			wantScale:            0,
			scanInto:             new(float32),
			assertScannedValue: func(t *testing.T, got interface{}) {
				gotValue := *got.(*float32)
				if gotValue != 1234.15 {
					t.Errorf("unexpected scanned value. got: %f", gotValue)
				}
			},
		},
		{
			description:          "it should decode a BINARY value",
			sql:                  `SELECT CAST('BIN' as BINARY) AS "COL_BINARY"`,
			wantColumnName:       "COL_BINARY",
			wantTypeName:         "COL_BINARY",
			wantDatabaseTypeName: "",
			wantNullableOk:       true,
			wantNullable:         true,
			wantLengthOk:         true,
			wantLength:           12,
			wantDecimalOk:        false,
			scanInto:             new([]uint8),
			assertScannedValue: func(t *testing.T, got interface{}) {
				gotValue := string(*got.(*[]uint8))
				if gotValue != "BIN" {
					t.Errorf("unexpected scanned value. got: %s", gotValue)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {

			_, conn, ctx, finish := newTestConnection(t)
			defer finish()

			rows, err := conn.QueryContext(ctx, test.sql)
			if err != nil {
				t.Fatalf("expected no error from query, got: %s", err.Error())
			}
			cols, err := rows.Columns()
			if err != nil {
				t.Fatalf("expected no error from columns, got: %s", err.Error())
			}
			if len(cols) != 1 {
				t.Errorf("expected one column returned, got: %+v", cols)
			}
			gotCol := cols[0]
			if gotCol != test.wantColumnName {
				t.Errorf("column was not expected. Got: %s", gotCol)
			}

			types, err := rows.ColumnTypes()
			if err != nil {
				t.Errorf("expected no error from columns, got: %s", err.Error())
			}

			if len(types) != 1 {
				t.Errorf("expected one type returned, got: %+v", types)
			}

			gotType := types[0]

			if gotType.Name() != test.wantTypeName {
				t.Errorf("unexpected type Name, got: %s", gotType.Name())
			}

			if gotType.DatabaseTypeName() != test.wantDatabaseTypeName {
				t.Errorf("unexpected DatabaseTypeName, got: %s", gotType.DatabaseTypeName())
			}

			if gotType.ScanType() != reflect.TypeOf(test.scanInto) {
				t.Errorf("unexpected ScanType name, got: %s", gotType.ScanType())
			}

			gotLength, gotLengthOk := gotType.Length()

			if gotLengthOk != test.wantLengthOk {
				t.Errorf("unexpected Length ok, got: %t", gotLengthOk)
			}

			if test.wantLengthOk {
				if gotLength != test.wantLength {
					t.Errorf("unexpected Length, got: %d", gotLength)
				}
			}

			gotNullable, gotNullableOk := gotType.Nullable()
			if gotNullableOk != test.wantNullableOk {
				t.Errorf("unexpected Nullable ok, got: %t", gotNullableOk)
			}
			if test.wantNullableOk {
				if gotNullable != test.wantNullable {
					t.Errorf("unexpected Nullable, got: %t", gotNullable)
				}
			}

			gotPrecision, gotScale, gotDecimalSizeOk := gotType.DecimalSize()
			if gotDecimalSizeOk != test.wantDecimalOk {
				t.Errorf("unexpected DecimalSize ok")
			}
			if test.wantDecimalOk {

				if gotPrecision != test.wantPrecision {
					t.Errorf("unexpected DecimalSize precision, got: %d", gotPrecision)
				}

				if gotScale != test.wantScale {
					t.Errorf("unexpected DecimalSize scale, got: %d", gotScale)
				}
			}

			if !rows.Next() {
				t.Fatalf("recordset was empty")
			}

			got := test.scanInto
			err = rows.Scan(got)
			if err != nil {
				t.Fatalf("expected no error from scan, got: %s", err.Error())
			}

			test.assertScannedValue(t, got)

			if rows.Next() {
				t.Fatalf("expected no more results from the recordset")
			}

		})

	}

}

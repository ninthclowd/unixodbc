package acceptance_test

import (
	"fmt"
	_ "github.com/ninthclowd/unixodbc"
	"math"
	"testing"
	"time"
)

func TestParameterBinding(t *testing.T) {

	tests := []struct {
		description string
		sql         string
		arg         interface{}
	}{
		{
			description: "it should bind a int8",
			sql:         fmt.Sprintf("SELECT 1 WHERE ? = %v", math.MaxInt8),
			arg:         int8(math.MaxInt8),
		},
		{
			description: "it should bind a int16",
			sql:         fmt.Sprintf("SELECT 1 WHERE ? = %v", math.MaxInt16),
			arg:         int16(math.MaxInt16),
		},
		{
			description: "it should bind a int32",
			sql:         fmt.Sprintf("SELECT 1 WHERE ? = %v", math.MaxInt32),
			arg:         int32(math.MaxInt32),
		},
		{
			description: "it should bind a int64",
			sql:         fmt.Sprintf("SELECT 1 WHERE ? = %v", math.MaxInt64),
			arg:         int64(math.MaxInt64),
		},
		{
			description: "it should bind a float32",
			sql:         "SELECT 1 WHERE ? = 52.42",
			arg:         float32(52.42),
		},
		{
			description: "it should bind a float64",
			sql:         "SELECT 1 WHERE ? = 52.42",
			arg:         52.42,
		},
		{
			description: "it should bind a date",
			sql:         "SELECT 1 WHERE CAST('2022/01/02' AS DATETIME) = ?",
			arg:         time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			description: "it should bind a boolean",
			sql:         `SELECT 1 WHERE ? = TRUE`,
			arg:         true,
		},
		{
			description: "it should bind a string",
			sql:         `SELECT 1 WHERE ? = 'foo'`,
			arg:         "foo",
		},
		{
			description: "it should bind an empty string",
			sql:         `SELECT 1 WHERE ? = ''`,
			arg:         "",
		},
		{
			description: "it should bind a nil value",
			sql:         `SELECT 1 WHERE ? IS NULL`,
			arg:         nil,
		},
		//TODO: binary
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {

			_, conn, ctx, finish := newTestConnection(t)
			defer finish()

			rows, err := conn.QueryContext(ctx, test.sql, test.arg)
			if err != nil {
				t.Fatalf("expected no error from QueryContext, got: %s", err.Error())
			}

			defer rows.Close()

			if rows.Next() != true {
				t.Errorf("'%s':%v returned no results", test.sql, test.arg)
			}
		})

	}

}

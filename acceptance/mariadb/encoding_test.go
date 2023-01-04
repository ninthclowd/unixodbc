package mariadb_test

import (
	"context"
	"database/sql"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tommy351/goldga"
	"time"
)

var _ = Describe("Encoding", func() {

	type DecodeTest struct {
		SQL    string
		NewOut func() interface{}
	}

	DescribeTable("Decoding Tests",
		func(test DecodeTest) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			rows, err := conn.QueryContext(ctx, test.SQL)
			Expect(err).To(BeNil())
			cols, err := rows.Columns()
			Expect(err).To(BeNil())
			Expect(cols).To(goldga.Match(goldga.WithDescription("expected cols")))

			types, err := rows.ColumnTypes()
			Expect(err).To(BeNil())
			Expect(CTypes(types)).To(goldga.Match(goldga.WithDescription("expected types")))

			Expect(rows.Next()).To(BeTrue())

			out := test.NewOut()
			err = rows.Scan(out)
			Expect(err).To(BeNil())
			Expect(out).To(goldga.Match(goldga.WithDescription("expected result")))
			Expect(rows.Next()).To(BeFalse())
		},
		Entry("DATE", DecodeTest{
			SQL:    `SELECT CAST(TIMESTAMP('2009-05-18') AS DATE) AS "COL_DATE"`,
			NewOut: func() interface{} { return new(time.Time) },
		}),
		Entry("DATETIME", DecodeTest{
			SQL:    `SELECT  CAST(TIMESTAMP('2009-05-18 12:03:05') AS DATE) AS "COL_DATETIME"`,
			NewOut: func() interface{} { return new(time.Time) },
		}),
		Entry("CHAR", DecodeTest{
			SQL:    `SELECT CAST('FOO' as CHAR) AS "COL_CHAR"`,
			NewOut: func() interface{} { return new(string) },
		}),
		Entry("NCHAR", DecodeTest{
			SQL:    `SELECT CAST('FOO' as NCHAR) AS "COL_NCHAR"`,
			NewOut: func() interface{} { return new(string) },
		}),

		Entry("SIGNED", DecodeTest{
			SQL:    `SELECT CAST(35 as SIGNED) AS "COL_SIGNED"`,
			NewOut: func() interface{} { return new(int64) },
		}),
		Entry("UNSIGNED", DecodeTest{
			SQL:    `SELECT CAST(35 as UNSIGNED) AS "COL_UNSIGNED"`,
			NewOut: func() interface{} { return new(uint64) },
		}),
		Entry("BINARY", DecodeTest{
			SQL: `SELECT CAST('BIN' as BINARY) AS "COL_BINARY"`,
			NewOut: func() interface{} {
				b := make([]byte, 100)
				return &b
			},
		}),
	)
})

type CType struct {
	ScanType         string
	DatabaseTypeName string
	Name             string
	Nullable         bool
	NullKnown        bool

	Length        int64
	LengthOk      bool
	Precision     int64
	Scale         int64
	DecimalSizeOk bool
}

func CTypes(arr []*sql.ColumnType) []*CType {
	c := make([]*CType, len(arr))
	for i, columnType := range arr {
		nullable, nullOk := columnType.Nullable()
		length, lengthOk := columnType.Length()
		precision, scale, decimalSizeOk := columnType.DecimalSize()
		c[i] = &CType{
			ScanType:         columnType.ScanType().Elem().Name(),
			DatabaseTypeName: columnType.DatabaseTypeName(),
			Name:             columnType.Name(),
			Nullable:         nullable,
			NullKnown:        nullOk,
			Length:           length,
			LengthOk:         lengthOk,
			Precision:        precision,
			Scale:            scale,
			DecimalSizeOk:    decimalSizeOk,
		}
	}
	return c
}

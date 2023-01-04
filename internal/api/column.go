package api

import (
	"reflect"
)

type Decoder interface {
	ScanType() reflect.Type
	Decode(c *Column) (out interface{}, err error)
}

type Column struct {
	pos           int
	stmt          *Statement
	name          string
	dataType      SQLSMALLINT
	decimalDigits int64
	nullable      bool
	nullableKnown bool
	columnSize    int64
	decoder       Decoder
}

type Columns []*Column

func (c Columns) Map() map[string]*Column {
	m := make(map[string]*Column)
	for _, column := range c {
		m[column.name] = column
	}
	return m
}

func (c Columns) Columns() []string {
	names := make([]string, len(c))
	for i, column := range c {
		names[i] = column.name
	}
	return names
}

func (c *Column) ScanType() reflect.Type {
	return c.decoder.ScanType()
}

func (c *Column) Precision() int64 {
	return c.decimalDigits
}

func (c *Column) Nullable() (nullable, ok bool) {
	return c.nullable, c.nullableKnown
}

func (c *Column) TypeInfo() (*TypeInfo, error) {
	return c.stmt.conn.TypeInfo(c.dataType)
}

func (c *Column) TypeLength() (length int64) {
	return c.columnSize
}
func (c *Column) VariableLengthType() bool {
	return false //TODO
}

func (c *Column) Decode() (out interface{}, err error) {
	return c.decoder.Decode(c)
}

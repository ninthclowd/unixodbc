package api

import "C"
import (
	"database/sql/driver"
)

type Encoder interface {
	Encode(p *Param, value interface{}) error
}

type Params []*Param

type Param struct {
	encoder       Encoder
	stmt          *Statement
	pos           int
	dataType      SQLSMALLINT
	parameterSize SQLULEN
	decimalDigits SQLSMALLINT
	nullable      bool
}

func (p *Param) Bind(value driver.Value) error {
	return p.encoder.Encode(p, value)
}

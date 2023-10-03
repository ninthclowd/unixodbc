package odbc

import (
	"database/sql/driver"
	"fmt"
	"github.com/ninthclowd/unixodbc/internal/api"
	"reflect"
	"unicode/utf16"
)

//go:generate mockgen -source=column.go -package odbc -destination column_mock_test.go
type Column interface {
	driver.Valuer
	Name() string
	VariableLength() (length int64, ok bool)
	Nullable() (nullable bool, ok bool)
	ScanType() reflect.Type
	Decimal() (precision int64, scale int64, ok bool)
}

type columnInfo struct {
	columnNumber  api.SQLUSMALLINT
	columnName    string
	dataType      api.SQLSMALLINT
	columnSize    api.SQLULEN
	decimalDigits api.SQLSMALLINT
	nullable      api.SQLSMALLINT
}

func (i *columnInfo) Name() string {
	return i.columnName
}
func (i *columnInfo) Nullable() (nullable bool, ok bool) {
	return i.nullable == api.SQL_NULLABLE, i.nullable != api.SQL_NULLABLE_UNKNOWN
}

type columnFactory func(info *columnInfo, hnd handle) Column

var registeredColumnFactory = map[api.SQLSMALLINT]columnFactory{}

func registerColumnFactoryForType(column columnFactory, types ...api.SQLSMALLINT) {
	for _, sqlsmallint := range types {
		if _, ok := registeredColumnFactory[sqlsmallint]; ok {
			panic(fmt.Sprintf("type %d is already registered", sqlsmallint))
		} else {
			registeredColumnFactory[sqlsmallint] = column
		}
	}
}

func columnsForStatement(h handle, loadColumn columnLoaderFN) (*columnsDetails, error) {
	var columnCount api.SQLSMALLINT
	if _, err := h.result(h.api().SQLNumResultCols((api.SQLHSTMT)(h.hnd()), &columnCount)); err != nil {
		return nil, fmt.Errorf("getting column count: %w", err)
	}
	details := &columnsDetails{
		names:   make([]string, columnCount),
		byIndex: make([]Column, columnCount),
		byName:  make(map[string]Column, columnCount),
	}

	for i := 0; i < int(columnCount); i++ {
		col, err := loadColumn(i)
		if err != nil {
			return nil, err
		}
		details.byIndex[i] = col
		details.names[i] = col.Name()
		details.byName[col.Name()] = col
	}
	return details, nil
}

type columnLoaderFN func(i int) (Column, error)

func newColumnLoader(h handle) columnLoaderFN {
	return func(i int) (Column, error) {
		info := &columnInfo{columnNumber: api.SQLUSMALLINT(i + 1)}
		name := make([]uint16, 100)
		var nameLength api.SQLSMALLINT
		_, err := h.result(h.api().SQLDescribeCol((api.SQLHSTMT)(h.hnd()),
			info.columnNumber,
			&name,
			api.SQLSMALLINT(len(name)),
			&nameLength,
			&info.dataType,
			&info.columnSize,
			&info.decimalDigits,
			&info.nullable))
		if err != nil {
			return nil, fmt.Errorf("describing column: %w", err)
		}
		info.columnName = string(utf16.Decode(name[:nameLength]))

		buildColumn, found := registeredColumnFactory[info.dataType]
		if !found {
			return nil, fmt.Errorf("no factory registered for column [%s] type [%d]", info.columnName, info.dataType)
		}

		return buildColumn(info, h), nil
	}
}

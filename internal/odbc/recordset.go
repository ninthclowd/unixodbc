package odbc

import (
	"fmt"
	"reflect"
)

//go:generate mockgen -source=recordset.go -package mocks -destination ../mocks/recordset.go
type RecordSet interface {
	Close() error
	Statement() Statement
	Fetch() (more bool, err error)
	Unmarshal(out interface{}) error
	ColumnWithName(name string) Column
	Column(index int) Column
	ColumnNames() []string
}

var _ RecordSet = (*recordSet)(nil)

type recordSet struct {
	stmt    *statement
	columns *columnsDetails
}

func (rs *recordSet) Close() error {
	if err := rs.stmt.closeCursor(); err != nil {
		return fmt.Errorf("closing cursor: %w", err)
	}
	rs.stmt = nil
	rs.columns = nil
	return nil
}

func (rs *recordSet) Statement() Statement {
	return rs.stmt
}

func (rs *recordSet) Fetch() (more bool, err error) {
	return rs.stmt.fetch()
}

func (rs *recordSet) Unmarshal(out interface{}) error {
	configStruct, ok := out.(reflect.Value)
	if !ok {
		configStruct = reflect.ValueOf(out).Elem()
	}

	for i := 0; i < configStruct.NumField(); i++ {
		property := configStruct.Type().Field(i)
		mappedColumnName := property.Tag.Get("col_name")
		if mappedColumnName == "" {
			continue
		}
		col, found := rs.columns.byName[mappedColumnName]
		if !found {
			continue
		}

		rawDBValue, err := col.Value()
		if err != nil {
			return fmt.Errorf("mapping [%s]: %w", property.Name, err)
		}
		dbVal := reflect.ValueOf(rawDBValue)
		if !dbVal.IsValid() {
			continue
		}

		propVal := configStruct.Field(i)
		propType := propVal.Type()
		if dbVal.Type() != propType {
			if dbVal.CanConvert(propType) {
				dbVal = dbVal.Convert(propType)
			} else {
				return fmt.Errorf("mapping [%s]: cannot convert database type [%s] to [%s]",
					property.Name,
					dbVal.Type().Name(),
					propType.Name())
			}
		}
		propVal.Set(dbVal)
	}

	return nil
}

func (rs *recordSet) ColumnWithName(name string) Column {
	if col, ok := rs.columns.byName[name]; ok {
		return col
	}
	return nil
}

func (rs *recordSet) Column(index int) Column {
	return rs.columns.byIndex[index]
}

func (rs *recordSet) ColumnNames() []string {
	return rs.columns.names
}

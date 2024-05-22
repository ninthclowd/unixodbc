package odbc

import (
	"fmt"
	"reflect"
)

type RecordSet struct {
	stmt    *Statement
	columns *columnsDetails
}

func (rs *RecordSet) Close() error {
	if err := rs.stmt.closeCursor(); err != nil {
		return fmt.Errorf("closing cursor: %w", err)
	}
	rs.stmt = nil
	rs.columns = nil
	return nil
}

func (rs *RecordSet) Statement() *Statement {
	return rs.stmt
}

func (rs *RecordSet) Fetch() (more bool, err error) {
	return rs.stmt.fetch()
}

func (rs *RecordSet) Unmarshal(out interface{}) error {
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

func (rs *RecordSet) ColumnWithName(name string) Column {
	if col, ok := rs.columns.byName[name]; ok {
		return col
	}
	return nil
}

func (rs *RecordSet) Column(index int) Column {
	return rs.columns.byIndex[index]
}

func (rs *RecordSet) ColumnNames() []string {
	return rs.columns.names
}

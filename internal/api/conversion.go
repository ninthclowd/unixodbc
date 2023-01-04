package api

import "C"
import (
	"errors"
)

var (
	ErrInvalidDataConversion = errors.New("invalid data conversion")
)

//type dbValue struct {
//	valueType SQLSMALLINT
//	value     SQLPOINTER
//	valueSize SQLLEN
//	outSize   *SQLLEN
//}

//
//
//func toCharArray(value driver.Value) (*dbValue, error) {
//	u := func(str string) *dbValue {
//		s := WCHAR(utf16.Encode([]rune(str)))
//		return &dbValue{
//			valueType: SQL_C_WCHAR,
//			value:     SQLPOINTER(s.Address()),
//			valueSize: SQLLEN(len(s)),
//			outSize:   new(SQLLEN),
//		}
//	}
//	switch value.(type) {
//	case []byte:
//		return u(string(value.([]byte))), nil
//	case string:
//		return u(value.(string)), nil
//	case decodeInt64:
//		return u(strconv.Itoa(int(value.(decodeInt64)))), nil
//	case boolEncoding:
//		return u(strconv.FormatBool(value.(boolEncoding))), nil
//	case decodeTime.Time:
//		return u(value.(decodeTime.Time).Format(decodeTime.RFC3339)), nil
//	case decodeFloat64:
//		return u(strconv.FormatFloat(value.(decodeFloat64), 'E', -1, 64)), nil
//	default:
//		return nil, ErrInvalidDataConversion
//	}
//}
//
//func toBit(value driver.Value) (*dbValue, error) {
//	var b byte = 0
//	switch value.(type) {
//	case string:
//		bo, err := strconv.ParseBool(value.(string))
//		if err != nil {
//			return nil, err
//		}
//		if bo {
//			b = 1
//		}
//	case decodeInt64:
//		if value.(decodeInt64) == 1 {
//			b = 1
//		}
//	case boolEncoding:
//		if value.(boolEncoding) {
//			b = 1
//		}
//	default:
//		return nil, ErrInvalidDataConversion
//	}
//	return &dbValue{
//		valueType: SQL_C_BIT,
//		value:     SQLPOINTER(&b),
//		valueSize: 1,
//		outSize:   nil,
//	}, nil
////}
//
//var toReflectType = map[SQLSMALLINT]reflect.Type{
//	SQL_CHAR:         byteType,
//	SQL_VARCHAR:      stringType,
//	SQL_LONGVARCHAR:  stringType,
//	SQL_WCHAR:        stringType,
//	SQL_WVARCHAR:     stringType,
//	SQL_WLONGVARCHAR: stringType,
//	SQL_DECIMAL:      float64Type,
//	//SQL_NUMERIC: reflect.in	NUMERIC(p,s)	Signed, exact, numeric value with a precision p and scale s (1 <= p <= 15; s <= p).[4]
//	SQL_SMALLINT: int16Type,
//	SQL_INTEGER:  int32Type,
//	//SQL_REAL	REAL	Signed, approximate, numeric value with a decodeBinary precision 24 (zero or absolute value 10[-38] to 10[38]).
//	SQL_FLOAT:   float64Type,
//	SQL_DOUBLE:  float64Type,
//	SQL_BIT:     boolType,
//	SQL_TINYINT: int8Type,
//	SQL_BIGINT:  int64Type,
//	//SQL_BINARY	BINARY(n)	Binary data of fixed length n.[9]
//	//SQL_VARBINARY	VARBINARY(n)	Variable length decodeBinary data of maximum length n. The maximum is set by the user.[9]
//	//SQL_LONGVARBINARY	LONG VARBINARY	Variable length decodeBinary data. Maximum length is data source-dependent.[9]
//	SQL_TYPE_DATE: timeType, //[6]	DATE	Year, month, and day fields, conforming to the rules of the Gregorian calendar. (See Constraints of the Gregorian Calendar, later in this appendix.)
//	//SQL_TYPE_TIME: 	TIME(p)	Hour, minute, and second fields, with valid values for hours of 00 to 23, valid values for minutes of 00 to 59, and valid values for seconds of 00 to 61. Precision p indicates the seconds precision.
//	SQL_TYPE_TIMESTAMP: timeType, //	TIMESTAMP(p)	Year, month, day, hour, minute, and second fields, with valid values as defined for the DATE and TIME data types.
//	//SQL_TYPE_UTCDATETIME	UTCDATETIME	Year, month, day, hour, minute, second, utchour, and utcminute fields. The utchour and utcminute fields have 1/10 microsecond precision.
//	//SQL_TYPE_UTCTIME	UTCTIME	Hour, minute, second, utchour, and utcminute fields. The utchour and utcminute fields have 1/10 microsecond precision..
//	//SQL_INTERVAL_MONTH[7]	INTERVAL MONTH(p)	Number of months between two dates; p is the interval leading precision.
//	//SQL_INTERVAL_YEAR[7]	INTERVAL YEAR(p)	Number of years between two dates; p is the interval leading precision.
//	//SQL_INTERVAL_YEAR_TO_MONTH[7]	INTERVAL YEAR(p) TO MONTH	Number of years and months between two dates; p is the interval leading precision.
//	//SQL_INTERVAL_DAY[7]	INTERVAL DAY(p)	Number of days between two dates; p is the interval leading precision.
//	//SQL_INTERVAL_HOUR[7]	INTERVAL HOUR(p)	Number of hours between two decodeDate/times; p is the interval leading precision.
//	//SQL_INTERVAL_MINUTE[7]	INTERVAL MINUTE(p)	Number of minutes between two decodeDate/times; p is the interval leading precision.
//	//SQL_INTERVAL_SECOND[7]	INTERVAL SECOND(p,q)	Number of seconds between two decodeDate/times; p is the interval leading precision and q is the interval seconds precision.
//	//SQL_INTERVAL_DAY_TO_HOUR[7]	INTERVAL DAY(p) TO HOUR	Number of days/hours between two decodeDate/times; p is the interval leading precision.
//	//SQL_INTERVAL_DAY_TO_MINUTE[7]	INTERVAL DAY(p) TO MINUTE	Number of days/hours/minutes between two decodeDate/times; p is the interval leading precision.
//	//SQL_INTERVAL_DAY_TO_SECOND[7]	INTERVAL DAY(p) TO SECOND(q)	Number of days/hours/minutes/seconds between two decodeDate/times; p is the interval leading precision and q is the interval seconds precision.
//	//SQL_INTERVAL_HOUR_TO_MINUTE[7]	INTERVAL HOUR(p) TO MINUTE	Number of hours/minutes between two decodeDate/times; p is the interval leading precision.
//	//SQL_INTERVAL_HOUR_TO_SECOND[7]	INTERVAL HOUR(p) TO SECOND(q)	Number of hours/minutes/seconds between two decodeDate/times; p is the interval leading precision and q is the interval seconds precision.
//	//SQL_INTERVAL_MINUTE_TO_SECOND[7]	INTERVAL MINUTE(p) TO SECOND(q)	Number of minutes/seconds between two decodeDate/times; p is the interval leading precision and q is the interval seconds precision.
//	//SQL_GUID	GUID	Fixed length GUID.
//
//}

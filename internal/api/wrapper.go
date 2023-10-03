package api

/*
#cgo linux LDFLAGS: -lodbc
#include "sql.h"
#include "sqlext.h"
#include "sqlucode.h"
#include "sqlspi.h"
#include <stdlib.h>
#include <stdint.h>

SQLPOINTER sqlLenDataAtExec(SQLLEN length) {
	return (SQLPOINTER)SQL_LEN_DATA_AT_EXEC(length);
}

SQLRETURN allocHandleNullInput(
      SQLSMALLINT   HandleType,
      SQLHANDLE *   OutputHandlePtr){
	return SQLAllocHandle(HandleType,SQL_NULL_HANDLE,OutputHandlePtr);
}

//SQLRETURN setEnvAttrConst(
//     SQLHENV      EnvironmentHandle,
//     SQLINTEGER   Attribute,
//     ulong ConstValue){
//	return SQLSetEnvAttr(EnvironmentHandle, Attribute, (SQLPOINTER)ConstValue, 0);
//}

SQLPOINTER constValue(ulong ConstValue){
	return (SQLPOINTER)ConstValue;
}

*/
import "C"
import (
	"unsafe"
)

//go:generate c-for-go -out ../  ./api.yml
type API struct{}

func (a *API) SQLEndTran(handleType SQLSMALLINT, handle SQLHANDLE, completionType SQLSMALLINT) SQLRETURN {
	return SQLRETURN(C.SQLEndTran(
		C.SQLSMALLINT(handleType),
		C.SQLHANDLE(handle),
		C.SQLSMALLINT(completionType)))
}

func (a *API) SQLExecute(statementHandle SQLHSTMT) SQLRETURN {
	return SQLRETURN(C.SQLExecute(C.SQLHSTMT(statementHandle)))
}

func (a *API) SQLExecDirect(statementHandle SQLHSTMT, statementText []uint16, textLength SQLINTEGER) SQLRETURN {
	return SQLRETURN(C.SQLExecDirectW(
		C.SQLHSTMT(statementHandle),
		(*C.SQLWCHAR)(unsafe.Pointer(&statementText[0])),
		C.SQLINTEGER(textLength),
	))
}

func (a *API) SQLGetData(statementHandle SQLHSTMT, colOrParamNum SQLUSMALLINT, targetType SQLSMALLINT, targetValuePtr SQLPOINTER, bufferLength SQLLEN, strLenOrIndPtr *SQLLEN) SQLRETURN {
	return SQLRETURN(C.SQLGetData(
		C.SQLHSTMT(statementHandle),
		C.SQLUSMALLINT(colOrParamNum),
		C.SQLSMALLINT(targetType),
		C.SQLPOINTER(targetValuePtr),
		C.SQLLEN(bufferLength),
		(*C.SQLLEN)(unsafe.Pointer(strLenOrIndPtr)),
	))
}

func (a *API) SQLGetTypeInfo(statementHandle SQLHSTMT, dataType SQLSMALLINT) SQLRETURN {
	return SQLRETURN(C.SQLGetTypeInfo(
		C.SQLHSTMT(statementHandle),
		C.SQLSMALLINT(dataType),
	))
}

func (a *API) SQLCancelHandle(handleType SQLSMALLINT, inputHandle SQLHANDLE) SQLRETURN {
	return SQLRETURN(C.SQLCancelHandle(
		C.SQLSMALLINT(handleType),
		C.SQLHANDLE(inputHandle),
	))
}

func (a *API) SQLBindParameter(statementHandle SQLHSTMT, parameterNumber SQLUSMALLINT, inputOutputType SQLSMALLINT, valueType SQLSMALLINT, parameterType SQLSMALLINT, columnSize SQLULEN, decimalDigits SQLSMALLINT, parameterValuePtr SQLPOINTER, bufferLength SQLLEN, strLenOrIndPtr *SQLLEN) SQLRETURN {
	return SQLRETURN(C.SQLBindParameter(
		C.SQLHSTMT(statementHandle),
		C.SQLUSMALLINT(parameterNumber),
		C.SQLSMALLINT(inputOutputType),
		C.SQLSMALLINT(valueType),
		C.SQLSMALLINT(parameterType),
		C.SQLULEN(columnSize),
		C.SQLSMALLINT(decimalDigits),
		C.SQLPOINTER(parameterValuePtr),
		C.SQLLEN(bufferLength),
		(*C.SQLLEN)(unsafe.Pointer(strLenOrIndPtr)),
	))
}

func (a *API) SQLDescribeCol(statementHandle SQLHSTMT, columnNumber SQLUSMALLINT, columnName *[]uint16, bufferLength SQLSMALLINT, nameLengthPtr *SQLSMALLINT, dataTypePtr *SQLSMALLINT, columnSizePtr *SQLULEN, decimalDigitsPtr *SQLSMALLINT, nullablePtr *SQLSMALLINT) SQLRETURN {
	return SQLRETURN(C.SQLDescribeColW(
		C.SQLHSTMT(statementHandle),
		C.SQLUSMALLINT(columnNumber),
		(*C.SQLWCHAR)(unsafe.Pointer(&(*columnName)[0])), //TODO unicode
		C.SQLSMALLINT(bufferLength),
		(*C.SQLSMALLINT)(unsafe.Pointer(nameLengthPtr)),
		(*C.SQLSMALLINT)(unsafe.Pointer(dataTypePtr)),
		(*C.SQLULEN)(unsafe.Pointer(columnSizePtr)),
		(*C.SQLSMALLINT)(unsafe.Pointer(decimalDigitsPtr)),
		(*C.SQLSMALLINT)(unsafe.Pointer(nullablePtr)),
	))
}

func (a *API) SQLSetEnvAttrConst(environmentHandle SQLHENV, attribute SQLINTEGER, value uint64) SQLRETURN {
	return SQLRETURN(C.SQLSetEnvAttr(C.SQLHENV(environmentHandle), C.SQLINTEGER(attribute), C.constValue(C.ulong(value)), SQL_IS_UINTEGER))
}

func (a *API) SQLSetEnvAttrStr(environmentHandle SQLHENV, attribute SQLINTEGER, value SQLPOINTER, stringLength SQLINTEGER) SQLRETURN {
	return SQLRETURN(C.SQLSetEnvAttr(C.SQLHENV(environmentHandle), C.SQLINTEGER(attribute), C.SQLPOINTER(value), C.SQLINTEGER(stringLength)))
}

func (a *API) SQLSetConnectAttrStr(connHandle SQLHDBC, attribute SQLINTEGER, value SQLPOINTER, stringLength SQLINTEGER) SQLRETURN {
	return SQLRETURN(C.SQLSetConnectAttr(C.SQLHDBC(connHandle), C.SQLINTEGER(attribute), C.SQLPOINTER(value), C.SQLINTEGER(stringLength)))
}

func (a *API) SQLSetConnectAttrConst(connHandle SQLHDBC, attribute SQLINTEGER, value uint64) SQLRETURN {
	return SQLRETURN(C.SQLSetConnectAttr(C.SQLHDBC(connHandle), C.SQLINTEGER(attribute), C.constValue(C.ulong(value)), SQL_IS_UINTEGER))
}

func (a *API) SQLDriverConnectW(connectionHandle SQLHDBC, windowHandle SQLHWND, inConnectionString []uint16, stringLength1 SQLSMALLINT, outConnectionString *[]uint16, bufferLength SQLSMALLINT, stringLength2Ptr *SQLSMALLINT, driverCompletion SQLUSMALLINT) SQLRETURN {
	var outStrPtr *uint16 = nil
	if outConnectionString != nil {
		outStrPtr = &(*outConnectionString)[0]
	}
	return SQLRETURN(C.SQLDriverConnectW(
		C.SQLHDBC(connectionHandle),
		C.SQLHWND(windowHandle),
		(*C.SQLWCHAR)(unsafe.Pointer(&inConnectionString[0])),
		C.SQLSMALLINT(stringLength1),
		(*C.SQLWCHAR)(unsafe.Pointer(outStrPtr)),
		C.SQLSMALLINT(bufferLength),
		(*C.SQLSMALLINT)(unsafe.Pointer(stringLength2Ptr)),
		C.SQLUSMALLINT(driverCompletion),
	))
}

func (a *API) SQLNumParams(statementHandle SQLHSTMT, parameterCountPtr *SQLSMALLINT) SQLRETURN {
	return SQLRETURN(C.SQLNumParams(
		C.SQLHSTMT(statementHandle),
		(*C.SQLSMALLINT)(unsafe.Pointer(parameterCountPtr)),
	))
}

func (a *API) SQLNumResultCols(statementHandle SQLHSTMT, columnCountPtr *SQLSMALLINT) SQLRETURN {
	return SQLRETURN(C.SQLNumResultCols(
		C.SQLHSTMT(statementHandle),
		(*C.SQLSMALLINT)(unsafe.Pointer(columnCountPtr)),
	))
}

func (a *API) SQLGetConnectAttr(connectionHandle SQLHDBC, attribute SQLINTEGER, valuePtr SQLPOINTER, bufferLength SQLINTEGER, stringLengthPtr *SQLINTEGER) SQLRETURN {
	return SQLRETURN(C.SQLGetConnectAttr(
		C.SQLHDBC(connectionHandle),
		C.SQLINTEGER(attribute),
		C.SQLPOINTER(valuePtr),
		C.SQLINTEGER(bufferLength),
		(*C.SQLINTEGER)(unsafe.Pointer(stringLengthPtr)),
	))
}

func (a *API) SQLPrepare(statementHandle SQLHSTMT, statementText []uint16, textLength SQLINTEGER) SQLRETURN {
	return SQLRETURN(C.SQLPrepareW(
		C.SQLHSTMT(statementHandle),
		(*C.SQLWCHAR)(unsafe.Pointer(&statementText[0])),
		C.SQLINTEGER(textLength),
	))
}

func (a *API) SQLFreeStmt(statementHandle SQLHSTMT, option SQLUSMALLINT) SQLRETURN {
	return SQLRETURN(C.SQLFreeStmt(
		C.SQLHSTMT(statementHandle),
		C.SQLUSMALLINT(option),
	))
}

func (a *API) SQLFetch(statementHandle SQLHSTMT) SQLRETURN {
	return SQLRETURN(C.SQLFetch(C.SQLHSTMT(statementHandle)))
}

func (a *API) SQLGetDiagRecW(handleType SQLSMALLINT, handle SQLHANDLE, recNumber SQLSMALLINT, sqlState *[]uint16, nativeErrorPtr *SQLINTEGER, messageText *[]uint16, bufferLength SQLSMALLINT, textLengthPtr *SQLSMALLINT) SQLRETURN {
	return SQLRETURN(C.SQLGetDiagRecW(
		C.SQLSMALLINT(handleType),
		C.SQLHANDLE(handle),
		C.SQLSMALLINT(recNumber),
		(*C.SQLWCHAR)(unsafe.Pointer(&(*sqlState)[0])),
		(*C.SQLINTEGER)(unsafe.Pointer(nativeErrorPtr)),
		(*C.SQLWCHAR)(unsafe.Pointer(&(*messageText)[0])),
		C.SQLSMALLINT(bufferLength),
		(*C.SQLSMALLINT)(unsafe.Pointer(textLengthPtr))))
}

func (a *API) SQLDisconnect(connectionHandle SQLHDBC) SQLRETURN {
	return SQLRETURN(C.SQLDisconnect(C.SQLHDBC(connectionHandle)))
}

func (a *API) SQLFreeHandle(handleType SQLSMALLINT, handle SQLHANDLE) SQLRETURN {
	return SQLRETURN(C.SQLFreeHandle(C.SQLSMALLINT(handleType), C.SQLHANDLE(handle)))
}

func (a *API) SQLCloseCursor(statementHandle SQLHSTMT) SQLRETURN {
	return SQLRETURN(C.SQLCloseCursor(C.SQLHSTMT(statementHandle)))
}

func (a *API) SQL_LEN_DATA_AT_EXEC(length SQLLEN) SQLPOINTER {
	return SQLPOINTER(C.sqlLenDataAtExec(C.SQLLEN(length)))
}

// SQLAllocHandle allocates an environment, connection, statement, or descriptor handle.
func (a *API) SQLAllocHandle(handleType SQLSMALLINT, inputHandle SQLHANDLE, outputHandle *SQLHANDLE) SQLRETURN {
	if inputHandle == nil {
		return SQLRETURN(C.allocHandleNullInput((C.SQLSMALLINT)(handleType), (*C.SQLHANDLE)(outputHandle)))
	}
	return SQLRETURN(C.SQLAllocHandle((C.SQLSMALLINT)(handleType), C.SQLHANDLE(inputHandle), (*C.SQLHANDLE)(outputHandle)))
}

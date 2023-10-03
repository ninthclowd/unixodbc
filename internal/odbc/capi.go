package odbc

import "github.com/ninthclowd/unixodbc/internal/api"

//go:generate mockgen -source=capi.go -package odbc -destination capi_mock_test.go -mock_names odbcAPI=MockAPI
type odbcAPI interface {
	SQLCancelHandle(handleType api.SQLSMALLINT, inputHandle api.SQLHANDLE) api.SQLRETURN
	SQLGetTypeInfo(statementHandle api.SQLHSTMT, dataType api.SQLSMALLINT) api.SQLRETURN
	SQLBindParameter(statementHandle api.SQLHSTMT, parameterNumber api.SQLUSMALLINT, inputOutputType api.SQLSMALLINT, valueType api.SQLSMALLINT, parameterType api.SQLSMALLINT, columnSize api.SQLULEN, decimalDigits api.SQLSMALLINT, parameterValuePtr api.SQLPOINTER, bufferLength api.SQLLEN, strLenOrIndPtr *api.SQLLEN) api.SQLRETURN
	SQLDescribeCol(statementHandle api.SQLHSTMT, columnNumber api.SQLUSMALLINT, columnName *[]uint16, bufferLength api.SQLSMALLINT, nameLengthPtr *api.SQLSMALLINT, dataTypePtr *api.SQLSMALLINT, columnSizePtr *api.SQLULEN, decimalDigitsPtr *api.SQLSMALLINT, nullablePtr *api.SQLSMALLINT) api.SQLRETURN
	SQLSetEnvAttrConst(environmentHandle api.SQLHENV, attribute api.SQLINTEGER, value uint64) api.SQLRETURN
	SQLSetConnectAttrConst(connectionHandle api.SQLHDBC, attribute api.SQLINTEGER, value uint64) api.SQLRETURN
	SQLDriverConnectW(connectionHandle api.SQLHDBC, windowHandle api.SQLHWND, inConnectionString []uint16, stringLength1 api.SQLSMALLINT, outConnectionString *[]uint16, bufferLength api.SQLSMALLINT, stringLength2Ptr *api.SQLSMALLINT, driverCompletion api.SQLUSMALLINT) api.SQLRETURN
	SQLPrepare(statementHandle api.SQLHSTMT, statementText []uint16, textLength api.SQLINTEGER) api.SQLRETURN
	SQLNumParams(statementHandle api.SQLHSTMT, parameterCountPtr *api.SQLSMALLINT) api.SQLRETURN
	SQLGetData(statementHandle api.SQLHSTMT, colOrParamNum api.SQLUSMALLINT, targetType api.SQLSMALLINT, targetValuePtr api.SQLPOINTER, bufferLength api.SQLLEN, strLenOrIndPtr *api.SQLLEN) api.SQLRETURN
	SQLNumResultCols(statementHandle api.SQLHSTMT, columnCountPtr *api.SQLSMALLINT) api.SQLRETURN
	SQLGetConnectAttr(connectionHandle api.SQLHDBC, attribute api.SQLINTEGER, valuePtr api.SQLPOINTER, bufferLength api.SQLINTEGER, stringLengthPtr *api.SQLINTEGER) api.SQLRETURN
	SQLFreeStmt(statementHandle api.SQLHSTMT, option api.SQLUSMALLINT) api.SQLRETURN
	SQLFetch(statementHandle api.SQLHSTMT) api.SQLRETURN
	SQLAllocHandle(handleType api.SQLSMALLINT, inputHandle api.SQLHANDLE, outputHandle *api.SQLHANDLE) api.SQLRETURN
	SQLGetDiagRecW(handleType api.SQLSMALLINT, handle api.SQLHANDLE, recNumber api.SQLSMALLINT, sqlState *[]uint16, nativeErrorPtr *api.SQLINTEGER, messageText *[]uint16, bufferLength api.SQLSMALLINT, textLengthPtr *api.SQLSMALLINT) api.SQLRETURN
	SQLDisconnect(connectionHandle api.SQLHDBC) api.SQLRETURN
	SQLFreeHandle(handleType api.SQLSMALLINT, handle api.SQLHANDLE) api.SQLRETURN
	SQLExecute(statementHandle api.SQLHSTMT) api.SQLRETURN
	SQLExecDirect(statementHandle api.SQLHSTMT, statementText []uint16, textLength api.SQLINTEGER) api.SQLRETURN
	SQLCloseCursor(statementHandle api.SQLHSTMT) api.SQLRETURN
	SQLSetEnvAttrStr(environmentHandle api.SQLHENV, attribute api.SQLINTEGER, value api.SQLPOINTER, stringLength api.SQLINTEGER) api.SQLRETURN
	SQLEndTran(handleType api.SQLSMALLINT, handle api.SQLHANDLE, completionType api.SQLSMALLINT) api.SQLRETURN
}

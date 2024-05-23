package api

/*
#cgo linux LDFLAGS: -lodbc
#include "sql.h"
#include "sqlext.h"
#include "sqlucode.h"
#include "sqlspi.h"
#include <stdlib.h>
#include <stdint.h>

SQLPOINTER* constValue(ulong ConstValue){
	return (SQLPOINTER*)ConstValue;
}

*/
import "C"

//go:generate c-for-go -out ../  ./api.yml

func Const(value uint64) *SQLPOINTER {
	return (*SQLPOINTER)(C.constValue(C.ulong(value)))
}

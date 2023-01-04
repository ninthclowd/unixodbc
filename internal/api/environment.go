package api

// #include <sql.h>
// #include <sqlext.h>
// #include <stdint.h>
/*
SQLPOINTER asPointer(uintptr_t valuePtr) {
	return (SQLPOINTER)valuePtr;
}
*/
import "C"
import (
	"fmt"
)

func NewEnvironment() (*Environment, error) {
	hnd := new(SQLHANDLE)
	if r := C.SQLAllocHandle(C.SQL_HANDLE_ENV, nil, (*C.SQLHANDLE)(hnd)); r == C.SQL_INVALID_HANDLE {
		return nil, ErrInvalidHandle
	} else if r == C.SQL_ERROR {
		return nil, fmt.Errorf("unable to alloc env Handle: %d", (int)(r))
	}
	return &Environment{

		Handle: &Handle{
			hType: C.SQL_HANDLE_ENV,
			hPtr:  C.SQLHANDLE(*hnd),
		},
	}, nil
}

type Environment struct {
	*Handle
}

func (env *Environment) SetAttrInt(attribute SQLINTEGER, value int32) error {
	if _, err := env.Result(C.SQLSetEnvAttr(env.hPtr, C.SQLINTEGER(attribute), C.asPointer(C.uintptr_t(uintptr(value))), 0)); err != nil {
		return fmt.Errorf("setting environment attribute [%d] to [%d]: %w", attribute, value, err)
	}
	return nil
}

func (env *Environment) Connection() (*Connection, error) {
	hnd := new(SQLHANDLE)
	if _, err := env.Result(C.SQLAllocHandle(C.SQL_HANDLE_DBC, env.hPtr, (*C.SQLHANDLE)(hnd))); err != nil {
		return nil, fmt.Errorf("alloc conn handle: %w", err)
	}
	return &Connection{
		typeInfo: make(map[SQLSMALLINT]*TypeInfo),
		Handle: &Handle{
			hType: C.SQL_HANDLE_DBC,
			hPtr:  C.SQLHANDLE(*hnd),
		},
	}, nil
}

package api

// #include <sql.h>
// #include <sqlext.h>
// #include <stdint.h>
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

var (
	ErrInvalidHandle = errors.New("invalid handle")
	ErrHandleFreed   = errors.New("attempt to double free")
	ErrClosed        = errors.New("connection closed")
	//ErrStillExecuting = errors.New("still executing")
	//ErrNoData         = errors.New("no data")
	//ErrNull           = errors.New("null data")
)

var _ error = (*SQLError)(nil)

type SQLError struct {
	diagRecords []*DiagRec
}

func (e *SQLError) Error() string {
	message := ""
	for i, record := range e.diagRecords {
		if i > 0 {
			message += "\n"
		}
		message += fmt.Sprintf("[%s:%d]", record.State, record.ErrorCode)
		if record.Message != "" {
			message += " " + record.Message
		}

	}
	return message
}

type DiagRec struct {
	State     string
	ErrorCode int
	Message   string
}

type Handle struct {
	hType C.SQLSMALLINT
	hPtr  C.SQLHANDLE
	freed bool
}

var errorMap = map[C.SQLRETURN]error{
	C.SQL_SUCCESS:           nil,
	C.SQL_SUCCESS_WITH_INFO: nil,
	C.SQL_NO_DATA:           nil,
	C.SQL_INVALID_HANDLE:    ErrInvalidHandle,
}

func (hnd *Handle) Result(ret C.SQLRETURN) (SQLRETURN, error) {
	r := (SQLRETURN)(ret)
	if hnd.freed {
		return r, ErrHandleFreed
	}
	if err, found := errorMap[ret]; found {
		return r, err
	}

	if d, err := hnd.getDiagRecs(); err != nil {
		return r, err
	} else {
		return r, &SQLError{diagRecords: d}
	}
}

func (hnd *Handle) getDiagRecs() ([]*DiagRec, error) {
	if hnd.freed {
		return nil, ErrHandleFreed
	}
	records := make([]*DiagRec, 0)
	for i := 1; ; i++ {
		var nativeError SQLINTEGER
		sqlState := WCHAR(make([]uint16, C.SQL_SQLSTATE_SIZE))
		messageText := WCHAR(make([]uint16, C.SQL_MAX_MESSAGE_LENGTH))
		var msgSize SQLSMALLINT
		recordNumber := C.SQLSMALLINT(i)

		ret := C.SQLGetDiagRecW(hnd.hType, hnd.hPtr,
			recordNumber,
			(*C.SQLWCHAR)(sqlState.Address()), (*C.SQLINTEGER)(unsafe.Pointer(&nativeError)),
			(*C.SQLWCHAR)(messageText.Address()), C.SQLSMALLINT(len(messageText)), (*C.SQLSMALLINT)(unsafe.Pointer(&msgSize)))

		if ret == C.SQL_NO_DATA {
			break
		}
		if err, _ := errorMap[ret]; err != nil {
			return nil, fmt.Errorf("getting diag record: %w", err)
		}
		if ret == C.SQL_ERROR {
			return nil, fmt.Errorf("SQLGetDiagRecW returned SQL_ERROR")
		}
		records = append(records, &DiagRec{
			State:     sqlState.String(),
			ErrorCode: int(nativeError),
			Message:   messageText[:msgSize].String(),
		})
	}
	return records, nil
}

func (hnd *Handle) Free() error {
	if hnd.freed {
		return ErrHandleFreed
	}
	if _, err := hnd.Result(C.SQLFreeHandle(hnd.hType, hnd.hPtr)); err != nil {
		return fmt.Errorf("freeing handle: %w", err)
	}
	hnd.freed = true
	return nil
}

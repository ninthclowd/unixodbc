package odbc

import (
	"errors"
	"fmt"
)

var (
	ErrNotImplemented = errors.New("not implemented")
	// ErrHandledAllocated = errors.New("handleImpl previously allocated")

	// ErrHandleFreed      = errors.New("attempt to double free")
	//ErrConnectionDead   = errors.New("Connection dead")
)

type MultipleErrors map[string]error

func (m MultipleErrors) Err() error {
	var combinedErrMsg string
	for desc, err := range m {
		if err != nil {
			//TODO in go1.20 use errors.join and implement unwrap
			if combinedErrMsg == "" {
				combinedErrMsg = fmt.Sprintf("%s: %s", desc, err.Error())
			} else {
				combinedErrMsg = fmt.Sprintf("\n%s: %s", desc, err.Error())
			}
		}
	}
	if combinedErrMsg != "" {
		return errors.New(combinedErrMsg)
	}
	return nil
}

type DiagRec struct {
	State     string
	ErrorCode int
	Message   string
}

var _ error = (*Error)(nil)

type Error struct {
	DiagRecords []*DiagRec
}

func (e *Error) Error() string {
	message := ""
	for i, record := range e.DiagRecords {
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

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

type ErrorMap map[string]error

func (m ErrorMap) Err() error {
	if len(m) == 0 {
		return nil
	}
	errs := make([]error, len(m))
	idx := 0
	for desc, err := range m {
		errs[idx] = fmt.Errorf("%s: %w", desc, err)
		idx++
	}
	return errors.Join(errs...)
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

package unixodbc

import (
	"errors"
	"fmt"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

type MultipleErrors map[string]error

func (m MultipleErrors) Error() error {
	if len(m) == 0 {
		return nil
	}
	var combinedErrMsg string
	for desc, err := range m {
		if err != nil {
			//TODO(ninthclowd): use errors.join and implement unwrap
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

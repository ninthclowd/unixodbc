//go:generate mockgen -source=handle.go -package odbc -destination handle_mock_test.go -mock_names handle=MockHandle
package odbc

import (
	"context"
	"errors"
	"fmt"
	"github.com/ninthclowd/unixodbc/internal/api"
	"runtime"
	"sync"
	"unicode/utf16"
)

func init() {
	runtime.LockOSThread()
}

var (
	ErrInvalidHandle = errors.New("invalid handleImpl")
	ErrHandleFreed   = errors.New("attempt to double free")
)

func cancelHandleOnContext(ctx context.Context, h handle) (done func()) {
	var wg sync.WaitGroup
	closer := make(chan struct{}, 1)

	done = func() {
		close(closer)
		wg.Wait()
	}

	wg.Add(1)
	go func() {
		select {
		case <-ctx.Done():
			_ = h.cancel()
		case <-closer:
		}
		wg.Done()
	}()

	return
}

var errorMap = map[api.SQLRETURN]error{
	api.SQL_SUCCESS:           nil,
	api.SQL_SUCCESS_WITH_INFO: nil,
	api.SQL_NO_DATA:           nil,
	api.SQL_INVALID_HANDLE:    ErrInvalidHandle,
}

type handle interface {
	result(r api.SQLRETURN) (api.SQLRETURN, error)
	hnd() api.SQLHANDLE
	hndType() api.SQLSMALLINT
	cancel() error
	free() error
	api() odbcAPI
	child(handleType api.SQLSMALLINT) (handle, error)
}

type handleImpl struct {
	cAPI    odbcAPI
	ptr     api.SQLHANDLE
	ptrType api.SQLSMALLINT
	ptrMux  sync.RWMutex
}

func newEnvHandle(cAPI odbcAPI) (handle, error) {
	hnd := &handleImpl{
		cAPI:    cAPI,
		ptrType: api.SQL_HANDLE_ENV,
	}

	if r := cAPI.SQLAllocHandle(api.SQL_HANDLE_ENV, nil, &hnd.ptr); r == api.SQL_ERROR {
		return nil, fmt.Errorf("unable to alloc env handle: %d", (int)(r))
	}
	return hnd, nil
}

func (h *handleImpl) api() odbcAPI {
	return h.cAPI
}

func (h *handleImpl) hnd() api.SQLHANDLE {
	h.ptrMux.RLock()
	defer h.ptrMux.RUnlock()
	return h.ptr
}

func (h *handleImpl) hndType() api.SQLSMALLINT {
	return h.ptrType
}

func (h *handleImpl) cancel() error {
	if code := h.cAPI.SQLCancelHandle(h.ptrType, h.hnd()); code != api.SQL_SUCCESS {
		return fmt.Errorf("received code %d when cancelling handle", code)
	}
	return nil
}

func (h *handleImpl) free() error {
	h.ptrMux.Lock()
	defer h.ptrMux.Unlock()
	if h.ptr == nil {
		return ErrHandleFreed
	}
	if code := h.cAPI.SQLFreeHandle(h.ptrType, h.ptr); code != api.SQL_SUCCESS {
		return fmt.Errorf("received code %d when freeing handle", code)
	}
	h.ptr = nil
	return nil
}

func (h *handleImpl) child(handleType api.SQLSMALLINT) (handle, error) {
	hnd := &handleImpl{cAPI: h.cAPI, ptrType: handleType}
	if _, err := h.result(h.cAPI.SQLAllocHandle(handleType, h.hnd(), &hnd.ptr)); err != nil {
		return nil, err
	}
	return hnd, nil
}

func (h *handleImpl) result(r api.SQLRETURN) (api.SQLRETURN, error) {
	if err, found := errorMap[r]; found {
		return r, err
	}
	if d, err := h.getDiagRecs(); err != nil {
		return r, err
	} else {
		return r, &Error{DiagRecords: d}
	}
}

func (h *handleImpl) getDiagRecs() ([]*DiagRec, error) {
	records := make([]*DiagRec, 0)
	for i := 1; ; i++ {
		var nativeError api.SQLINTEGER
		sqlState := make([]uint16, api.SQL_SQLSTATE_SIZE)
		messageText := make([]uint16, api.SQL_MAX_MESSAGE_LENGTH)

		var msgSize api.SQLSMALLINT
		recordNumber := api.SQLSMALLINT(i)

		ret := h.cAPI.SQLGetDiagRecW(h.ptrType, h.hnd(),
			recordNumber,
			&sqlState,
			&nativeError,
			&messageText,
			api.SQLSMALLINT(len(messageText)),
			&msgSize,
		)

		if ret == api.SQL_NO_DATA {
			break
		}
		if ret == api.SQL_ERROR {
			return nil, fmt.Errorf("SQLGetDiagRecW returned SQL_ERROR")
		}

		rec := &DiagRec{
			State:     string(utf16.Decode(sqlState)),
			ErrorCode: int(nativeError),
		}

		if msgSize == 0 {
			rec.Message = nullTerminatedString(utf16.Decode(messageText))
		} else {
			rec.Message = string(utf16.Decode(messageText[:msgSize]))
		}

		records = append(records, rec)
	}
	return records, nil
}

func nullTerminatedString(s []rune) string {
	for i, r := range s {
		if r == 0 {
			return string(s[:i])
		}
	}
	return string(s)
}

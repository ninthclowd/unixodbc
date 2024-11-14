package odbc

import (
	"context"
	"errors"
	"fmt"
	"github.com/ninthclowd/unixodbc/internal/api"
	"runtime"
	"sync"
	"sync/atomic"
	"unicode/utf16"
)

func init() {
	runtime.LockOSThread()
}

var (
	ErrInvalidHandle = errors.New("invalid handle")
	ErrHandleFreed   = errors.New("attempt to double free")
)

var openHandleCount atomic.Int64

func OpenHandles() int64 {
	return openHandleCount.Load()
}

func cancelHandleOnContext(ctx context.Context, h *handle) (done func()) {
	var wg sync.WaitGroup
	closer := make(chan struct{}, 1)

	done = func() {
		close(closer)
		wg.Wait()
	}

	wg.Add(1)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
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

type handle struct {
	ptr     *api.SQLHANDLE
	ptrType api.SQLSMALLINT
	ptrMux  sync.RWMutex
}

func newEnvHandle() (*handle, error) {
	hnd := &handle{
		ptrType: api.SQL_HANDLE_ENV,
	}

	if r := api.SQLAllocHandle(api.SQL_HANDLE_ENV, nil, &hnd.ptr); r == api.SQL_ERROR {
		return nil, fmt.Errorf("unable to alloc env handle: %d", (int)(r))
	}
	openHandleCount.Add(1)
	return hnd, nil
}

func (h *handle) hnd() *api.SQLHANDLE {
	h.ptrMux.RLock()
	defer h.ptrMux.RUnlock()
	return h.ptr
}

func (h *handle) hndType() api.SQLSMALLINT {
	return h.ptrType
}

func (h *handle) cancel() error {
	if code := api.SQLCancelHandle(h.ptrType, h.hnd()); code != api.SQL_SUCCESS {
		return fmt.Errorf("received code %d when cancelling handle", code)
	}
	return nil
}

func (h *handle) free() error {
	h.ptrMux.Lock()
	defer h.ptrMux.Unlock()
	if h.ptr == nil {
		return ErrHandleFreed
	}
	if code := api.SQLFreeHandle(h.ptrType, h.ptr); code != api.SQL_SUCCESS {
		return fmt.Errorf("received code %d when freeing handle", code)
	}
	openHandleCount.Add(-1)
	h.ptr = nil
	return nil
}

func (h *handle) child(handleType api.SQLSMALLINT) (*handle, error) {
	hnd := &handle{ptrType: handleType}
	if _, err := h.result(api.SQLAllocHandle(handleType, h.hnd(), &hnd.ptr)); err != nil {
		return nil, err
	}
	openHandleCount.Add(1)

	return hnd, nil
}

func (h *handle) result(r api.SQLRETURN) (api.SQLRETURN, error) {
	if err, found := errorMap[r]; found {
		return r, err
	}
	if d, err := h.getDiagRecs(); err != nil {
		return r, err
	} else {
		return r, &Error{DiagRecords: d}
	}
}

func (h *handle) getDiagRecs() ([]*DiagRec, error) {
	records := make([]*DiagRec, 0)
	for i := 1; ; i++ {
		var nativeError api.SQLINTEGER
		sqlState := make([]uint16, api.SQL_SQLSTATE_SIZE)
		messageText := make([]uint16, api.SQL_MAX_MESSAGE_LENGTH)

		var msgSize api.SQLSMALLINT
		recordNumber := api.SQLSMALLINT(i)

		ret := api.SQLGetDiagRecW(h.ptrType, h.hnd(),
			recordNumber,
			(*api.SQLWCHAR)(&sqlState[0]),
			&nativeError,
			(*api.SQLWCHAR)(&messageText[0]),
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

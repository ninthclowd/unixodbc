// THE AUTOGENERATED LICENSE. ALL THE RIGHTS ARE RESERVED BY ROBOTS.

// WARNING: This file has automatically been generated on Sat, 18 Feb 2023 19:13:59 MST.
// Code generated by https://git.io/c-for-go. DO NOT EDIT.

package api

/*
#cgo linux LDFLAGS: -lodbc
#include "sql.h"
#include "sqlext.h"
#include "stdint.h"
#include "sqlucode.h"
#include <stdlib.h>
#include "cgo_helpers.h"
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

// cgoAllocMap stores pointers to C allocated memory for future reference.
type cgoAllocMap struct {
	mux sync.RWMutex
	m   map[unsafe.Pointer]struct{}
}

var cgoAllocsUnknown = new(cgoAllocMap)

func (a *cgoAllocMap) Add(ptr unsafe.Pointer) {
	a.mux.Lock()
	if a.m == nil {
		a.m = make(map[unsafe.Pointer]struct{})
	}
	a.m[ptr] = struct{}{}
	a.mux.Unlock()
}

func (a *cgoAllocMap) IsEmpty() bool {
	a.mux.RLock()
	isEmpty := len(a.m) == 0
	a.mux.RUnlock()
	return isEmpty
}

func (a *cgoAllocMap) Borrow(b *cgoAllocMap) {
	if b == nil || b.IsEmpty() {
		return
	}
	b.mux.Lock()
	a.mux.Lock()
	for ptr := range b.m {
		if a.m == nil {
			a.m = make(map[unsafe.Pointer]struct{})
		}
		a.m[ptr] = struct{}{}
		delete(b.m, ptr)
	}
	a.mux.Unlock()
	b.mux.Unlock()
}

func (a *cgoAllocMap) Free() {
	a.mux.Lock()
	for ptr := range a.m {
		C.free(ptr)
		delete(a.m, ptr)
	}
	a.mux.Unlock()
}

// allocSQL_DATE_STRUCTMemory allocates memory for type C.SQL_DATE_STRUCT in C.
// The caller is responsible for freeing the this memory via C.free.
func allocSQL_DATE_STRUCTMemory(n int) unsafe.Pointer {
	mem, err := C.calloc(C.size_t(n), (C.size_t)(sizeOfSQL_DATE_STRUCTValue))
	if mem == nil {
		panic(fmt.Sprintln("memory alloc error: ", err))
	}
	return mem
}

const sizeOfSQL_DATE_STRUCTValue = unsafe.Sizeof([1]C.SQL_DATE_STRUCT{})

// Ref returns the underlying reference to C object or nil if struct is nil.
func (x *SQL_DATE_STRUCT) Ref() *C.SQL_DATE_STRUCT {
	if x == nil {
		return nil
	}
	return x.refecf4b83f
}

// Free invokes alloc map's free mechanism that cleanups any allocated memory using C free.
// Does nothing if struct is nil or has no allocation map.
func (x *SQL_DATE_STRUCT) Free() {
	if x != nil && x.allocsecf4b83f != nil {
		x.allocsecf4b83f.(*cgoAllocMap).Free()
		x.refecf4b83f = nil
	}
}

// NewSQL_DATE_STRUCTRef creates a new wrapper struct with underlying reference set to the original C object.
// Returns nil if the provided pointer to C object is nil too.
func NewSQL_DATE_STRUCTRef(ref unsafe.Pointer) *SQL_DATE_STRUCT {
	if ref == nil {
		return nil
	}
	obj := new(SQL_DATE_STRUCT)
	obj.refecf4b83f = (*C.SQL_DATE_STRUCT)(unsafe.Pointer(ref))
	return obj
}

// PassRef returns the underlying C object, otherwise it will allocate one and set its values
// from this wrapping struct, counting allocations into an allocation map.
func (x *SQL_DATE_STRUCT) PassRef() (*C.SQL_DATE_STRUCT, *cgoAllocMap) {
	if x == nil {
		return nil, nil
	} else if x.refecf4b83f != nil {
		return x.refecf4b83f, nil
	}
	memecf4b83f := allocSQL_DATE_STRUCTMemory(1)
	refecf4b83f := (*C.SQL_DATE_STRUCT)(memecf4b83f)
	allocsecf4b83f := new(cgoAllocMap)
	allocsecf4b83f.Add(memecf4b83f)

	var cyear_allocs *cgoAllocMap
	refecf4b83f.year, cyear_allocs = (C.SQLSMALLINT)(x.Year), cgoAllocsUnknown
	allocsecf4b83f.Borrow(cyear_allocs)

	var cmonth_allocs *cgoAllocMap
	refecf4b83f.month, cmonth_allocs = (C.SQLUSMALLINT)(x.Month), cgoAllocsUnknown
	allocsecf4b83f.Borrow(cmonth_allocs)

	var cday_allocs *cgoAllocMap
	refecf4b83f.day, cday_allocs = (C.SQLUSMALLINT)(x.Day), cgoAllocsUnknown
	allocsecf4b83f.Borrow(cday_allocs)

	x.refecf4b83f = refecf4b83f
	x.allocsecf4b83f = allocsecf4b83f
	return refecf4b83f, allocsecf4b83f

}

// PassValue does the same as PassRef except that it will try to dereference the returned pointer.
func (x SQL_DATE_STRUCT) PassValue() (C.SQL_DATE_STRUCT, *cgoAllocMap) {
	if x.refecf4b83f != nil {
		return *x.refecf4b83f, nil
	}
	ref, allocs := x.PassRef()
	return *ref, allocs
}

// Deref uses the underlying reference to C object and fills the wrapping struct with values.
// Do not forget to call this method whether you get a struct for C object and want to read its values.
func (x *SQL_DATE_STRUCT) Deref() {
	if x.refecf4b83f == nil {
		return
	}
	x.Year = (SQLSMALLINT)(x.refecf4b83f.year)
	x.Month = (SQLUSMALLINT)(x.refecf4b83f.month)
	x.Day = (SQLUSMALLINT)(x.refecf4b83f.day)
}

// allocSQL_TIME_STRUCTMemory allocates memory for type C.SQL_TIME_STRUCT in C.
// The caller is responsible for freeing the this memory via C.free.
func allocSQL_TIME_STRUCTMemory(n int) unsafe.Pointer {
	mem, err := C.calloc(C.size_t(n), (C.size_t)(sizeOfSQL_TIME_STRUCTValue))
	if mem == nil {
		panic(fmt.Sprintln("memory alloc error: ", err))
	}
	return mem
}

const sizeOfSQL_TIME_STRUCTValue = unsafe.Sizeof([1]C.SQL_TIME_STRUCT{})

// Ref returns the underlying reference to C object or nil if struct is nil.
func (x *SQL_TIME_STRUCT) Ref() *C.SQL_TIME_STRUCT {
	if x == nil {
		return nil
	}
	return x.ref51b6c82a
}

// Free invokes alloc map's free mechanism that cleanups any allocated memory using C free.
// Does nothing if struct is nil or has no allocation map.
func (x *SQL_TIME_STRUCT) Free() {
	if x != nil && x.allocs51b6c82a != nil {
		x.allocs51b6c82a.(*cgoAllocMap).Free()
		x.ref51b6c82a = nil
	}
}

// NewSQL_TIME_STRUCTRef creates a new wrapper struct with underlying reference set to the original C object.
// Returns nil if the provided pointer to C object is nil too.
func NewSQL_TIME_STRUCTRef(ref unsafe.Pointer) *SQL_TIME_STRUCT {
	if ref == nil {
		return nil
	}
	obj := new(SQL_TIME_STRUCT)
	obj.ref51b6c82a = (*C.SQL_TIME_STRUCT)(unsafe.Pointer(ref))
	return obj
}

// PassRef returns the underlying C object, otherwise it will allocate one and set its values
// from this wrapping struct, counting allocations into an allocation map.
func (x *SQL_TIME_STRUCT) PassRef() (*C.SQL_TIME_STRUCT, *cgoAllocMap) {
	if x == nil {
		return nil, nil
	} else if x.ref51b6c82a != nil {
		return x.ref51b6c82a, nil
	}
	mem51b6c82a := allocSQL_TIME_STRUCTMemory(1)
	ref51b6c82a := (*C.SQL_TIME_STRUCT)(mem51b6c82a)
	allocs51b6c82a := new(cgoAllocMap)
	allocs51b6c82a.Add(mem51b6c82a)

	var chour_allocs *cgoAllocMap
	ref51b6c82a.hour, chour_allocs = (C.SQLUSMALLINT)(x.Hour), cgoAllocsUnknown
	allocs51b6c82a.Borrow(chour_allocs)

	var cminute_allocs *cgoAllocMap
	ref51b6c82a.minute, cminute_allocs = (C.SQLUSMALLINT)(x.Minute), cgoAllocsUnknown
	allocs51b6c82a.Borrow(cminute_allocs)

	var csecond_allocs *cgoAllocMap
	ref51b6c82a.second, csecond_allocs = (C.SQLUSMALLINT)(x.Second), cgoAllocsUnknown
	allocs51b6c82a.Borrow(csecond_allocs)

	x.ref51b6c82a = ref51b6c82a
	x.allocs51b6c82a = allocs51b6c82a
	return ref51b6c82a, allocs51b6c82a

}

// PassValue does the same as PassRef except that it will try to dereference the returned pointer.
func (x SQL_TIME_STRUCT) PassValue() (C.SQL_TIME_STRUCT, *cgoAllocMap) {
	if x.ref51b6c82a != nil {
		return *x.ref51b6c82a, nil
	}
	ref, allocs := x.PassRef()
	return *ref, allocs
}

// Deref uses the underlying reference to C object and fills the wrapping struct with values.
// Do not forget to call this method whether you get a struct for C object and want to read its values.
func (x *SQL_TIME_STRUCT) Deref() {
	if x.ref51b6c82a == nil {
		return
	}
	x.Hour = (SQLUSMALLINT)(x.ref51b6c82a.hour)
	x.Minute = (SQLUSMALLINT)(x.ref51b6c82a.minute)
	x.Second = (SQLUSMALLINT)(x.ref51b6c82a.second)
}

// allocSQL_TIMESTAMP_STRUCTMemory allocates memory for type C.SQL_TIMESTAMP_STRUCT in C.
// The caller is responsible for freeing the this memory via C.free.
func allocSQL_TIMESTAMP_STRUCTMemory(n int) unsafe.Pointer {
	mem, err := C.calloc(C.size_t(n), (C.size_t)(sizeOfSQL_TIMESTAMP_STRUCTValue))
	if mem == nil {
		panic(fmt.Sprintln("memory alloc error: ", err))
	}
	return mem
}

const sizeOfSQL_TIMESTAMP_STRUCTValue = unsafe.Sizeof([1]C.SQL_TIMESTAMP_STRUCT{})

// Ref returns the underlying reference to C object or nil if struct is nil.
func (x *SQL_TIMESTAMP_STRUCT) Ref() *C.SQL_TIMESTAMP_STRUCT {
	if x == nil {
		return nil
	}
	return x.ref863f74dc
}

// Free invokes alloc map's free mechanism that cleanups any allocated memory using C free.
// Does nothing if struct is nil or has no allocation map.
func (x *SQL_TIMESTAMP_STRUCT) Free() {
	if x != nil && x.allocs863f74dc != nil {
		x.allocs863f74dc.(*cgoAllocMap).Free()
		x.ref863f74dc = nil
	}
}

// NewSQL_TIMESTAMP_STRUCTRef creates a new wrapper struct with underlying reference set to the original C object.
// Returns nil if the provided pointer to C object is nil too.
func NewSQL_TIMESTAMP_STRUCTRef(ref unsafe.Pointer) *SQL_TIMESTAMP_STRUCT {
	if ref == nil {
		return nil
	}
	obj := new(SQL_TIMESTAMP_STRUCT)
	obj.ref863f74dc = (*C.SQL_TIMESTAMP_STRUCT)(unsafe.Pointer(ref))
	return obj
}

// PassRef returns the underlying C object, otherwise it will allocate one and set its values
// from this wrapping struct, counting allocations into an allocation map.
func (x *SQL_TIMESTAMP_STRUCT) PassRef() (*C.SQL_TIMESTAMP_STRUCT, *cgoAllocMap) {
	if x == nil {
		return nil, nil
	} else if x.ref863f74dc != nil {
		return x.ref863f74dc, nil
	}
	mem863f74dc := allocSQL_TIMESTAMP_STRUCTMemory(1)
	ref863f74dc := (*C.SQL_TIMESTAMP_STRUCT)(mem863f74dc)
	allocs863f74dc := new(cgoAllocMap)
	allocs863f74dc.Add(mem863f74dc)

	var cyear_allocs *cgoAllocMap
	ref863f74dc.year, cyear_allocs = (C.SQLSMALLINT)(x.Year), cgoAllocsUnknown
	allocs863f74dc.Borrow(cyear_allocs)

	var cmonth_allocs *cgoAllocMap
	ref863f74dc.month, cmonth_allocs = (C.SQLUSMALLINT)(x.Month), cgoAllocsUnknown
	allocs863f74dc.Borrow(cmonth_allocs)

	var cday_allocs *cgoAllocMap
	ref863f74dc.day, cday_allocs = (C.SQLUSMALLINT)(x.Day), cgoAllocsUnknown
	allocs863f74dc.Borrow(cday_allocs)

	var chour_allocs *cgoAllocMap
	ref863f74dc.hour, chour_allocs = (C.SQLUSMALLINT)(x.Hour), cgoAllocsUnknown
	allocs863f74dc.Borrow(chour_allocs)

	var cminute_allocs *cgoAllocMap
	ref863f74dc.minute, cminute_allocs = (C.SQLUSMALLINT)(x.Minute), cgoAllocsUnknown
	allocs863f74dc.Borrow(cminute_allocs)

	var csecond_allocs *cgoAllocMap
	ref863f74dc.second, csecond_allocs = (C.SQLUSMALLINT)(x.Second), cgoAllocsUnknown
	allocs863f74dc.Borrow(csecond_allocs)

	var cfraction_allocs *cgoAllocMap
	ref863f74dc.fraction, cfraction_allocs = (C.SQLUINTEGER)(x.Fraction), cgoAllocsUnknown
	allocs863f74dc.Borrow(cfraction_allocs)

	x.ref863f74dc = ref863f74dc
	x.allocs863f74dc = allocs863f74dc
	return ref863f74dc, allocs863f74dc

}

// PassValue does the same as PassRef except that it will try to dereference the returned pointer.
func (x SQL_TIMESTAMP_STRUCT) PassValue() (C.SQL_TIMESTAMP_STRUCT, *cgoAllocMap) {
	if x.ref863f74dc != nil {
		return *x.ref863f74dc, nil
	}
	ref, allocs := x.PassRef()
	return *ref, allocs
}

// Deref uses the underlying reference to C object and fills the wrapping struct with values.
// Do not forget to call this method whether you get a struct for C object and want to read its values.
func (x *SQL_TIMESTAMP_STRUCT) Deref() {
	if x.ref863f74dc == nil {
		return
	}
	x.Year = (SQLSMALLINT)(x.ref863f74dc.year)
	x.Month = (SQLUSMALLINT)(x.ref863f74dc.month)
	x.Day = (SQLUSMALLINT)(x.ref863f74dc.day)
	x.Hour = (SQLUSMALLINT)(x.ref863f74dc.hour)
	x.Minute = (SQLUSMALLINT)(x.ref863f74dc.minute)
	x.Second = (SQLUSMALLINT)(x.ref863f74dc.second)
	x.Fraction = (SQLUINTEGER)(x.ref863f74dc.fraction)
}

// allocSQL_YEAR_MONTH_STRUCTMemory allocates memory for type C.SQL_YEAR_MONTH_STRUCT in C.
// The caller is responsible for freeing the this memory via C.free.
func allocSQL_YEAR_MONTH_STRUCTMemory(n int) unsafe.Pointer {
	mem, err := C.calloc(C.size_t(n), (C.size_t)(sizeOfSQL_YEAR_MONTH_STRUCTValue))
	if mem == nil {
		panic(fmt.Sprintln("memory alloc error: ", err))
	}
	return mem
}

const sizeOfSQL_YEAR_MONTH_STRUCTValue = unsafe.Sizeof([1]C.SQL_YEAR_MONTH_STRUCT{})

// Ref returns the underlying reference to C object or nil if struct is nil.
func (x *SQL_YEAR_MONTH_STRUCT) Ref() *C.SQL_YEAR_MONTH_STRUCT {
	if x == nil {
		return nil
	}
	return x.ref8e33cf55
}

// Free invokes alloc map's free mechanism that cleanups any allocated memory using C free.
// Does nothing if struct is nil or has no allocation map.
func (x *SQL_YEAR_MONTH_STRUCT) Free() {
	if x != nil && x.allocs8e33cf55 != nil {
		x.allocs8e33cf55.(*cgoAllocMap).Free()
		x.ref8e33cf55 = nil
	}
}

// NewSQL_YEAR_MONTH_STRUCTRef creates a new wrapper struct with underlying reference set to the original C object.
// Returns nil if the provided pointer to C object is nil too.
func NewSQL_YEAR_MONTH_STRUCTRef(ref unsafe.Pointer) *SQL_YEAR_MONTH_STRUCT {
	if ref == nil {
		return nil
	}
	obj := new(SQL_YEAR_MONTH_STRUCT)
	obj.ref8e33cf55 = (*C.SQL_YEAR_MONTH_STRUCT)(unsafe.Pointer(ref))
	return obj
}

// PassRef returns the underlying C object, otherwise it will allocate one and set its values
// from this wrapping struct, counting allocations into an allocation map.
func (x *SQL_YEAR_MONTH_STRUCT) PassRef() (*C.SQL_YEAR_MONTH_STRUCT, *cgoAllocMap) {
	if x == nil {
		return nil, nil
	} else if x.ref8e33cf55 != nil {
		return x.ref8e33cf55, nil
	}
	mem8e33cf55 := allocSQL_YEAR_MONTH_STRUCTMemory(1)
	ref8e33cf55 := (*C.SQL_YEAR_MONTH_STRUCT)(mem8e33cf55)
	allocs8e33cf55 := new(cgoAllocMap)
	allocs8e33cf55.Add(mem8e33cf55)

	var cyear_allocs *cgoAllocMap
	ref8e33cf55.year, cyear_allocs = (C.SQLUINTEGER)(x.Year), cgoAllocsUnknown
	allocs8e33cf55.Borrow(cyear_allocs)

	var cmonth_allocs *cgoAllocMap
	ref8e33cf55.month, cmonth_allocs = (C.SQLUINTEGER)(x.Month), cgoAllocsUnknown
	allocs8e33cf55.Borrow(cmonth_allocs)

	x.ref8e33cf55 = ref8e33cf55
	x.allocs8e33cf55 = allocs8e33cf55
	return ref8e33cf55, allocs8e33cf55

}

// PassValue does the same as PassRef except that it will try to dereference the returned pointer.
func (x SQL_YEAR_MONTH_STRUCT) PassValue() (C.SQL_YEAR_MONTH_STRUCT, *cgoAllocMap) {
	if x.ref8e33cf55 != nil {
		return *x.ref8e33cf55, nil
	}
	ref, allocs := x.PassRef()
	return *ref, allocs
}

// Deref uses the underlying reference to C object and fills the wrapping struct with values.
// Do not forget to call this method whether you get a struct for C object and want to read its values.
func (x *SQL_YEAR_MONTH_STRUCT) Deref() {
	if x.ref8e33cf55 == nil {
		return
	}
	x.Year = (SQLUINTEGER)(x.ref8e33cf55.year)
	x.Month = (SQLUINTEGER)(x.ref8e33cf55.month)
}

// allocSQL_DAY_SECOND_STRUCTMemory allocates memory for type C.SQL_DAY_SECOND_STRUCT in C.
// The caller is responsible for freeing the this memory via C.free.
func allocSQL_DAY_SECOND_STRUCTMemory(n int) unsafe.Pointer {
	mem, err := C.calloc(C.size_t(n), (C.size_t)(sizeOfSQL_DAY_SECOND_STRUCTValue))
	if mem == nil {
		panic(fmt.Sprintln("memory alloc error: ", err))
	}
	return mem
}

const sizeOfSQL_DAY_SECOND_STRUCTValue = unsafe.Sizeof([1]C.SQL_DAY_SECOND_STRUCT{})

// Ref returns the underlying reference to C object or nil if struct is nil.
func (x *SQL_DAY_SECOND_STRUCT) Ref() *C.SQL_DAY_SECOND_STRUCT {
	if x == nil {
		return nil
	}
	return x.ref74d37172
}

// Free invokes alloc map's free mechanism that cleanups any allocated memory using C free.
// Does nothing if struct is nil or has no allocation map.
func (x *SQL_DAY_SECOND_STRUCT) Free() {
	if x != nil && x.allocs74d37172 != nil {
		x.allocs74d37172.(*cgoAllocMap).Free()
		x.ref74d37172 = nil
	}
}

// NewSQL_DAY_SECOND_STRUCTRef creates a new wrapper struct with underlying reference set to the original C object.
// Returns nil if the provided pointer to C object is nil too.
func NewSQL_DAY_SECOND_STRUCTRef(ref unsafe.Pointer) *SQL_DAY_SECOND_STRUCT {
	if ref == nil {
		return nil
	}
	obj := new(SQL_DAY_SECOND_STRUCT)
	obj.ref74d37172 = (*C.SQL_DAY_SECOND_STRUCT)(unsafe.Pointer(ref))
	return obj
}

// PassRef returns the underlying C object, otherwise it will allocate one and set its values
// from this wrapping struct, counting allocations into an allocation map.
func (x *SQL_DAY_SECOND_STRUCT) PassRef() (*C.SQL_DAY_SECOND_STRUCT, *cgoAllocMap) {
	if x == nil {
		return nil, nil
	} else if x.ref74d37172 != nil {
		return x.ref74d37172, nil
	}
	mem74d37172 := allocSQL_DAY_SECOND_STRUCTMemory(1)
	ref74d37172 := (*C.SQL_DAY_SECOND_STRUCT)(mem74d37172)
	allocs74d37172 := new(cgoAllocMap)
	allocs74d37172.Add(mem74d37172)

	var cday_allocs *cgoAllocMap
	ref74d37172.day, cday_allocs = (C.SQLUINTEGER)(x.Day), cgoAllocsUnknown
	allocs74d37172.Borrow(cday_allocs)

	var chour_allocs *cgoAllocMap
	ref74d37172.hour, chour_allocs = (C.SQLUINTEGER)(x.Hour), cgoAllocsUnknown
	allocs74d37172.Borrow(chour_allocs)

	var cminute_allocs *cgoAllocMap
	ref74d37172.minute, cminute_allocs = (C.SQLUINTEGER)(x.Minute), cgoAllocsUnknown
	allocs74d37172.Borrow(cminute_allocs)

	var csecond_allocs *cgoAllocMap
	ref74d37172.second, csecond_allocs = (C.SQLUINTEGER)(x.Second), cgoAllocsUnknown
	allocs74d37172.Borrow(csecond_allocs)

	var cfraction_allocs *cgoAllocMap
	ref74d37172.fraction, cfraction_allocs = (C.SQLUINTEGER)(x.Fraction), cgoAllocsUnknown
	allocs74d37172.Borrow(cfraction_allocs)

	x.ref74d37172 = ref74d37172
	x.allocs74d37172 = allocs74d37172
	return ref74d37172, allocs74d37172

}

// PassValue does the same as PassRef except that it will try to dereference the returned pointer.
func (x SQL_DAY_SECOND_STRUCT) PassValue() (C.SQL_DAY_SECOND_STRUCT, *cgoAllocMap) {
	if x.ref74d37172 != nil {
		return *x.ref74d37172, nil
	}
	ref, allocs := x.PassRef()
	return *ref, allocs
}

// Deref uses the underlying reference to C object and fills the wrapping struct with values.
// Do not forget to call this method whether you get a struct for C object and want to read its values.
func (x *SQL_DAY_SECOND_STRUCT) Deref() {
	if x.ref74d37172 == nil {
		return
	}
	x.Day = (SQLUINTEGER)(x.ref74d37172.day)
	x.Hour = (SQLUINTEGER)(x.ref74d37172.hour)
	x.Minute = (SQLUINTEGER)(x.ref74d37172.minute)
	x.Second = (SQLUINTEGER)(x.ref74d37172.second)
	x.Fraction = (SQLUINTEGER)(x.ref74d37172.fraction)
}

// allocSQL_INTERVAL_STRUCTMemory allocates memory for type C.SQL_INTERVAL_STRUCT in C.
// The caller is responsible for freeing the this memory via C.free.
func allocSQL_INTERVAL_STRUCTMemory(n int) unsafe.Pointer {
	mem, err := C.calloc(C.size_t(n), (C.size_t)(sizeOfSQL_INTERVAL_STRUCTValue))
	if mem == nil {
		panic(fmt.Sprintln("memory alloc error: ", err))
	}
	return mem
}

const sizeOfSQL_INTERVAL_STRUCTValue = unsafe.Sizeof([1]C.SQL_INTERVAL_STRUCT{})

// Ref returns the underlying reference to C object or nil if struct is nil.
func (x *SQL_INTERVAL_STRUCT) Ref() *C.SQL_INTERVAL_STRUCT {
	if x == nil {
		return nil
	}
	return x.ref88eb89bb
}

// Free invokes alloc map's free mechanism that cleanups any allocated memory using C free.
// Does nothing if struct is nil or has no allocation map.
func (x *SQL_INTERVAL_STRUCT) Free() {
	if x != nil && x.allocs88eb89bb != nil {
		x.allocs88eb89bb.(*cgoAllocMap).Free()
		x.ref88eb89bb = nil
	}
}

// NewSQL_INTERVAL_STRUCTRef creates a new wrapper struct with underlying reference set to the original C object.
// Returns nil if the provided pointer to C object is nil too.
func NewSQL_INTERVAL_STRUCTRef(ref unsafe.Pointer) *SQL_INTERVAL_STRUCT {
	if ref == nil {
		return nil
	}
	obj := new(SQL_INTERVAL_STRUCT)
	obj.ref88eb89bb = (*C.SQL_INTERVAL_STRUCT)(unsafe.Pointer(ref))
	return obj
}

// PassRef returns the underlying C object, otherwise it will allocate one and set its values
// from this wrapping struct, counting allocations into an allocation map.
func (x *SQL_INTERVAL_STRUCT) PassRef() (*C.SQL_INTERVAL_STRUCT, *cgoAllocMap) {
	if x == nil {
		return nil, nil
	} else if x.ref88eb89bb != nil {
		return x.ref88eb89bb, nil
	}
	mem88eb89bb := allocSQL_INTERVAL_STRUCTMemory(1)
	ref88eb89bb := (*C.SQL_INTERVAL_STRUCT)(mem88eb89bb)
	allocs88eb89bb := new(cgoAllocMap)
	allocs88eb89bb.Add(mem88eb89bb)

	var cinterval_type_allocs *cgoAllocMap
	ref88eb89bb.interval_type, cinterval_type_allocs = (C.SQLINTERVAL)(x.Interval_type), cgoAllocsUnknown
	allocs88eb89bb.Borrow(cinterval_type_allocs)

	var cinterval_sign_allocs *cgoAllocMap
	ref88eb89bb.interval_sign, cinterval_sign_allocs = (C.SQLSMALLINT)(x.Interval_sign), cgoAllocsUnknown
	allocs88eb89bb.Borrow(cinterval_sign_allocs)

	x.ref88eb89bb = ref88eb89bb
	x.allocs88eb89bb = allocs88eb89bb
	return ref88eb89bb, allocs88eb89bb

}

// PassValue does the same as PassRef except that it will try to dereference the returned pointer.
func (x SQL_INTERVAL_STRUCT) PassValue() (C.SQL_INTERVAL_STRUCT, *cgoAllocMap) {
	if x.ref88eb89bb != nil {
		return *x.ref88eb89bb, nil
	}
	ref, allocs := x.PassRef()
	return *ref, allocs
}

// Deref uses the underlying reference to C object and fills the wrapping struct with values.
// Do not forget to call this method whether you get a struct for C object and want to read its values.
func (x *SQL_INTERVAL_STRUCT) Deref() {
	if x.ref88eb89bb == nil {
		return
	}
	x.Interval_type = (SQLINTERVAL)(x.ref88eb89bb.interval_type)
	x.Interval_sign = (SQLSMALLINT)(x.ref88eb89bb.interval_sign)
}

// allocSQL_NUMERIC_STRUCTMemory allocates memory for type C.SQL_NUMERIC_STRUCT in C.
// The caller is responsible for freeing the this memory via C.free.
func allocSQL_NUMERIC_STRUCTMemory(n int) unsafe.Pointer {
	mem, err := C.calloc(C.size_t(n), (C.size_t)(sizeOfSQL_NUMERIC_STRUCTValue))
	if mem == nil {
		panic(fmt.Sprintln("memory alloc error: ", err))
	}
	return mem
}

const sizeOfSQL_NUMERIC_STRUCTValue = unsafe.Sizeof([1]C.SQL_NUMERIC_STRUCT{})

// Ref returns the underlying reference to C object or nil if struct is nil.
func (x *SQL_NUMERIC_STRUCT) Ref() *C.SQL_NUMERIC_STRUCT {
	if x == nil {
		return nil
	}
	return x.ref56ec20d2
}

// Free invokes alloc map's free mechanism that cleanups any allocated memory using C free.
// Does nothing if struct is nil or has no allocation map.
func (x *SQL_NUMERIC_STRUCT) Free() {
	if x != nil && x.allocs56ec20d2 != nil {
		x.allocs56ec20d2.(*cgoAllocMap).Free()
		x.ref56ec20d2 = nil
	}
}

// NewSQL_NUMERIC_STRUCTRef creates a new wrapper struct with underlying reference set to the original C object.
// Returns nil if the provided pointer to C object is nil too.
func NewSQL_NUMERIC_STRUCTRef(ref unsafe.Pointer) *SQL_NUMERIC_STRUCT {
	if ref == nil {
		return nil
	}
	obj := new(SQL_NUMERIC_STRUCT)
	obj.ref56ec20d2 = (*C.SQL_NUMERIC_STRUCT)(unsafe.Pointer(ref))
	return obj
}

// PassRef returns the underlying C object, otherwise it will allocate one and set its values
// from this wrapping struct, counting allocations into an allocation map.
func (x *SQL_NUMERIC_STRUCT) PassRef() (*C.SQL_NUMERIC_STRUCT, *cgoAllocMap) {
	if x == nil {
		return nil, nil
	} else if x.ref56ec20d2 != nil {
		return x.ref56ec20d2, nil
	}
	mem56ec20d2 := allocSQL_NUMERIC_STRUCTMemory(1)
	ref56ec20d2 := (*C.SQL_NUMERIC_STRUCT)(mem56ec20d2)
	allocs56ec20d2 := new(cgoAllocMap)
	allocs56ec20d2.Add(mem56ec20d2)

	var cprecision_allocs *cgoAllocMap
	ref56ec20d2.precision, cprecision_allocs = (C.SQLCHAR)(x.Precision), cgoAllocsUnknown
	allocs56ec20d2.Borrow(cprecision_allocs)

	var cscale_allocs *cgoAllocMap
	ref56ec20d2.scale, cscale_allocs = (C.SQLSCHAR)(x.Scale), cgoAllocsUnknown
	allocs56ec20d2.Borrow(cscale_allocs)

	var csign_allocs *cgoAllocMap
	ref56ec20d2.sign, csign_allocs = (C.SQLCHAR)(x.Sign), cgoAllocsUnknown
	allocs56ec20d2.Borrow(csign_allocs)

	var cval_allocs *cgoAllocMap
	ref56ec20d2.val, cval_allocs = *(*[16]C.SQLCHAR)(unsafe.Pointer(&x.Val)), cgoAllocsUnknown
	allocs56ec20d2.Borrow(cval_allocs)

	x.ref56ec20d2 = ref56ec20d2
	x.allocs56ec20d2 = allocs56ec20d2
	return ref56ec20d2, allocs56ec20d2

}

// PassValue does the same as PassRef except that it will try to dereference the returned pointer.
func (x SQL_NUMERIC_STRUCT) PassValue() (C.SQL_NUMERIC_STRUCT, *cgoAllocMap) {
	if x.ref56ec20d2 != nil {
		return *x.ref56ec20d2, nil
	}
	ref, allocs := x.PassRef()
	return *ref, allocs
}

// Deref uses the underlying reference to C object and fills the wrapping struct with values.
// Do not forget to call this method whether you get a struct for C object and want to read its values.
func (x *SQL_NUMERIC_STRUCT) Deref() {
	if x.ref56ec20d2 == nil {
		return
	}
	x.Precision = (SQLCHAR)(x.ref56ec20d2.precision)
	x.Scale = (SQLSCHAR)(x.ref56ec20d2.scale)
	x.Sign = (SQLCHAR)(x.ref56ec20d2.sign)
	x.Val = *(*[16]SQLCHAR)(unsafe.Pointer(&x.ref56ec20d2.val))
}

// allocSQLGUIDMemory allocates memory for type C.SQLGUID in C.
// The caller is responsible for freeing the this memory via C.free.
func allocSQLGUIDMemory(n int) unsafe.Pointer {
	mem, err := C.calloc(C.size_t(n), (C.size_t)(sizeOfSQLGUIDValue))
	if mem == nil {
		panic(fmt.Sprintln("memory alloc error: ", err))
	}
	return mem
}

const sizeOfSQLGUIDValue = unsafe.Sizeof([1]C.SQLGUID{})

// Ref returns the underlying reference to C object or nil if struct is nil.
func (x *SQLGUID) Ref() *C.SQLGUID {
	if x == nil {
		return nil
	}
	return x.refe93ebf54
}

// Free invokes alloc map's free mechanism that cleanups any allocated memory using C free.
// Does nothing if struct is nil or has no allocation map.
func (x *SQLGUID) Free() {
	if x != nil && x.allocse93ebf54 != nil {
		x.allocse93ebf54.(*cgoAllocMap).Free()
		x.refe93ebf54 = nil
	}
}

// NewSQLGUIDRef creates a new wrapper struct with underlying reference set to the original C object.
// Returns nil if the provided pointer to C object is nil too.
func NewSQLGUIDRef(ref unsafe.Pointer) *SQLGUID {
	if ref == nil {
		return nil
	}
	obj := new(SQLGUID)
	obj.refe93ebf54 = (*C.SQLGUID)(unsafe.Pointer(ref))
	return obj
}

// PassRef returns the underlying C object, otherwise it will allocate one and set its values
// from this wrapping struct, counting allocations into an allocation map.
func (x *SQLGUID) PassRef() (*C.SQLGUID, *cgoAllocMap) {
	if x == nil {
		return nil, nil
	} else if x.refe93ebf54 != nil {
		return x.refe93ebf54, nil
	}
	meme93ebf54 := allocSQLGUIDMemory(1)
	refe93ebf54 := (*C.SQLGUID)(meme93ebf54)
	allocse93ebf54 := new(cgoAllocMap)
	allocse93ebf54.Add(meme93ebf54)

	var cData1_allocs *cgoAllocMap
	refe93ebf54.Data1, cData1_allocs = (C.DWORD)(x.Data1), cgoAllocsUnknown
	allocse93ebf54.Borrow(cData1_allocs)

	var cData2_allocs *cgoAllocMap
	refe93ebf54.Data2, cData2_allocs = (C.WORD)(x.Data2), cgoAllocsUnknown
	allocse93ebf54.Borrow(cData2_allocs)

	var cData3_allocs *cgoAllocMap
	refe93ebf54.Data3, cData3_allocs = (C.WORD)(x.Data3), cgoAllocsUnknown
	allocse93ebf54.Borrow(cData3_allocs)

	var cData4_allocs *cgoAllocMap
	refe93ebf54.Data4, cData4_allocs = *(*[8]C.BYTE)(unsafe.Pointer(&x.Data4)), cgoAllocsUnknown
	allocse93ebf54.Borrow(cData4_allocs)

	x.refe93ebf54 = refe93ebf54
	x.allocse93ebf54 = allocse93ebf54
	return refe93ebf54, allocse93ebf54

}

// PassValue does the same as PassRef except that it will try to dereference the returned pointer.
func (x SQLGUID) PassValue() (C.SQLGUID, *cgoAllocMap) {
	if x.refe93ebf54 != nil {
		return *x.refe93ebf54, nil
	}
	ref, allocs := x.PassRef()
	return *ref, allocs
}

// Deref uses the underlying reference to C object and fills the wrapping struct with values.
// Do not forget to call this method whether you get a struct for C object and want to read its values.
func (x *SQLGUID) Deref() {
	if x.refe93ebf54 == nil {
		return
	}
	x.Data1 = (uint32)(x.refe93ebf54.Data1)
	x.Data2 = (uint16)(x.refe93ebf54.Data2)
	x.Data3 = (uint16)(x.refe93ebf54.Data3)
	x.Data4 = *(*[8]byte)(unsafe.Pointer(&x.refe93ebf54.Data4))
}
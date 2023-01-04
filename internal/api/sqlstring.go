package api

import "C"
import (
	"unicode/utf16"
)

type WCHAR []uint16

func (w WCHAR) Address() *SQLWCHAR {
	return (*SQLWCHAR)(&w[0])
}

func (w WCHAR) String() string {
	return nts(string(utf16.Decode(w)))
}

type CHAR []uint8

func (w CHAR) Address() *SQLCHAR {
	return (*SQLCHAR)(&w[0])
}

func (w CHAR) String() string {
	return nts(string(w))
}

func nts(str string) string {
	for i := range str {
		if int(str[i]) == 0 {
			return str[:i]
		}
	}
	return str
}

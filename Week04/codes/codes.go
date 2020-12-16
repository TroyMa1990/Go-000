package hsbcodes

import (
	"google.golang.org/grpc/codes"
	st "google.golang.org/grpc/status"
)

type Code uint32

const (
	OK       Code = 0
	Unknown  Code = 1
	Continue Code = 2
	Nothing  Code = 3

	MyServerLoadStageFailed Code = 100
)

// Return Code Message
func (c Code) String() string {
	m := make(map[Code]string)
	m[OK] = "OK"
	m[Unknown] = "Unknown"
	m[Continue] = "Continue"
	m[Nothing] = "Nothing"

	m[MyServerLoadStageFailed] = "MyServerLoadStageFailed"

	return m[c]
}

// Return Grpc Code
func (c Code) GrpcCode() codes.Code {
	m := make(map[Code]codes.Code)
	m[OK] = codes.OK
	m[Unknown] = codes.Unknown
	m[Continue] = codes.OK
	m[Nothing] = codes.OK

	m[MyServerLoadStageFailed] = codes.DataLoss

	if cs, ok := m[c]; ok {
		return cs
	}

	return codes.Unknown
}

// Convert String To Code
func StringToCode(s string) Code {
	m := make(map[string]Code)
	m["OK"] = OK
	m["Unknown"] = Unknown
	m["Continue"] = Continue
	m["Nothing"] = Nothing

	m["MyServerLoadStageFailed"] = MyServerLoadStageFailed

	if c, ok := m[s]; ok {
		return c
	}

	return Unknown
}

// Return GRPC Error
func (c Code) GrpcError() error {
	return st.Error(c.GrpcCode(), c.String())
}

// Error Interface Want
func (c Code) Error() string {
	return c.String()
}

// Convert error To Code
func FromError(err error) Code {
	if err == nil {
		return OK
	}

	if s, ok := st.FromError(err); ok {
		return StringToCode(s.Message())
	}
	return Unknown
}

const (
	FromUnknow          int = 0
	FromCommon          int = 1
	FromMyServerControl int = 2
)

// Return Error Come From
// 1. Common; 2. MyServer;
func ErrorComeFrom(err error) int {
	if s, ok := st.FromError(err); ok {
		myCode := StringToCode(s.Message())
		if int(myCode) >= 100 {
			return FromMyServerControl
		} else if int(myCode) >= 0 {
			return FromCommon
		}
	}
	return FromUnknow
}

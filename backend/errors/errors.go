package errors

import (
	"fmt"
	"net/http"
)

type Scope int

const (
	ScopeSystem   Scope = 1
	ScopeBusiness Scope = 2
)

type OPError struct {
	HttpCode int
	Message  string
	Code     int
	Scope    Scope
}

func (O OPError) Error() string {
	var scope string
	switch O.Scope {
	case ScopeSystem:
		scope = "system"
	case ScopeBusiness:
		scope = "business"
	}
	return fmt.Sprintf("%s error: [%d]%s.", scope, O.Code, O.Message)
}

func NewSystemOPError(code int, message string, httpCodes ...int) *OPError {
	httpCode := http.StatusOK
	if len(httpCodes) > 0 {
		httpCode = httpCodes[0]
	}
	return NewOpError(code, message, ScopeSystem, httpCode)
}
func NewOpError(code int, message string, scope Scope, httpCode int) *OPError {
	return &OPError{
		Message:  message,
		Code:     code,
		Scope:    scope,
		HttpCode: httpCode,
	}
}
func NewBusinessError(code int, message string, httpCodes ...int) *OPError {
	httpCode := http.StatusOK
	if len(httpCodes) > 0 {
		httpCode = httpCodes[0]
	}
	return NewOpError(code, message, ScopeBusiness, httpCode)
}

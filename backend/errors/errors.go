package errors

import "fmt"

type Scope int

const (
	ScopeSystem   Scope = 1
	ScopeBusiness Scope = 2
)

type OPError struct {
	Message string
	Code    int
	Scope   Scope
}

func (O OPError) Error() string {
	var scope string
	switch O.Scope {
	case ScopeSystem:
		scope = "system"
		break
	case ScopeBusiness:
		scope = "business"
	}
	return fmt.Sprintf("%s : %d -> %s.", scope, O.Code, O.Message)
}

func NewSystemOPError(code int, message string) *OPError {
	return &OPError{
		Message: message,
		Code:    code,
		Scope:   ScopeSystem,
	}
}
func NewBusinessError(code int, message string) *OPError {
	return &OPError{
		Message: message,
		Code:    code,
		Scope:   ScopeBusiness,
	}
}

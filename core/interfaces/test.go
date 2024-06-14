package interfaces

import "testing"

type Test interface {
	Setup(*testing.T)
	Cleanup()
}

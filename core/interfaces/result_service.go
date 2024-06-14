package interfaces

import (
	"github.com/crawlab-team/crawlab-db/generic"
	"time"
)

type ResultService interface {
	Insert(records ...interface{}) (err error)
	List(query generic.ListQuery, opts *generic.ListOptions) (results []interface{}, err error)
	Count(query generic.ListQuery) (n int, err error)
	Index(fields []string)
	SetTime(t time.Time)
	GetTime() (t time.Time)
}

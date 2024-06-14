package interfaces

import "time"

type TaskStat interface {
	Model
	GetCreateTs() (ts time.Time)
	SetCreateTs(ts time.Time)
	GetStartTs() (ts time.Time)
	SetStartTs(ts time.Time)
	GetEndTs() (ts time.Time)
	SetEndTs(ts time.Time)
	GetWaitDuration() (d int64)
	SetWaitDuration(d int64)
	GetRuntimeDuration() (d int64)
	SetRuntimeDuration(d int64)
	GetTotalDuration() (d int64)
	SetTotalDuration(d int64)
	GetResultCount() (c int64)
	SetResultCount(c int64)
	GetErrorLogCount() (c int64)
	SetErrorLogCount(c int64)
}

package interfaces

import "time"

type Export interface {
	GetId() string
	GetType() string
	GetTarget() string
	GetFilter() Filter
	GetStatus() string
	GetStartTs() time.Time
	GetEndTs() time.Time
	GetDownloadPath() string
}

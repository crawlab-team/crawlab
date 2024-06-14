package entity

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"time"
)

type Export struct {
	Id           string            `json:"id"`
	Type         string            `json:"type"`
	Target       string            `json:"target"`
	Filter       interfaces.Filter `json:"filter"`
	Status       string            `json:"status"`
	StartTs      time.Time         `json:"start_ts"`
	EndTs        time.Time         `json:"end_ts"`
	FileName     string            `json:"file_name"`
	DownloadPath string            `json:"-"`
	Limit        int               `json:"-"`
}

func (e *Export) GetId() string {
	return e.Id
}

func (e *Export) GetType() string {
	return e.Type
}

func (e *Export) GetTarget() string {
	return e.Target
}

func (e *Export) GetFilter() interfaces.Filter {
	return e.Filter
}

func (e *Export) GetStatus() string {
	return e.Status
}

func (e *Export) GetStartTs() time.Time {
	return e.StartTs
}

func (e *Export) GetEndTs() time.Time {
	return e.EndTs
}

func (e *Export) GetDownloadPath() string {
	return e.DownloadPath
}

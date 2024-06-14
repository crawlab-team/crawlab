package interfaces

import (
	"github.com/robfig/cron/v3"
	"time"
)

type ScheduleService interface {
	WithConfigPath
	Module
	GetLocation() (loc *time.Location)
	SetLocation(loc *time.Location)
	GetDelay() (delay bool)
	SetDelay(delay bool)
	GetSkip() (skip bool)
	SetSkip(skip bool)
	GetUpdateInterval() (interval time.Duration)
	SetUpdateInterval(interval time.Duration)
	Enable(s Schedule, args ...interface{}) (err error)
	Disable(s Schedule, args ...interface{}) (err error)
	Update()
	GetCron() (c *cron.Cron)
}

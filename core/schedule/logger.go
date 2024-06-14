package schedule

import (
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/go-trace"
	"github.com/robfig/cron/v3"
	"strings"
)

type Logger struct {
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	p := l.getPlaceholder(len(keysAndValues))
	log.Infof(fmt.Sprintf("cron: %s %s", msg, p), keysAndValues...)
}

func (l *Logger) Error(err error, msg string, keysAndValues ...interface{}) {
	p := l.getPlaceholder(len(keysAndValues))
	log.Errorf(fmt.Sprintf("cron: %s %s", msg, p), keysAndValues...)
	trace.PrintError(err)
}

func (l *Logger) getPlaceholder(n int) (s string) {
	var arr []string
	for i := 0; i < n; i++ {
		arr = append(arr, "%v")
	}
	return strings.Join(arr, " ")
}

func NewLogger() cron.Logger {
	return &Logger{}
}

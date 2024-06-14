package utils

import (
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/go-trace"
	"time"
)

func BackoffErrorNotify(prefix string) backoff.Notify {
	return func(err error, duration time.Duration) {
		log.Errorf("%s error: %v. reattempt in %.1f seconds...", prefix, err, duration.Seconds())
		trace.PrintError(err)
	}
}

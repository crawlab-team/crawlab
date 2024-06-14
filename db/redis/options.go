package redis

import (
	"github.com/crawlab-team/crawlab/db"
	"time"
)

type Option func(c db.RedisClient)

func WithBackoffMaxInterval(interval time.Duration) Option {
	return func(c db.RedisClient) {
		c.SetBackoffMaxInterval(interval)
	}
}

func WithTimeout(timeout int) Option {
	return func(c db.RedisClient) {
		c.SetTimeout(timeout)
	}
}

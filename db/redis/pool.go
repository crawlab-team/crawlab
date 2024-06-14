package redis

import (
	"github.com/crawlab-team/go-trace"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"time"
)

func NewRedisPool() *redis.Pool {
	var address = viper.GetString("redis.address")
	var port = viper.GetString("redis.port")
	var database = viper.GetString("redis.database")
	var password = viper.GetString("redis.password")

	// normalize params
	if address == "" {
		address = "localhost"
	}
	if port == "" {
		port = "6379"
	}
	if database == "" {
		database = "1"
	}

	var url string
	if password == "" {
		url = "redis://" + address + ":" + port + "/" + database
	} else {
		url = "redis://x:" + password + "@" + address + ":" + port + "/" + database
	}
	return &redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			return redis.DialURL(url,
				redis.DialConnectTimeout(time.Second*10),
				redis.DialReadTimeout(time.Second*600),
				redis.DialWriteTimeout(time.Second*10),
			)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return trace.TraceError(err)
		},
		MaxIdle:         10,
		MaxActive:       0,
		IdleTimeout:     300 * time.Second,
		Wait:            false,
		MaxConnLifetime: 0,
	}
}

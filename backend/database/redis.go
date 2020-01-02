package database

import (
	"context"
	"crawlab/entity"
	"crawlab/utils"
	"errors"
	"github.com/apex/log"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"runtime/debug"
	"strings"
	"time"
)

var RedisClient *Redis

type Redis struct {
	pool *redis.Pool
}

type Mutex struct {
	Name   string
	expiry time.Duration
	tries  int
	delay  time.Duration
	value  string
}

func NewRedisClient() *Redis {
	return &Redis{pool: NewRedisPool()}
}

func (r *Redis) RPush(collection string, value interface{}) error {
	c := r.pool.Get()
	defer utils.Close(c)

	if _, err := c.Do("RPUSH", collection, value); err != nil {
		debug.PrintStack()
		return err
	}
	return nil
}

func (r *Redis) LPush(collection string, value interface{}) error {
	c := r.pool.Get()
	defer utils.Close(c)

	if _, err := c.Do("RPUSH", collection, value); err != nil {
		debug.PrintStack()
		return err
	}
	return nil
}

func (r *Redis) LPop(collection string) (string, error) {
	c := r.pool.Get()
	defer utils.Close(c)

	value, err2 := redis.String(c.Do("LPOP", collection))
	if err2 != nil {
		return value, err2
	}
	return value, nil
}

func (r *Redis) HSet(collection string, key string, value string) error {
	c := r.pool.Get()
	defer utils.Close(c)

	if _, err := c.Do("HSET", collection, key, value); err != nil {
		debug.PrintStack()
		return err
	}
	return nil
}

func (r *Redis) HGet(collection string, key string) (string, error) {
	c := r.pool.Get()
	defer utils.Close(c)

	value, err2 := redis.String(c.Do("HGET", collection, key))
	if err2 != nil {
		return value, err2
	}
	return value, nil
}

func (r *Redis) HDel(collection string, key string) error {
	c := r.pool.Get()
	defer utils.Close(c)

	if _, err := c.Do("HDEL", collection, key); err != nil {
		return err
	}
	return nil
}

func (r *Redis) HKeys(collection string) ([]string, error) {
	c := r.pool.Get()
	defer utils.Close(c)

	value, err2 := redis.Strings(c.Do("HKeys", collection))
	if err2 != nil {
		return []string{}, err2
	}
	return value, nil
}

func (r *Redis) BRPop(collection string, timeout int) (string, error) {
	if timeout <= 0 {
		timeout = 60
	}
	c := r.pool.Get()
	defer utils.Close(c)

	values, err := redis.Strings(c.Do("BRPOP", collection, timeout))
	if err != nil {
		return "", err
	}
	return values[1], nil
}

func NewRedisPool() *redis.Pool {
	var address = viper.GetString("redis.address")
	var port = viper.GetString("redis.port")
	var database = viper.GetString("redis.database")
	var password = viper.GetString("redis.password")

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
			return err
		},
		MaxIdle:         10,
		MaxActive:       0,
		IdleTimeout:     300 * time.Second,
		Wait:            false,
		MaxConnLifetime: 0,
	}
}

func InitRedis() error {
	RedisClient = NewRedisClient()
	return nil
}

func Pub(channel string, msg entity.NodeMessage) error {
	if _, err := RedisClient.Publish(channel, utils.GetJson(msg)); err != nil {
		log.Errorf("publish redis error: %s", err.Error())
		debug.PrintStack()
		return err
	}
	return nil
}

func Sub(channel string, consume ConsumeFunc) error {
	ctx := context.Background()
	if err := RedisClient.Subscribe(ctx, consume, channel); err != nil {
		log.Errorf("subscribe redis error: %s", err.Error())
		debug.PrintStack()
		return err
	}
	return nil
}

// 构建同步锁key
func (r *Redis) getLockKey(lockKey string) string {
	lockKey = strings.ReplaceAll(lockKey, ":", "-")
	return "nodes:lock:" + lockKey
}

// 获得锁
func (r *Redis) Lock(lockKey string) (int64, error) {
	c := r.pool.Get()
	defer utils.Close(c)
	lockKey = r.getLockKey(lockKey)

	ts := time.Now().Unix()
	ok, err := c.Do("SET", lockKey, ts, "NX", "PX", 30000)
	if err != nil {
		log.Errorf("get lock fail with error: %s", err.Error())
		debug.PrintStack()
		return 0, err
	}
	if err == nil && ok == nil {
		log.Errorf("the lockKey is locked: key=%s", lockKey)
		return 0, errors.New("the lockKey is locked")
	}
	return ts, nil
}

func (r *Redis) UnLock(lockKey string, value int64) {
	c := r.pool.Get()
	defer utils.Close(c)
	lockKey = r.getLockKey(lockKey)

	getValue, err := redis.Int64(c.Do("GET", lockKey))
	if err != nil {
		log.Errorf("get lockKey error: %s", err.Error())
		debug.PrintStack()
		return
	}

	if getValue != value {
		log.Errorf("the lockKey value diff: %d, %d", value, getValue)
		return
	}

	v, err := redis.Int64(c.Do("DEL", lockKey))
	if err != nil {
		log.Errorf("unlock failed, error: %s", err.Error())
		debug.PrintStack()
		return
	}

	if v == 0 {
		log.Errorf("unlock failed: key=%s", lockKey)
		return
	}
}

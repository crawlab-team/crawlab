package redis

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/db"
	"github.com/crawlab-team/crawlab/db/errors"
	"github.com/crawlab-team/crawlab/db/utils"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/gomodule/redigo/redis"
	"reflect"
	"strings"
	"time"
)

type Client struct {
	// settings
	backoffMaxInterval time.Duration
	timeout            int

	// internals
	pool *redis.Pool
}

func (client *Client) Ping() error {
	c := client.pool.Get()
	defer utils.Close(c)
	if _, err := redis.String(c.Do("PING")); err != nil {
		if err != redis.ErrNil {
			return trace.TraceError(err)
		}
		return err
	}
	return nil
}

func (client *Client) Keys(pattern string) (values []string, err error) {
	c := client.pool.Get()
	defer utils.Close(c)

	values, err = redis.Strings(c.Do("KEYS", pattern))
	if err != nil {
		return nil, trace.TraceError(err)
	}
	return values, nil
}

func (client *Client) AllKeys() (values []string, err error) {
	return client.Keys("*")
}

func (client *Client) Get(collection string) (value string, err error) {
	c := client.pool.Get()
	defer utils.Close(c)

	value, err = redis.String(c.Do("GET", collection))
	if err != nil {
		return "", trace.TraceError(err)
	}
	return value, nil
}

func (client *Client) Set(collection string, value string) (err error) {
	c := client.pool.Get()
	defer utils.Close(c)

	value, err = redis.String(c.Do("SET", collection, value))
	if err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (client *Client) Del(collection string) error {
	c := client.pool.Get()
	defer utils.Close(c)

	if _, err := c.Do("DEL", collection); err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (client *Client) RPush(collection string, value interface{}) error {
	c := client.pool.Get()
	defer utils.Close(c)

	if _, err := c.Do("RPUSH", collection, value); err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (client *Client) LPush(collection string, value interface{}) error {
	c := client.pool.Get()
	defer utils.Close(c)

	if _, err := c.Do("LPUSH", collection, value); err != nil {
		if err != redis.ErrNil {
			return trace.TraceError(err)
		}
		return err
	}
	return nil
}

func (client *Client) LPop(collection string) (string, error) {
	c := client.pool.Get()
	defer utils.Close(c)

	value, err := redis.String(c.Do("LPOP", collection))
	if err != nil {
		if err != redis.ErrNil {
			return value, trace.TraceError(err)
		}
		return value, err
	}
	return value, nil
}

func (client *Client) RPop(collection string) (string, error) {
	c := client.pool.Get()
	defer utils.Close(c)

	value, err := redis.String(c.Do("RPOP", collection))
	if err != nil {
		if err != redis.ErrNil {
			return value, trace.TraceError(err)
		}
		return value, err
	}
	return value, nil
}

func (client *Client) LLen(collection string) (int, error) {
	c := client.pool.Get()
	defer utils.Close(c)

	value, err := redis.Int(c.Do("LLEN", collection))
	if err != nil {
		return 0, trace.TraceError(err)
	}
	return value, nil
}

func (client *Client) BRPop(collection string, timeout int) (value string, err error) {
	if timeout <= 0 {
		timeout = 60
	}
	c := client.pool.Get()
	defer utils.Close(c)

	values, err := redis.Strings(c.Do("BRPOP", collection, timeout))
	if err != nil {
		if err != redis.ErrNil {
			return value, trace.TraceError(err)
		}
		return value, err
	}
	return values[1], nil
}

func (client *Client) BLPop(collection string, timeout int) (value string, err error) {
	if timeout <= 0 {
		timeout = 60
	}
	c := client.pool.Get()
	defer utils.Close(c)

	values, err := redis.Strings(c.Do("BLPOP", collection, timeout))
	if err != nil {
		if err != redis.ErrNil {
			return value, trace.TraceError(err)
		}
		return value, err
	}
	return values[1], nil
}

func (client *Client) HSet(collection string, key string, value string) error {
	c := client.pool.Get()
	defer utils.Close(c)

	if _, err := c.Do("HSET", collection, key, value); err != nil {
		if err != redis.ErrNil {
			return trace.TraceError(err)
		}
		return err
	}
	return nil
}

func (client *Client) HGet(collection string, key string) (string, error) {
	c := client.pool.Get()
	defer utils.Close(c)
	value, err := redis.String(c.Do("HGET", collection, key))
	if err != nil && err != redis.ErrNil {
		if err != redis.ErrNil {
			return value, trace.TraceError(err)
		}
		return value, err
	}
	return value, nil
}

func (client *Client) HDel(collection string, key string) error {
	c := client.pool.Get()
	defer utils.Close(c)

	if _, err := c.Do("HDEL", collection, key); err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (client *Client) HScan(collection string) (results map[string]string, err error) {
	c := client.pool.Get()
	defer utils.Close(c)

	var (
		cursor int64
		items  []string
	)

	results = map[string]string{}

	for {
		values, err := redis.Values(c.Do("HSCAN", collection, cursor))
		if err != nil {
			if err != redis.ErrNil {
				return nil, trace.TraceError(err)
			}
			return nil, err
		}

		values, err = redis.Scan(values, &cursor, &items)
		if err != nil {
			if err != redis.ErrNil {
				return nil, trace.TraceError(err)
			}
			return nil, err
		}
		for i := 0; i < len(items); i += 2 {
			key := items[i]
			value := items[i+1]
			results[key] = value
		}
		if cursor == 0 {
			break
		}
	}
	return results, nil
}

func (client *Client) HKeys(collection string) (results []string, err error) {
	c := client.pool.Get()
	defer utils.Close(c)

	results, err = redis.Strings(c.Do("HKEYS", collection))
	if err != nil {
		if err != redis.ErrNil {
			return results, trace.TraceError(err)
		}
		return results, err
	}
	return results, nil
}

func (client *Client) ZAdd(collection string, score float32, value interface{}) (err error) {
	c := client.pool.Get()
	defer utils.Close(c)

	if _, err := c.Do("ZADD", collection, score, value); err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (client *Client) ZCount(collection string, min string, max string) (count int, err error) {
	c := client.pool.Get()
	defer utils.Close(c)

	count, err = redis.Int(c.Do("ZCOUNT", collection, min, max))
	if err != nil {
		return 0, trace.TraceError(err)
	}
	return count, nil
}

func (client *Client) ZCountAll(collection string) (count int, err error) {
	return client.ZCount(collection, "-inf", "+inf")
}

func (client *Client) ZScan(collection string, pattern string, count int) (values []string, err error) {
	c := client.pool.Get()
	defer utils.Close(c)

	values, err = redis.Strings(c.Do("ZSCAN", collection, 0, pattern, count))
	if err != nil {
		if err != redis.ErrNil {
			return nil, trace.TraceError(err)
		}
		return nil, err
	}
	return values, nil
}

func (client *Client) ZPopMax(collection string, count int) (results []string, err error) {
	c := client.pool.Get()
	defer utils.Close(c)

	results = []string{}

	values, err := redis.Strings(c.Do("ZPOPMAX", collection, count))
	if err != nil {
		if err != redis.ErrNil {
			return nil, trace.TraceError(err)
		}
		return nil, err
	}

	for i := 0; i < len(values); i += 2 {
		v := values[i]
		results = append(results, v)
	}

	return results, nil
}

func (client *Client) ZPopMin(collection string, count int) (results []string, err error) {
	c := client.pool.Get()
	defer utils.Close(c)

	results = []string{}

	values, err := redis.Strings(c.Do("ZPOPMIN", collection, count))
	if err != nil {
		if err != redis.ErrNil {
			return nil, trace.TraceError(err)
		}
		return nil, err
	}

	for i := 0; i < len(values); i += 2 {
		v := values[i]
		results = append(results, v)
	}

	return results, nil
}

func (client *Client) ZPopMaxOne(collection string) (value string, err error) {
	c := client.pool.Get()
	defer utils.Close(c)

	values, err := client.ZPopMax(collection, 1)
	if err != nil {
		return "", err
	}
	if values == nil || len(values) == 0 {
		return "", nil
	}
	return values[0], nil
}

func (client *Client) ZPopMinOne(collection string) (value string, err error) {
	c := client.pool.Get()
	defer utils.Close(c)

	values, err := client.ZPopMin(collection, 1)
	if err != nil {
		return "", err
	}
	if values == nil || len(values) == 0 {
		return "", nil
	}
	return values[0], nil
}

func (client *Client) BZPopMax(collection string, timeout int) (value string, err error) {
	c := client.pool.Get()
	defer utils.Close(c)

	values, err := redis.Strings(c.Do("BZPOPMAX", collection, timeout))
	if err != nil {
		if err != redis.ErrNil {
			return "", trace.TraceError(err)
		}
		return "", err
	}
	if len(values) < 3 {
		return "", trace.TraceError(errors.ErrorRedisInvalidType)
	}
	return values[1], nil
}

func (client *Client) BZPopMin(collection string, timeout int) (value string, err error) {
	c := client.pool.Get()
	defer utils.Close(c)

	values, err := redis.Strings(c.Do("BZPOPMIN", collection, timeout))
	if err != nil {
		if err != redis.ErrNil {
			return "", trace.TraceError(err)
		}
		return "", err
	}
	if len(values) < 3 {
		return "", trace.TraceError(errors.ErrorRedisInvalidType)
	}
	return values[1], nil
}

func (client *Client) Lock(lockKey string) (value int64, err error) {
	c := client.pool.Get()
	defer utils.Close(c)
	lockKey = client.getLockKey(lockKey)

	ts := time.Now().Unix()
	ok, err := c.Do("SET", lockKey, ts, "NX", "PX", 30000)
	if err != nil {
		if err != redis.ErrNil {
			return value, trace.TraceError(err)
		}
		return value, err
	}
	if ok == nil {
		return 0, trace.TraceError(errors.ErrorRedisLocked)
	}
	return ts, nil
}

func (client *Client) UnLock(lockKey string, value int64) {
	c := client.pool.Get()
	defer utils.Close(c)
	lockKey = client.getLockKey(lockKey)

	getValue, err := redis.Int64(c.Do("GET", lockKey))
	if err != nil {
		log.Errorf("get lockKey error: %s", err.Error())
		return
	}

	if getValue != value {
		log.Errorf("the lockKey value diff: %d, %d", value, getValue)
		return
	}

	v, err := redis.Int64(c.Do("DEL", lockKey))
	if err != nil {
		log.Errorf("unlock failed, error: %s", err.Error())
		return
	}

	if v == 0 {
		log.Errorf("unlock failed: key=%s", lockKey)
		return
	}
}

func (client *Client) MemoryStats() (stats map[string]int64, err error) {
	stats = map[string]int64{}
	c := client.pool.Get()
	defer utils.Close(c)
	values, err := redis.Values(c.Do("MEMORY", "STATS"))
	for i, v := range values {
		t := reflect.TypeOf(v)
		if t.Kind() == reflect.Slice {
			vc, _ := redis.String(v, err)
			if utils.ContainsString(MemoryStatsMetrics, vc) {
				stats[vc], _ = redis.Int64(values[i+1], err)
			}
		}
	}
	if err != nil {
		if err != redis.ErrNil {
			return stats, trace.TraceError(err)
		}
		return stats, err
	}
	return stats, nil
}

func (client *Client) SetBackoffMaxInterval(interval time.Duration) {
	client.backoffMaxInterval = interval
}

func (client *Client) SetTimeout(timeout int) {
	client.timeout = timeout
}

func (client *Client) init() (err error) {
	b := backoff.NewExponentialBackOff()
	b.MaxInterval = client.backoffMaxInterval
	if err := backoff.Retry(func() error {
		err := client.Ping()
		if err != nil {
			log.WithError(err).Warnf("waiting for redis pool active connection. will after %f seconds try again.", b.NextBackOff().Seconds())
		}
		return nil
	}, b); err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (client *Client) getLockKey(lockKey string) string {
	lockKey = strings.ReplaceAll(lockKey, ":", "-")
	return "nodes:lock:" + lockKey
}

func (client *Client) getTimeout(timeout int) (res int) {
	if timeout == 0 {
		return client.timeout
	}
	return timeout
}

var client db.RedisClient

func NewRedisClient(opts ...Option) (client *Client, err error) {
	// client
	client = &Client{
		backoffMaxInterval: 20 * time.Second,
		pool:               NewRedisPool(),
	}

	// apply options
	for _, opt := range opts {
		opt(client)
	}

	// init
	if err := client.init(); err != nil {
		return nil, err
	}

	return client, nil
}

func GetRedisClient() (c db.RedisClient, err error) {
	if client != nil {
		return client, nil
	}
	c, err = NewRedisClient()
	if err != nil {
		return nil, err
	}

	return c, nil
}

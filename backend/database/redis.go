package database

import (
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"runtime/debug"
)

var RedisClient = Redis{}

type ConsumeFunc func(channel string, message []byte) error

type Redis struct {
}

func (r *Redis) RPush(collection string, value interface{}) error {
	c, err := GetRedisConn()
	if err != nil {
		debug.PrintStack()
		return err
	}
	defer c.Close()

	if _, err := c.Do("RPUSH", collection, value); err != nil {
		debug.PrintStack()
		return err
	}
	return nil
}

func (r *Redis) LPop(collection string) (string, error) {
	c, err := GetRedisConn()
	if err != nil {
		debug.PrintStack()
		return "", err
	}
	defer c.Close()

	value, err2 := redis.String(c.Do("LPOP", collection))
	if err2 != nil {
		return value, err2
	}
	return value, nil
}

func (r *Redis) HSet(collection string, key string, value string) error {
	c, err := GetRedisConn()
	if err != nil {
		debug.PrintStack()
		return err
	}
	defer c.Close()

	if _, err := c.Do("HSET", collection, key, value); err != nil {
		debug.PrintStack()
		return err
	}
	return nil
}

func (r *Redis) HGet(collection string, key string) (string, error) {
	c, err := GetRedisConn()
	if err != nil {
		debug.PrintStack()
		return "", err
	}
	defer c.Close()

	value, err2 := redis.String(c.Do("HGET", collection, key))
	if err2 != nil {
		return value, err2
	}
	return value, nil
}

func (r *Redis) HDel(collection string, key string) error {
	c, err := GetRedisConn()
	if err != nil {
		debug.PrintStack()
		return err
	}
	defer c.Close()

	if _, err := c.Do("HDEL", collection, key); err != nil {
		return err
	}
	return nil
}

func (r *Redis) HKeys(collection string) ([]string, error) {
	c, err := GetRedisConn()
	if err != nil {
		debug.PrintStack()
		return []string{}, err
	}
	defer c.Close()

	value, err2 := redis.Strings(c.Do("HKeys", collection))
	if err2 != nil {
		return []string{}, err2
	}
	return value, nil
}

func GetRedisConn() (redis.Conn, error) {
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
	c, err := redis.DialURL(url)
	if err != nil {
		debug.PrintStack()
		return c, err
	}
	return c, nil
}

func InitRedis() error {
	return nil
}

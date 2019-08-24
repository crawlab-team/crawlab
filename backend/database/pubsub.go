package database

import (
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/gomodule/redigo/redis"
	"time"
	"unsafe"
)

type SubscribeCallback func(channel, message string)

type Subscriber struct {
	client redis.PubSubConn
	cbMap  map[string]SubscribeCallback
}

func (c *Subscriber) Connect() {
	conn, err := GetRedisConn()
	if err != nil {
		log.Fatalf("redis dial failed.")
	}

	c.client = redis.PubSubConn{Conn: conn}
	c.cbMap = make(map[string]SubscribeCallback)

	//retry connect redis 5 times, or panic
	index := 0
	go func(i int) {
		for {
			log.Debug("wait...")
			switch res := c.client.Receive().(type) {
			case redis.Message:
				i = 0
				channel := (*string)(unsafe.Pointer(&res.Channel))
				message := (*string)(unsafe.Pointer(&res.Data))
				c.cbMap[*channel](*channel, *message)
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", res.Channel, res.Kind, res.Count)
			case error:
				log.Error("error handle redis connection...")

				time.Sleep(2 * time.Second)
				if i > 5 {
					panic(errors.New("redis connection failed too many times, panic"))
				}
				con, err := GetRedisConn()
				if err != nil {
					log.Error("redis dial failed")
					continue
				}
				c.client = redis.PubSubConn{Conn: con}
				i += 1

				continue
			}
		}
	}(index)

}

func (c *Subscriber) Close() {
	err := c.client.Close()
	if err != nil {
		log.Errorf("redis close error.")
	}
}

func (c *Subscriber) Subscribe(channel interface{}, cb SubscribeCallback) {
	err := c.client.Subscribe(channel)
	if err != nil {
		log.Fatalf("redis Subscribe error.")
	}

	c.cbMap[channel.(string)] = cb
}

func Publish(channel string, msg string) error {
	c, err := GetRedisConn()
	if err != nil {
		return err
	}

	if _, err := c.Do("PUBLISH", channel, msg); err != nil {
		return err
	}

	return nil
}

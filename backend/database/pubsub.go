package database

import (
	"fmt"
	"github.com/apex/log"
	"github.com/gomodule/redigo/redis"
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

	go func() {
		for {
			log.Debug("wait...")
			switch res := c.client.Receive().(type) {
			case redis.Message:
				channel := (*string)(unsafe.Pointer(&res.Channel))
				message := (*string)(unsafe.Pointer(&res.Data))
				c.cbMap[*channel](*channel, *message)
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", res.Channel, res.Kind, res.Count)
			case error:
				log.Error("error handle...")
				continue
			}
		}
	}()

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

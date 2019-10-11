package database

import (
	"context"
	"fmt"
	"github.com/apex/log"
	"github.com/gomodule/redigo/redis"
	errors2 "github.com/pkg/errors"
	"time"
)

type ConsumeFunc func(message redis.Message) error

func (r *Redis) Close() {
	err := r.pool.Close()
	if err != nil {
		log.Errorf("redis close error.")
	}
}
func (r *Redis) subscribe(ctx context.Context, consume ConsumeFunc, channel ...string) error {
	psc := redis.PubSubConn{Conn: r.pool.Get()}
	if err := psc.Subscribe(redis.Args{}.AddFlat(channel)...); err != nil {
		return err
	}
	done := make(chan error, 1)
	tick := time.NewTicker(time.Second * 3)
	defer tick.Stop()
	go func() {
		defer func() { _ = psc.Close() }()
		for {
			switch msg := psc.Receive().(type) {
			case error:
				done <- fmt.Errorf("redis pubsub receive err: %v", msg)
				return
			case redis.Message:
				if err := consume(msg); err != nil {
					fmt.Printf("redis pubsub consume message err: %v", err)
					continue
				}
			case redis.Subscription:
				fmt.Println(msg)

				if msg.Count == 0 {
					// all channels are unsubscribed
					return
				}
			}

		}
	}()
	// start a new goroutine to receive message
	for {
		select {
		case <-ctx.Done():
			if err := psc.Unsubscribe(); err != nil {
				fmt.Printf("redis pubsub unsubscribe err: %v \n", err)
			}
			done <- nil
		case <-tick.C:
			//fmt.Printf("ping message  \n")
			if err := psc.Ping(""); err != nil {
				done <- err
			}
		case err := <-done:
			close(done)
			return err
		}
	}

}
func (r *Redis) Subscribe(ctx context.Context, consume ConsumeFunc, channel ...string) error {
	index := 0
	go func() {
		for {
			err := r.subscribe(ctx, consume, channel...)
			fmt.Println(err)

			if err == nil {
				break
			}
			time.Sleep(5 * time.Second)
			index += 1
			fmt.Printf("try reconnect %d times \n", index)
		}
	}()
	return nil
}
func (r *Redis) Publish(channel, message string) (n int, err error) {
	conn := r.pool.Get()
	defer func() { _ = conn.Close() }()
	n, err = redis.Int(conn.Do("PUBLISH", channel, message))
	if err != nil {
		return 0, errors2.Wrapf(err, "redis publish %s %s", channel, message)
	}
	return
}

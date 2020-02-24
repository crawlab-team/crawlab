package utils

import (
	"sync"
)

var TaskExecChanMap = NewChanMap()

type ChanMap struct {
	m sync.Map
}

func NewChanMap() *ChanMap {
	return &ChanMap{m: sync.Map{}}
}

func (cm *ChanMap) Chan(key string) chan string {
	if ch, ok := cm.m.Load(key); ok {
		return ch.(interface{}).(chan string)
	}
	ch := make(chan string, 10)
	cm.m.Store(key, ch)
	return ch
}

func (cm *ChanMap) ChanBlocked(key string) chan string {
	if ch, ok := cm.m.Load(key); ok {
		return ch.(interface{}).(chan string)
	}
	ch := make(chan string)
	cm.m.Store(key, ch)
	return ch
}

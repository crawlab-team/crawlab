package utils

var TaskExecChanMap = NewChanMap()

type ChanMap struct {
	m map[string]chan string
}

func NewChanMap() *ChanMap {
	return &ChanMap{m: make(map[string]chan string)}
}

func (cm *ChanMap) Chan(key string) chan string {
	if ch, ok := cm.m[key]; ok {
		return ch
	}
	ch := make(chan string, 10)
	cm.m[key] = ch
	return ch
}

func (cm *ChanMap) ChanBlocked(key string) chan string {
	if ch, ok := cm.m[key]; ok {
		return ch
	}
	ch := make(chan string)
	cm.m[key] = ch
	return ch
}

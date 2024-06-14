package interfaces

type EventFn func(data ...interface{}) (err error)

type EventService interface {
	Register(key, include, exclude string, ch *chan EventData)
	Unregister(key string)
	SendEvent(eventName string, data ...interface{})
}

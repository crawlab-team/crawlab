package event

func SendEvent(eventName string, data ...interface{}) {
	svc := NewEventService()
	svc.SendEvent(eventName, data...)
}

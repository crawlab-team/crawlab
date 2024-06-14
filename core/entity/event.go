package entity

type EventData struct {
	Event string
	Data  interface{}
}

func (d *EventData) GetEvent() string {
	return d.Event
}

func (d *EventData) GetData() interface{} {
	return d.Data
}

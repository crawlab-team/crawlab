package interfaces

type EventData interface {
	GetEvent() string
	GetData() interface{}
}

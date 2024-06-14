package parser

type Entity interface {
	GetType() string
	SetType(string)
	GetName() string
	SetName(string)
}

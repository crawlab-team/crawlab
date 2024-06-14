package interfaces

type Address interface {
	Entity
	String() string
	IsEmpty() bool
}

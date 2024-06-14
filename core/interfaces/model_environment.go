package interfaces

type Environment interface {
	Model
	GetKey() (key string)
	SetKey(key string)
	GetValue() (value string)
	SetValue(value string)
}

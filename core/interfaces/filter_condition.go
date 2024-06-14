package interfaces

type FilterCondition interface {
	GetKey() (key string)
	SetKey(key string)
	GetOp() (op string)
	SetOp(op string)
	GetValue() (value interface{})
	SetValue(value interface{})
}

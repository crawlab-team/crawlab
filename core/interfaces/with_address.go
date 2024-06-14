package interfaces

type WithAddress interface {
	GetAddress() (address Address)
	SetAddress(address Address)
}

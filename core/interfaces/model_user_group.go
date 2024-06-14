package interfaces

type UserGroup interface {
	Model
	GetUsers() (users []User, err error)
}

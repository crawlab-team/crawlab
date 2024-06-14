package interfaces

type User interface {
	Model
	GetUsername() (name string)
	GetPassword() (p string)
	GetRole() (r string)
	GetEmail() (email string)
}

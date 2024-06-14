package interfaces

type GrpcBase interface {
	WithConfigPath
	Init() (err error)
	Start() (err error)
	Stop() (err error)
	Register() (err error)
}

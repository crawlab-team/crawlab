package interfaces

type WithConfigPath interface {
	GetConfigPath() (path string)
	SetConfigPath(path string)
}

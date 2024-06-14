package interfaces

type ControllerParams interface {
	IsZero() (ok bool)
	IsDefault() (ok bool)
}

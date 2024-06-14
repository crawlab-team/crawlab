package interfaces

type ModelDelegateMethod string

type ModelDelegate interface {
	Add() error
	Save() error
	Delete() error
	GetArtifact() (ModelArtifact, error)
	GetModel() Model
	Refresh() error
	ToBytes(interface{}) ([]byte, error)
}

const (
	ModelDelegateMethodAdd         = "add"
	ModelDelegateMethodSave        = "save"
	ModelDelegateMethodDelete      = "delete"
	ModelDelegateMethodGetArtifact = "get-artifact"
	ModelDelegateMethodRefresh     = "refresh"
	ModelDelegateMethodChange      = "change"
)

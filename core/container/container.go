package container

import (
	"go.uber.org/dig"
)

var c = dig.New()

func GetContainer() *dig.Container {
	return c
}

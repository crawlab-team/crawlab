package apps

import (
	"github.com/apex/log"
)

type Handler struct {
}

func (app *Handler) Init() {
}

func (app *Handler) Start() {
}

func (app *Handler) Wait() {
	DefaultWait()
}

func (app *Handler) Stop() {
	log.Info("handler has stopped")
}

func NewHandler() *Handler {
	return &Handler{}
}

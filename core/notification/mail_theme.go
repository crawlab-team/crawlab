package notification

import "github.com/matcornic/hermes/v2"

type MailTheme interface {
	hermes.Theme
	GetStyle() string
}

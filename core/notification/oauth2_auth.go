package notification

import (
	"fmt"
	"net/smtp"
)

type XOAuth2Auth struct {
	Username, Token string
}

func (a *XOAuth2Auth) Start(_ *smtp.ServerInfo) (string, []byte, error) {
	return "XOAUTH2", []byte("user=" + a.Username + "\x01auth=Bearer " + a.Token + "\x01\x01"), nil
}

func (a *XOAuth2Auth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		return nil, fmt.Errorf("unexpected server challenge: %s", fromServer)
	}
	return nil, nil
}

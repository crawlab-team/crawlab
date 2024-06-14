package vcs

import (
	"os/user"
	"path/filepath"
)

func getDefaultPublicKeyPath() (path string) {
	u, err := user.Current()
	if err != nil {
		return path
	}
	path = filepath.Join(u.HomeDir, ".ssh", "id_rsa")
	return
}

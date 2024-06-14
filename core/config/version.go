package config

import "strings"

const Version = "v0.6.3"

func GetVersion() (v string) {
	if strings.HasPrefix(Version, "v") {
		return Version
	}
	return "v" + Version
}

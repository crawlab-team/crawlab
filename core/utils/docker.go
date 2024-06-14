package utils

func IsDocker() (ok bool) {
	return EnvIsTrue("docker", false)
}

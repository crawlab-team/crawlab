package utils

func IsMaster() bool {
	return EnvIsTrue("node.master", false)
}

func GetNodeType() string {
	if IsMaster() {
		return "master"
	} else {
		return "worker"
	}
}

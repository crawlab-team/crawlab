package entity

type SystemInfo struct {
	Edition string `json:"edition"` // edition. e.g. community / pro
	Version string `json:"version"` // version. e.g. v0.6.0
}

package entity

type SystemInfo struct {
	ARCH        string       `json:"arch"`
	OS          string       `json:"os"`
	Hostname    string       `json:"host_name"`
	NumCpu      int          `json:"num_cpu"`
	Executables []Executable `json:"executables"`
}

type Executable struct {
	Path        string `json:"path"`
	FileName    string `json:"file_name"`
	DisplayName string `json:"display_name"`
}

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

type Lang struct {
	Name              string   `json:"name"`
	ExecutableName    string   `json:"executable_name"`
	ExecutablePaths   []string `json:"executable_paths"`
	DepExecutablePath string   `json:"dep_executable_path"`
	LockPath          string   `json:"lock_path"`
	InstallScript     string   `json:"install_script"`
	InstallStatus     string   `json:"install_status"`
	DepFileName       string   `json:"dep_file_name"`
	InstallDepArgs    string   `json:"install_dep_cmd"`
}

type Dependency struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Installed   bool   `json:"installed"`
}

type PackageJson struct {
	Dependencies map[string]string `json:"dependencies"`
}

package utils

import "crawlab/entity"

func GetLangList() []entity.Lang {
	list := []entity.Lang{
		{
			Name:              "Python",
			ExecutableName:    "python",
			ExecutablePaths:   []string{"/usr/bin/python", "/usr/local/bin/python"},
			DepExecutablePath: "/usr/local/bin/pip",
			LockPath:          "/tmp/install-python.lock",
		},
		{
			Name:              "Node.js",
			ExecutableName:    "node",
			ExecutablePaths:   []string{"/usr/bin/node", "/usr/local/bin/node"},
			DepExecutablePath: "/usr/local/bin/npm",
			LockPath:          "/tmp/install-nodejs.lock",
			InstallScript:     "install-nodejs.sh",
		},
		{
			Name:            "Java",
			ExecutableName:  "java",
			ExecutablePaths: []string{"/usr/bin/java", "/usr/local/bin/java"},
			LockPath:        "/tmp/install-java.lock",
			InstallScript:   "install-java.sh",
		},
		{
			Name:            ".Net Core",
			ExecutableName:  "dotnet",
			ExecutablePaths: []string{"/usr/bin/dotnet", "/usr/local/bin/dotnet"},
			LockPath:        "/tmp/install-dotnet.lock",
			InstallScript:   "install-dotnet.sh",
		},
		{
			Name:            "PHP",
			ExecutableName:  "php",
			ExecutablePaths: []string{"/usr/bin/php", "/usr/local/bin/php"},
			LockPath:        "/tmp/install-php.lock",
			InstallScript:   "install-php.sh",
		},
	}
	return list
}

// 获取语言列表
func GetLangListPlain() []entity.Lang {
	list := GetLangList()
	return list
}

// 根据语言名获取语言实例，不包含状态
func GetLangFromLangNamePlain(name string) entity.Lang {
	langList := GetLangListPlain()
	for _, lang := range langList {
		if lang.ExecutableName == name {
			return lang
		}
	}
	return entity.Lang{}
}

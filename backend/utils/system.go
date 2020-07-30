package utils

import (
	"crawlab/entity"
	"encoding/json"
	"github.com/apex/log"
	"io/ioutil"
	"runtime/debug"
)

func GetLangList() []entity.Lang {
	list := []entity.Lang{
		{
			Name:              "Python",
			ExecutableName:    "python",
			ExecutablePaths:   []string{"/usr/bin/python", "/usr/local/bin/python"},
			DepExecutablePath: "/usr/local/bin/pip",
			LockPath:          "/tmp/install-python.lock",
			DepFileName:       "requirements.txt",
			InstallDepArgs:    "install -i https://pypi.tuna.tsinghua.edu.cn/simple -r requirements.txt",
		},
		{
			Name:              "Node.js",
			ExecutableName:    "node",
			ExecutablePaths:   []string{"/usr/bin/node", "/usr/local/bin/node"},
			DepExecutablePath: "/usr/local/bin/npm",
			LockPath:          "/tmp/install-nodejs.lock",
			InstallScript:     "install-nodejs.sh",
			DepFileName:       "package.json",
			InstallDepArgs:    "install -g --registry=https://registry.npm.taobao.org",
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
		{
			Name:            "Golang",
			ExecutableName:  "go",
			ExecutablePaths: []string{"/usr/bin/go", "/usr/local/bin/go"},
			LockPath:        "/tmp/install-go.lock",
			InstallScript:   "install-go.sh",
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

func GetPackageJsonDeps(filepath string) (deps []string, err error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Errorf("get package.json deps error: " + err.Error())
		debug.PrintStack()
		return deps, err
	}
	var packageJson entity.PackageJson
	if err := json.Unmarshal(data, &packageJson); err != nil {
		log.Errorf("get package.json deps error: " + err.Error())
		debug.PrintStack()
		return deps, err
	}

	for d, v := range packageJson.Dependencies {
		deps = append(deps, d+"@"+v)
	}

	return deps, nil
}

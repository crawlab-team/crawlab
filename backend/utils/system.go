package utils

import (
	"crawlab/constants"
	"crawlab/entity"
	"encoding/json"
	"github.com/apex/log"
	"github.com/spf13/viper"
	"io/ioutil"
	"path"
	"runtime/debug"
	"strings"
)

func GetLangList() []entity.Lang {
	list := []entity.Lang{
		// 语言
		{
			Name:              "Python",
			ExecutableName:    "python",
			ExecutablePaths:   []string{"/usr/bin/python", "/usr/local/bin/python"},
			DepExecutablePath: "/usr/local/bin/pip",
			LockPath:          "/tmp/install-python.lock",
			DepFileName:       "requirements.txt",
			InstallDepArgs:    "install -i https://pypi.tuna.tsinghua.edu.cn/simple -r requirements.txt",
			Type:              constants.LangTypeLang,
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
			Type:              constants.LangTypeLang,
		},
		{
			Name:            "Java",
			ExecutableName:  "java",
			ExecutablePaths: []string{"/usr/bin/java", "/usr/local/bin/java"},
			LockPath:        "/tmp/install-java.lock",
			InstallScript:   "install-java.sh",
			Type:            constants.LangTypeLang,
		},
		{
			Name:            ".Net Core",
			ExecutableName:  "dotnet",
			ExecutablePaths: []string{"/usr/bin/dotnet", "/usr/local/bin/dotnet"},
			LockPath:        "/tmp/install-dotnet.lock",
			InstallScript:   "install-dotnet.sh",
			Type:            constants.LangTypeLang,
		},
		{
			Name:            "PHP",
			ExecutableName:  "php",
			ExecutablePaths: []string{"/usr/bin/php", "/usr/local/bin/php"},
			LockPath:        "/tmp/install-php.lock",
			InstallScript:   "install-php.sh",
			Type:            constants.LangTypeLang,
		},
		{
			Name:            "Golang",
			ExecutableName:  "go",
			ExecutablePaths: []string{"/usr/bin/go", "/usr/local/bin/go"},
			LockPath:        "/tmp/install-go.lock",
			InstallScript:   "install-go.sh",
			Type:            constants.LangTypeLang,
		},
		// WebDriver
		{
			Name:            "Chrome Driver",
			ExecutableName:  "chromedriver",
			ExecutablePaths: []string{"/usr/bin/chromedriver", "/usr/local/bin/chromedriver"},
			LockPath:        "/tmp/install-chromedriver.lock",
			InstallScript:   "install-chromedriver.sh",
			Type:            constants.LangTypeWebDriver,
		},
		{
			Name:            "Firefox",
			ExecutableName:  "firefox",
			ExecutablePaths: []string{"/usr/bin/firefox", "/usr/local/bin/firefox"},
			LockPath:        "/tmp/install-firefox.lock",
			InstallScript:   "install-firefox.sh",
			Type:            constants.LangTypeWebDriver,
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

// 获取系统脚本列表
func GetSystemScripts() (res []string) {
	scriptsPath := viper.GetString("server.scripts")
	for _, fInfo := range ListDir(scriptsPath) {
		if !fInfo.IsDir() && strings.HasSuffix(fInfo.Name(), ".sh") {
			res = append(res, fInfo.Name())
		}
	}
	return res
}

func GetSystemScriptPath(scriptName string) string {
	scriptsPath := viper.GetString("server.scripts")
	for _, name := range GetSystemScripts() {
		if name == scriptName {
			return path.Join(scriptsPath, name)
		}
	}
	return ""
}

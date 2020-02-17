package services

import (
	"bytes"
	"crawlab/model"
	"encoding/json"
	"github.com/apex/log"
	"os/exec"
	"runtime/debug"
	"strings"
)

func GetScrapySpiderNames(s model.Spider) ([]string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("scrapy", "list")
	cmd.Dir = s.Src
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return []string{}, err
	}

	spiderNames := strings.Split(stdout.String(), "\n")

	var res []string
	for _, sn := range spiderNames {
		if sn != "" {
			res = append(res, sn)
		}
	}

	return res, nil
}

func GetScrapySettings(s model.Spider) (res []map[string]interface{}, err error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("crawlab", "settings")
	cmd.Dir = s.Src
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Errorf(err.Error())
		log.Errorf(stderr.String())
		debug.PrintStack()
		return res, err
	}

	log.Infof(stdout.String())
	if err := json.Unmarshal([]byte(stdout.String()), &res); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return res, err
	}

	return res, nil
}

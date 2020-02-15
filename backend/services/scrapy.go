package services

import (
	"bytes"
	"crawlab/model"
	"os/exec"
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

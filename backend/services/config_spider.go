package services

import (
	"crawlab/constants"
	"crawlab/entity"
	"crawlab/model"
	"crawlab/model/config_spider"
	"errors"
	"fmt"
	"strings"
)

func GenerateConfigSpiderFiles(spider model.Spider, configData entity.ConfigSpiderData) error {
	// 校验Spiderfile正确性
	if err := ValidateSpiderfile(configData); err != nil {
		return err
	}

	// 构造代码生成器
	generator := config_spider.ScrapyGenerator{
		Spider:     spider,
		ConfigData: configData,
	}

	// 生成代码
	if err := generator.Generate(); err != nil {
		return err
	}

	return nil
}

// 验证Spiderfile
func ValidateSpiderfile(configData entity.ConfigSpiderData) error {
	// 获取所有字段
	fields := config_spider.GetAllFields(configData)

	// 校验是否存在 start_url
	if configData.StartUrl == "" {
		return errors.New("spiderfile start_url is empty")
	}

	// 校验是否存在 stages
	if len(configData.Stages) == 0 {
		return errors.New("spiderfile stages is empty")
	}

	// 校验stages
	dict := map[string]int{}
	for stageName, stage := range configData.Stages {
		// stage 名称不能为空
		if stageName == "" {
			return errors.New("spiderfile stage name is empty")
		}

		// stage 名称不能为保留字符串
		// NOTE: 如果有其他Engine，可以扩展，默认为Scrapy
		if configData.Engine == "" || configData.Engine == constants.EngineScrapy {
			if strings.Contains(constants.ScrapyProtectedStageNames, stageName) {
				return errors.New(fmt.Sprintf("spiderfile stage name '%s' is protected", stageName))
			}
		} else if configData.Engine == constants.EngineColly {
			return errors.New(fmt.Sprintf("engine '%s' is not implemented", stageName))
		}

		// stage 名称不能重复
		if dict[stageName] == 1 {
			return errors.New("spiderfile stage name should be unique")
		}
		dict[stageName] = 1

		// stage 字段不能为空
		if len(stage.Fields) == 0 {
			return errors.New(fmt.Sprintf("spiderfile stage '%s' has no fields", stageName))
		}

		// stage 的下一个 stage 只能有一个
		hasNextStage := false
		for _, field := range stage.Fields {
			if field.NextStage != "" {
				if hasNextStage {
					return errors.New("spiderfile stage fields should have only 1 next_stage")
				}
				hasNextStage = true
			}
		}

		// 如果 stage 的 is_list 为 true 但 list_css 为空，报错
		if stage.IsList && stage.ListCss == "" {
			return errors.New("spiderfile stage with is_list = true should have list_css being set")
		}
	}

	// 校验字段唯一性
	if !IsUniqueConfigSpiderFields(fields) {
		return errors.New("spiderfile fields not unique")
	}

	// 字段名称不能为保留字符串
	for _, field := range fields {
		if strings.Contains(constants.ScrapyProtectedFieldNames, field.Name) {
			return errors.New(fmt.Sprintf("spiderfile field name '%s' is protected", field.Name))
		}
	}

	return nil
}

func IsUniqueConfigSpiderFields(fields []entity.Field) bool {
	dict := map[string]int{}
	for _, field := range fields {
		if dict[field.Name] == 1 {
			return false
		}
		dict[field.Name] = 1
	}
	return true
}

package config_spider

import "github.com/crawlab-team/crawlab/core/entity"

func GetAllFields(data entity.ConfigSpiderData) []entity.Field {
	var fields []entity.Field
	for _, stage := range data.Stages {
		fields = append(fields, stage.Fields...)
	}
	return fields
}

func GetStartStageName(data entity.ConfigSpiderData) string {
	// 如果 start_stage 设置了且在 stages 里，则返回
	if data.StartStage != "" {
		return data.StartStage
	}

	// 否则返回第一个 stage
	for _, stage := range data.Stages {
		return stage.Name
	}
	return ""
}

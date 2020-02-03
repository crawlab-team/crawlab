package config_spider

import "crawlab/entity"

func GetAllFields(data entity.ConfigSpiderData) []entity.Field {
	var fields []entity.Field
	for _, stage := range data.Stages {
		for _, field := range stage.Fields {
			fields = append(fields, field)
		}
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

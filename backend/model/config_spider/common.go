package config_spider

import "crawlab/entity"

func GetAllFields(data entity.ConfigSpiderData) []entity.Field {
	var fields []entity.Field
	for _, stage := range data.Stages {
		if stage.IsList {
			for _, field := range stage.Fields {
				fields = append(fields, field)
			}
		}
	}
	return fields
}

func GetStartStageName(data entity.ConfigSpiderData) string {
	for stageName := range data.Stages {
		return stageName
	}
	return ""
}

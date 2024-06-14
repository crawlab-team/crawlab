package entity

import "github.com/crawlab-team/crawlab/core/interfaces"

type ModelDelegate struct {
	Id       interfaces.ModelId       `json:"id"`
	ColName  string                   `json:"col_name"`
	Doc      interfaces.Model         `json:"doc"`
	Artifact interfaces.ModelArtifact `json:"a"`
	User     interfaces.User          `json:"u"`
}

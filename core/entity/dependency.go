package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type DependencyResult struct {
	Name          string               `json:"name,omitempty" bson:"name,omitempty"`
	NodeIds       []primitive.ObjectID `json:"node_ids,omitempty" bson:"node_ids,omitempty"`
	Versions      []string             `json:"versions,omitempty" bson:"versions,omitempty"`
	LatestVersion string               `json:"latest_version" bson:"latest_version"`
	Count         int                  `json:"count,omitempty" bson:"count,omitempty"`
	Upgradable    bool                 `json:"upgradable" bson:"upgradable"`
	Downgradable  bool                 `json:"downgradable" bson:"downgradable"`
	Installable   bool                 `json:"installable" bson:"installable"`
}

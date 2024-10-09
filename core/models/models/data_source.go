package models

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataSource struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Type        string             `json:"type" bson:"type"`
	Description string             `json:"description" bson:"description"`
	Host        string             `json:"host" bson:"host"`
	Port        int                `json:"port" bson:"port"`
	Url         string             `json:"url" bson:"url"`
	Hosts       []string           `json:"hosts" bson:"hosts"`
	Database    string             `json:"database" bson:"database"`
	Username    string             `json:"username" bson:"username"`
	Password    string             `json:"password,omitempty" bson:"-"`
	ConnectType string             `json:"connect_type" bson:"connect_type"`
	Status      string             `json:"status" bson:"status"`
	Error       string             `json:"error" bson:"error"`
	Extra       map[string]string  `json:"extra,omitempty" bson:"extra,omitempty"`
}

func (ds *DataSource) GetId() (id primitive.ObjectID) {
	return ds.Id
}

func (ds *DataSource) SetId(id primitive.ObjectID) {
	ds.Id = id
}

type DataSourceList []DataSource

func (l *DataSourceList) GetModels() (res []interfaces.Model) {
	for i := range *l {
		d := (*l)[i]
		res = append(res, &d)
	}
	return res
}

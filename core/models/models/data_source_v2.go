package models

type DataSourceV2 struct {
	any                     `collection:"data_sources"`
	BaseModelV2[DataSource] `bson:",inline"`
	Name                    string            `json:"name" bson:"name"`
	Type                    string            `json:"type" bson:"type"`
	Description             string            `json:"description" bson:"description"`
	Host                    string            `json:"host" bson:"host"`
	Port                    string            `json:"port" bson:"port"`
	Url                     string            `json:"url" bson:"url"`
	Hosts                   []string          `json:"hosts" bson:"hosts"`
	Database                string            `json:"database" bson:"database"`
	Username                string            `json:"username" bson:"username"`
	Password                string            `json:"password,omitempty" bson:"-"`
	ConnectType             string            `json:"connect_type" bson:"connect_type"`
	Status                  string            `json:"status" bson:"status"`
	Error                   string            `json:"error" bson:"error"`
	Extra                   map[string]string `json:"extra,omitempty" bson:"extra,omitempty"`
}

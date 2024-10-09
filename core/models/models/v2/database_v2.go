package models

import (
	"time"
)

type DatabaseV2 struct {
	any                     `collection:"databases"`
	BaseModelV2[DatabaseV2] `bson:",inline"`
	Name                    string    `json:"name" bson:"name"`
	Description             string    `json:"description" bson:"description"`
	DataSource              string    `json:"data_source" bson:"data_source"`
	Host                    string    `json:"host" bson:"host"`
	Port                    int       `json:"port" bson:"port"`
	URI                     string    `json:"uri,omitempty" bson:"uri,omitempty"`
	Database                string    `json:"database,omitempty" bson:"database,omitempty"`
	Username                string    `json:"username,omitempty" bson:"username,omitempty"`
	Password                string    `json:"password,omitempty" bson:"-"`
	EncryptedPassword       string    `json:"-,omitempty" bson:"encrypted_password,omitempty"`
	Status                  string    `json:"status" bson:"status"`
	Error                   string    `json:"error" bson:"error"`
	Active                  bool      `json:"active" bson:"active"`
	ActiveAt                time.Time `json:"active_ts" bson:"active_ts"`
	IsDefault               bool      `json:"is_default" bson:"-"`

	MongoParams *struct {
		AuthSource    string `json:"auth_source,omitempty" bson:"auth_source,omitempty"`
		AuthMechanism string `json:"auth_mechanism,omitempty" bson:"auth_mechanism,omitempty"`
	} `json:"mongo_params,omitempty" bson:"mongo_params,omitempty"`
	PostgresParams *struct {
		SSLMode string `json:"ssl_mode,omitempty" bson:"ssl_mode,omitempty"`
	} `json:"postgres_params,omitempty" bson:"postgres_params,omitempty"`
	SnowflakeParams *struct {
		Account   string `json:"account,omitempty" bson:"account,omitempty"`
		Schema    string `json:"schema,omitempty" bson:"schema,omitempty"`
		Warehouse string `json:"warehouse,omitempty" bson:"warehouse,omitempty"`
		Role      string `json:"role,omitempty" bson:"role,omitempty"`
	} `json:"snowflake_params,omitempty" bson:"snowflake_params,omitempty"`
	CassandraParams *struct {
		Keyspace string `json:"keyspace,omitempty" bson:"keyspace,omitempty"`
	} `json:"cassandra_params,omitempty" bson:"cassandra_params,omitempty"`
	HiveParams *struct {
		Auth string `json:"auth,omitempty" bson:"auth,omitempty"`
	} `json:"hive_params,omitempty" bson:"hive_params,omitempty"`
	RedisParams *struct {
		DB int `json:"db,omitempty" bson:"db,omitempty"`
	} `json:"redis_params,omitempty" bson:"redis_params,omitempty"`
}

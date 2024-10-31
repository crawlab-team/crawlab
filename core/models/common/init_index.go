package common

import (
	models2 "github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitIndexes() {
	// nodes
	RecreateIndexes(mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.NodeV2{})), []mongo2.IndexModel{
		{Keys: bson.M{"key": 1}},       // key
		{Keys: bson.M{"name": 1}},      // name
		{Keys: bson.M{"is_master": 1}}, // is_master
		{Keys: bson.M{"status": 1}},    // status
		{Keys: bson.M{"enabled": 1}},   // enabled
		{Keys: bson.M{"active": 1}},    // active
	})

	// projects
	RecreateIndexes(mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.ProjectV2{})), []mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
	})

	// spiders
	RecreateIndexes(mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.SpiderV2{})), []mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
		{Keys: bson.M{"type": 1}},
		{Keys: bson.M{"col_id": 1}},
		{Keys: bson.M{"project_id": 1}},
	})

	// tasks
	RecreateIndexes(mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.TaskV2{})), []mongo2.IndexModel{
		{Keys: bson.M{"spider_id": 1}},
		{Keys: bson.M{"status": 1}},
		{Keys: bson.M{"node_id": 1}},
		{Keys: bson.M{"schedule_id": 1}},
		{Keys: bson.M{"type": 1}},
		{Keys: bson.M{"mode": 1}},
		{Keys: bson.M{"priority": 1}},
		{Keys: bson.M{"parent_id": 1}},
		{Keys: bson.M{"has_sub": 1}},
		{Keys: bson.M{"created_ts": -1}, Options: (&options.IndexOptions{}).SetExpireAfterSeconds(60 * 60 * 24 * 30)},
		{Keys: bson.M{"node_id": 1, "status": 1}},
	})

	// task stats
	RecreateIndexes(mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.TaskStatV2{})), []mongo2.IndexModel{
		{Keys: bson.M{"created_ts": -1}, Options: (&options.IndexOptions{}).SetExpireAfterSeconds(60 * 60 * 24 * 30)},
	})

	// schedules
	RecreateIndexes(mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.ScheduleV2{})), []mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
		{Keys: bson.M{"spider_id": 1}},
		{Keys: bson.M{"enabled": 1}},
	})

	// users
	RecreateIndexes(mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.UserV2{})), []mongo2.IndexModel{
		{Keys: bson.M{"username": 1}},
		{Keys: bson.M{"role": 1}},
		{Keys: bson.M{"email": 1}},
	})

	// settings
	RecreateIndexes(mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.SettingV2{})), []mongo2.IndexModel{
		{Keys: bson.M{"key": 1}, Options: options.Index().SetUnique(true)},
	})

	// tokens
	RecreateIndexes(mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.TokenV2{})), []mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
	})

	// data collections
	RecreateIndexes(mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.DataCollectionV2{})), []mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
	})
}

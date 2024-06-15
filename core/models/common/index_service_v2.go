package common

import (
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexesV2() {
	// nodes
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models.NodeV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"key": 1}},       // key
		{Keys: bson.M{"name": 1}},      // name
		{Keys: bson.M{"is_master": 1}}, // is_master
		{Keys: bson.M{"status": 1}},    // status
		{Keys: bson.M{"enabled": 1}},   // enabled
		{Keys: bson.M{"active": 1}},    // active
	})

	// projects
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models.ProjectV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
	})

	// spiders
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models.SpiderV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
		{Keys: bson.M{"type": 1}},
		{Keys: bson.M{"col_id": 1}},
		{Keys: bson.M{"project_id": 1}},
	})

	// tasks
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models.TaskV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"spider_id": 1}},
		{Keys: bson.M{"status": 1}},
		{Keys: bson.M{"node_id": 1}},
		{Keys: bson.M{"schedule_id": 1}},
		{Keys: bson.M{"type": 1}},
		{Keys: bson.M{"mode": 1}},
		{Keys: bson.M{"priority": 1}},
		{Keys: bson.M{"parent_id": 1}},
		{Keys: bson.M{"has_sub": 1}},
		{Keys: bson.M{"create_ts": -1}},
	})

	// task stats
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models.TaskStatV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"create_ts": 1}},
	})

	// schedules
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models.ScheduleV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
		{Keys: bson.M{"spider_id": 1}},
		{Keys: bson.M{"enabled": 1}},
	})

	// users
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models.UserV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"username": 1}},
		{Keys: bson.M{"role": 1}},
		{Keys: bson.M{"email": 1}},
	})

	// settings
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models.SettingV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"key": 1}},
	})

	// tokens
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models.TokenV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
	})

	// variables
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models.VariableV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"key": 1}},
	})

	// data sources
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models.DataSourceV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
	})

	// data collections
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models.DataCollectionV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
	})

	// roles
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models.RoleV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.D{{"key", 1}}, Options: options.Index().SetUnique(true)},
	})

	// user role relations
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models.UserRoleV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.D{{"user_id", 1}, {"role_id", 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{"role_id", 1}, {"user_id", 1}}, Options: options.Index().SetUnique(true)},
	})

	// permissions
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models.PermissionV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.D{{"key", 1}}, Options: options.Index().SetUnique(true)},
	})

	// role permission relations
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models.RolePermissionV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.D{{"role_id", 1}, {"permission_id", 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{"permission_id", 1}, {"role_id", 1}}, Options: options.Index().SetUnique(true)},
	})

	// dependencies
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models.DependencyV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{
			Keys: bson.D{
				{"type", 1},
				{"node_id", 1},
				{"name", 1},
			},
			Options: (&options.IndexOptions{}).SetUnique(true),
		},
	})

	// dependency settings
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models.DependencySettingV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{
			Keys: bson.D{
				{"type", 1},
				{"node_id", 1},
				{"name", 1},
			},
			Options: (&options.IndexOptions{}).SetUnique(true),
		},
	})

	// dependency logs
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models.DependencyLogV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{
			Keys: bson.D{{"task_id", 1}},
		},
		{
			Keys:    bson.D{{"update_ts", 1}},
			Options: (&options.IndexOptions{}).SetExpireAfterSeconds(60 * 60 * 24),
		},
	})

	// dependency tasks
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models.DependencyTaskV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{
			Keys: bson.D{
				{"update_ts", 1},
			},
			Options: (&options.IndexOptions{}).SetExpireAfterSeconds(60 * 60 * 24),
		},
	})
}

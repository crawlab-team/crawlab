package common

import (
	models2 "github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexesV2() {
	// nodes
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.NodeV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"key": 1}},       // key
		{Keys: bson.M{"name": 1}},      // name
		{Keys: bson.M{"is_master": 1}}, // is_master
		{Keys: bson.M{"status": 1}},    // status
		{Keys: bson.M{"enabled": 1}},   // enabled
		{Keys: bson.M{"active": 1}},    // active
	})

	// projects
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.ProjectV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
	})

	// spiders
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.SpiderV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
		{Keys: bson.M{"type": 1}},
		{Keys: bson.M{"col_id": 1}},
		{Keys: bson.M{"project_id": 1}},
	})

	// tasks
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.TaskV2{})).MustCreateIndexes([]mongo2.IndexModel{
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
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.TaskStatV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"create_ts": 1}},
	})

	// schedules
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.ScheduleV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
		{Keys: bson.M{"spider_id": 1}},
		{Keys: bson.M{"enabled": 1}},
	})

	// users
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.UserV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"username": 1}},
		{Keys: bson.M{"role": 1}},
		{Keys: bson.M{"email": 1}},
	})

	// settings
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.SettingV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.D{{"key", 1}}, Options: options.Index().SetUnique(true)},
	})

	// tokens
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.TokenV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
	})

	// variables
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.VariableV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"key": 1}},
	})

	// data sources
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.DatabaseV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
	})

	// data collections
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.DataCollectionV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
	})

	// roles
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.RoleV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.D{{"key", 1}}, Options: options.Index().SetUnique(true)},
	})

	// user role relations
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.UserRoleV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.D{{"user_id", 1}, {"role_id", 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{"role_id", 1}, {"user_id", 1}}, Options: options.Index().SetUnique(true)},
	})

	// permissions
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.PermissionV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.D{{"key", 1}}, Options: options.Index().SetUnique(true)},
	})

	// role permission relations
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.RolePermissionV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.D{{"role_id", 1}, {"permission_id", 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{"permission_id", 1}, {"role_id", 1}}, Options: options.Index().SetUnique(true)},
	})

	// dependencies
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.DependencyV2{})).MustCreateIndexes([]mongo2.IndexModel{
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
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.DependencySettingV2{})).MustCreateIndexes([]mongo2.IndexModel{
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
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.DependencyLogV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{
			Keys: bson.D{{"task_id", 1}},
		},
		{
			Keys:    bson.D{{"update_ts", 1}},
			Options: (&options.IndexOptions{}).SetExpireAfterSeconds(60 * 60 * 24),
		},
	})

	// dependency tasks
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.DependencyTaskV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{
			Keys: bson.D{
				{"update_ts", 1},
			},
			Options: (&options.IndexOptions{}).SetExpireAfterSeconds(60 * 60 * 24),
		},
	})

	// metrics
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.MetricV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{
			Keys: bson.D{
				{"created_ts", -1},
			},
			Options: (&options.IndexOptions{}).SetExpireAfterSeconds(60 * 60 * 24 * 30),
		},
		{
			Keys: bson.D{
				{"node_id", 1},
			},
		},
		{
			Keys: bson.D{
				{"type", 1},
			},
		},
	})

	// notification requests
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.NotificationRequestV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{
			Keys: bson.D{
				{"created_ts", -1},
			},
			Options: (&options.IndexOptions{}).SetExpireAfterSeconds(60 * 60 * 24 * 7),
		},
		{
			Keys: bson.D{
				{"channel_id", 1},
			},
		},
		{
			Keys: bson.D{
				{"setting_id", 1},
			},
		},
		{
			Keys: bson.D{
				{"status", 1},
			},
		},
	})

	// databases
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.DatabaseV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{
			Keys: bson.D{
				{"data_source_id", 1},
			},
		},
	})

	// database metrics
	mongo.GetMongoCol(service.GetCollectionNameByInstance(models2.DatabaseMetricV2{})).MustCreateIndexes([]mongo2.IndexModel{
		{
			Keys: bson.D{
				{"created_ts", -1},
			},
			Options: (&options.IndexOptions{}).SetExpireAfterSeconds(60 * 60 * 24 * 30),
		},
		{
			Keys: bson.D{
				{"database_id", 1},
			},
		},
		{
			Keys: bson.D{
				{"type", 1},
			},
		},
	})
}

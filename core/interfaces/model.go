package interfaces

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model interface {
	GetId() (id primitive.ObjectID)
	SetId(id primitive.ObjectID)
}

type ModelV2 interface {
	GetId() (id primitive.ObjectID)
	SetId(id primitive.ObjectID)
	SetCreated(by primitive.ObjectID)
	SetUpdated(by primitive.ObjectID)
}

type ModelId int

const (
	ModelIdArtifact = iota
	ModelIdTag
	ModelIdNode
	ModelIdProject
	ModelIdSpider
	ModelIdTask
	ModelIdJob
	ModelIdSchedule
	ModelIdUser
	ModelIdSetting
	ModelIdToken
	ModelIdVariable
	ModelIdTaskQueue
	ModelIdTaskStat
	ModelIdSpiderStat
	ModelIdDataSource
	ModelIdDataCollection
	ModelIdResult
	ModelIdPassword
	ModelIdExtraValue
	ModelIdGit
	ModelIdRole
	ModelIdUserRole
	ModelIdPermission
	ModelIdRolePermission
	ModelIdEnvironment
	ModelIdDependencySetting
)

const (
	ModelColNameArtifact          = "artifacts"
	ModelColNameTag               = "tags"
	ModelColNameNode              = "nodes"
	ModelColNameProject           = "projects"
	ModelColNameSpider            = "spiders"
	ModelColNameTask              = "tasks"
	ModelColNameJob               = "jobs"
	ModelColNameSchedule          = "schedules"
	ModelColNameUser              = "users"
	ModelColNameSetting           = "settings"
	ModelColNameToken             = "tokens"
	ModelColNameVariable          = "variables"
	ModelColNameTaskQueue         = "task_queue"
	ModelColNameTaskStat          = "task_stats"
	ModelColNameSpiderStat        = "spider_stats"
	ModelColNameDataSource        = "data_sources"
	ModelColNameDataCollection    = "data_collections"
	ModelColNamePasswords         = "passwords"
	ModelColNameExtraValues       = "extra_values"
	ModelColNameGit               = "gits"
	ModelColNameRole              = "roles"
	ModelColNameUserRole          = "user_roles"
	ModelColNamePermission        = "permissions"
	ModelColNameRolePermission    = "role_permissions"
	ModelColNameEnvironment       = "environments"
	ModelColNameDependencySetting = "dependency_settings"
)

type ModelWithTags interface {
	Model
	SetTags(tags []Tag)
	GetTags() (tags []Tag)
}

type ModelWithNameDescription interface {
	Model
	GetName() (name string)
	SetName(name string)
	GetDescription() (description string)
	SetDescription(description string)
}

type ModelWithKey interface {
	Model
	GetKey() (key string)
	SetKey(key string)
}

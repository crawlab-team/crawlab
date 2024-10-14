package models

import (
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseModelV2[T any] struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	CreatedAt time.Time          `json:"created_ts,omitempty" bson:"created_ts,omitempty"`
	CreatedBy primitive.ObjectID `json:"created_by,omitempty" bson:"created_by,omitempty"`
	UpdatedAt time.Time          `json:"updated_ts,omitempty" bson:"updated_ts,omitempty"`
	UpdatedBy primitive.ObjectID `json:"updated_by,omitempty" bson:"updated_by,omitempty"`
}

func (m *BaseModelV2[T]) GetId() primitive.ObjectID {
	return m.Id
}

func (m *BaseModelV2[T]) SetId(id primitive.ObjectID) {
	m.Id = id
}

func (m *BaseModelV2[T]) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m *BaseModelV2[T]) SetCreatedAt(t time.Time) {
	m.CreatedAt = t
}

func (m *BaseModelV2[T]) GetCreatedBy() primitive.ObjectID {
	return m.CreatedBy
}

func (m *BaseModelV2[T]) SetCreatedBy(id primitive.ObjectID) {
	m.CreatedBy = id
}

func (m *BaseModelV2[T]) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}

func (m *BaseModelV2[T]) SetUpdatedAt(t time.Time) {
	m.UpdatedAt = t
}

func (m *BaseModelV2[T]) GetUpdatedBy() primitive.ObjectID {
	return m.UpdatedBy
}

func (m *BaseModelV2[T]) SetUpdatedBy(id primitive.ObjectID) {
	m.UpdatedBy = id
}

func (m *BaseModelV2[T]) SetCreated(id primitive.ObjectID) {
	m.SetCreatedAt(time.Now())
	m.SetCreatedBy(id)
}

func (m *BaseModelV2[T]) SetUpdated(id primitive.ObjectID) {
	m.SetUpdatedAt(time.Now())
	m.SetUpdatedBy(id)
}

func GetModelInstances() []any {
	return []any{
		*new(TestModelV2),
		*new(DataCollectionV2),
		*new(DatabaseV2),
		*new(DatabaseMetricV2),
		*new(DependencyV2),
		*new(DependencyLogV2),
		*new(DependencySettingV2),
		*new(DependencyTaskV2),
		*new(EnvironmentV2),
		*new(GitV2),
		*new(MetricV2),
		*new(NodeV2),
		*new(NotificationChannelV2),
		*new(NotificationRequestV2),
		*new(NotificationSettingV2),
		*new(PermissionV2),
		*new(ProjectV2),
		*new(RolePermissionV2),
		*new(RoleV2),
		*new(ScheduleV2),
		*new(SettingV2),
		*new(SpiderV2),
		*new(SpiderStatV2),
		*new(TaskQueueItemV2),
		*new(TaskStatV2),
		*new(TaskV2),
		*new(TokenV2),
		*new(UserRoleV2),
		*new(UserV2),
	}
}

func GetSystemModelColNames() []string {
	colNames := make([]string, 0)
	for _, instance := range GetModelInstances() {
		colName := GetCollectionNameByInstance(instance)
		if colName != "" {
			colNames = append(colNames, colName)
		}
	}
	return colNames
}

func GetCollectionNameByInstance(v any) string {
	t := reflect.TypeOf(v)
	field := t.Field(0)
	return field.Tag.Get("collection")
}

// Add this new function
func GetSystemModelColNamesMap() map[string]bool {
	colNamesMap := make(map[string]bool)
	for _, instance := range GetModelInstances() {
		colName := GetCollectionNameByInstance(instance)
		if colName != "" {
			colNamesMap[colName] = true
		}
	}
	return colNamesMap
}

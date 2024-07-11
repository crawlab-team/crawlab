package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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

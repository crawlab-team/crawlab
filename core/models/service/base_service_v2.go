package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"sync"

	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	instanceMap    = make(map[string]any)
	onceMap        = make(map[string]*sync.Once)
	onceColNameMap = make(map[string]*sync.Once)
	mu             sync.Mutex
)

type ModelServiceV2[T any] struct {
	col *mongo.Col
}

func (svc *ModelServiceV2[T]) GetById(id primitive.ObjectID) (model *T, err error) {
	var result T
	err = svc.col.FindId(id).One(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (svc *ModelServiceV2[T]) GetByIdContext(ctx context.Context, id primitive.ObjectID) (model *T, err error) {
	var result T
	err = svc.col.GetCollection().FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (svc *ModelServiceV2[T]) GetOne(query bson.M, options *mongo.FindOptions) (model *T, err error) {
	var result T
	err = svc.col.Find(query, options).One(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (svc *ModelServiceV2[T]) GetOneContext(ctx context.Context, query bson.M, opts *mongo.FindOptions) (model *T, err error) {
	var result T
	_opts := &options.FindOneOptions{}
	if opts != nil {
		if opts.Skip != 0 {
			skipInt64 := int64(opts.Skip)
			_opts.Skip = &skipInt64
		}
		if opts.Sort != nil {
			_opts.Sort = opts.Sort
		}
	}
	err = svc.col.GetCollection().FindOne(ctx, query, _opts).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (svc *ModelServiceV2[T]) GetMany(query bson.M, options *mongo.FindOptions) (models []T, err error) {
	var result []T
	err = svc.col.Find(query, options).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (svc *ModelServiceV2[T]) GetManyContext(ctx context.Context, query bson.M, opts *mongo.FindOptions) (models []T, err error) {
	var result []T
	_opts := &options.FindOptions{}
	if opts != nil {
		if opts.Skip != 0 {
			skipInt64 := int64(opts.Skip)
			_opts.Skip = &skipInt64
		}
		if opts.Limit != 0 {
			limitInt64 := int64(opts.Limit)
			_opts.Limit = &limitInt64
		}
		if opts.Sort != nil {
			_opts.Sort = opts.Sort
		}
	}
	cur, err := svc.col.GetCollection().Find(ctx, query, _opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var model T
		if err := cur.Decode(&model); err != nil {
			return nil, err
		}
		result = append(result, model)
	}
	return result, nil
}

func (svc *ModelServiceV2[T]) DeleteById(id primitive.ObjectID) (err error) {
	return svc.col.DeleteId(id)
}

func (svc *ModelServiceV2[T]) DeleteByIdContext(ctx context.Context, id primitive.ObjectID) (err error) {
	_, err = svc.col.GetCollection().DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (svc *ModelServiceV2[T]) DeleteOne(query bson.M) (err error) {
	_, err = svc.col.GetCollection().DeleteOne(svc.col.GetContext(), query)
	return err
}

func (svc *ModelServiceV2[T]) DeleteOneContext(ctx context.Context, query bson.M) (err error) {
	_, err = svc.col.GetCollection().DeleteOne(ctx, query)
	return err
}

func (svc *ModelServiceV2[T]) DeleteMany(query bson.M) (err error) {
	_, err = svc.col.GetCollection().DeleteMany(svc.col.GetContext(), query, nil)
	return err
}

func (svc *ModelServiceV2[T]) DeleteManyContext(ctx context.Context, query bson.M) (err error) {
	_, err = svc.col.GetCollection().DeleteMany(ctx, query, nil)
	return err
}

func (svc *ModelServiceV2[T]) UpdateById(id primitive.ObjectID, update bson.M) (err error) {
	return svc.col.UpdateId(id, update)
}

func (svc *ModelServiceV2[T]) UpdateByIdContext(ctx context.Context, id primitive.ObjectID, update bson.M) (err error) {
	_, err = svc.col.GetCollection().UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (svc *ModelServiceV2[T]) UpdateOne(query bson.M, update bson.M) (err error) {
	_, err = svc.col.GetCollection().UpdateOne(svc.col.GetContext(), query, update)
	return err
}

func (svc *ModelServiceV2[T]) UpdateOneContext(ctx context.Context, query bson.M, update bson.M) (err error) {
	_, err = svc.col.GetCollection().UpdateOne(ctx, query, update)
	return err
}

func (svc *ModelServiceV2[T]) UpdateMany(query bson.M, update bson.M) (err error) {
	_, err = svc.col.GetCollection().UpdateMany(svc.col.GetContext(), query, update)
	return err
}

func (svc *ModelServiceV2[T]) UpdateManyContext(ctx context.Context, query bson.M, update bson.M) (err error) {
	_, err = svc.col.GetCollection().UpdateMany(ctx, query, update)
	return err
}

func (svc *ModelServiceV2[T]) ReplaceById(id primitive.ObjectID, model T) (err error) {
	_, err = svc.col.GetCollection().ReplaceOne(svc.col.GetContext(), bson.M{"_id": id}, model)
	return err
}

func (svc *ModelServiceV2[T]) ReplaceByIdContext(ctx context.Context, id primitive.ObjectID, model T) (err error) {
	_, err = svc.col.GetCollection().ReplaceOne(ctx, bson.M{"_id": id}, model)
	return err
}

func (svc *ModelServiceV2[T]) ReplaceOne(query bson.M, model T) (err error) {
	_, err = svc.col.GetCollection().ReplaceOne(svc.col.GetContext(), query, model)
	return err
}

func (svc *ModelServiceV2[T]) ReplaceOneContext(ctx context.Context, query bson.M, model T) (err error) {
	_, err = svc.col.GetCollection().ReplaceOne(ctx, query, model)
	return err
}

func (svc *ModelServiceV2[T]) InsertOne(model T) (id primitive.ObjectID, err error) {
	m := any(&model).(interfaces.Model)
	if m.GetId().IsZero() {
		m.SetId(primitive.NewObjectID())
	}
	res, err := svc.col.GetCollection().InsertOne(svc.col.GetContext(), m)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (svc *ModelServiceV2[T]) InsertOneContext(ctx context.Context, model T) (id primitive.ObjectID, err error) {
	m := any(&model).(interfaces.Model)
	if m.GetId().IsZero() {
		m.SetId(primitive.NewObjectID())
	}
	res, err := svc.col.GetCollection().InsertOne(ctx, m)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (svc *ModelServiceV2[T]) InsertMany(models []T) (ids []primitive.ObjectID, err error) {
	var _models []any
	for _, model := range models {
		m := any(&model).(interfaces.Model)
		if m.GetId().IsZero() {
			m.SetId(primitive.NewObjectID())
		}
		_models = append(_models, m)
	}
	res, err := svc.col.GetCollection().InsertMany(svc.col.GetContext(), _models)
	if err != nil {
		return nil, err
	}
	for _, v := range res.InsertedIDs {
		ids = append(ids, v.(primitive.ObjectID))
	}
	return ids, nil
}

func (svc *ModelServiceV2[T]) InsertManyContext(ctx context.Context, models []T) (ids []primitive.ObjectID, err error) {
	var _models []any
	for _, model := range models {
		m := any(&model).(interfaces.Model)
		if m.GetId().IsZero() {
			m.SetId(primitive.NewObjectID())
		}
		_models = append(_models, m)
	}
	res, err := svc.col.GetCollection().InsertMany(ctx, _models)
	if err != nil {
		return nil, err
	}
	for _, v := range res.InsertedIDs {
		ids = append(ids, v.(primitive.ObjectID))
	}
	return ids, nil
}

func (svc *ModelServiceV2[T]) Count(query bson.M) (total int, err error) {
	return svc.col.Count(query)
}

func (svc *ModelServiceV2[T]) GetCol() (col *mongo.Col) {
	return svc.col
}

func GetCollectionNameByInstance(v any) string {
	t := reflect.TypeOf(v)
	field := t.Field(0)
	return field.Tag.Get("collection")
}

func getCollectionName[T any]() string {
	var instance T
	t := reflect.TypeOf(instance)
	field := t.Field(0)
	return field.Tag.Get("collection")
}

// NewModelServiceV2 return singleton instance of ModelServiceV2
func NewModelServiceV2[T any]() *ModelServiceV2[T] {
	typeName := fmt.Sprintf("%T", *new(T))

	mu.Lock()
	defer mu.Unlock()

	if _, exists := onceMap[typeName]; !exists {
		onceMap[typeName] = new(sync.Once)
	}

	var instance *ModelServiceV2[T]

	onceMap[typeName].Do(func() {
		collectionName := getCollectionName[T]()
		collection := mongo.GetMongoCol(collectionName)
		instance = &ModelServiceV2[T]{col: collection}
		instanceMap[typeName] = instance
	})

	return instanceMap[typeName].(*ModelServiceV2[T])
}

func NewModelServiceV2WithColName[T any](colName string) *ModelServiceV2[T] {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := onceColNameMap[colName]; !exists {
		onceColNameMap[colName] = new(sync.Once)
	}

	var instance *ModelServiceV2[T]

	onceColNameMap[colName].Do(func() {
		collection := mongo.GetMongoCol(colName)
		instance = &ModelServiceV2[T]{col: collection}
		instanceMap[colName] = instance
	})

	return instanceMap[colName].(*ModelServiceV2[T])
}

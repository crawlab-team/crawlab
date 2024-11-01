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

type ModelService[T any] struct {
	col *mongo.Col
}

func (svc *ModelService[T]) GetById(id primitive.ObjectID) (model *T, err error) {
	var result T
	err = svc.col.FindId(id).One(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (svc *ModelService[T]) GetByIdContext(ctx context.Context, id primitive.ObjectID) (model *T, err error) {
	var result T
	err = svc.col.GetCollection().FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (svc *ModelService[T]) GetOne(query bson.M, options *mongo.FindOptions) (model *T, err error) {
	var result T
	err = svc.col.Find(query, options).One(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (svc *ModelService[T]) GetOneContext(ctx context.Context, query bson.M, opts *mongo.FindOptions) (model *T, err error) {
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

func (svc *ModelService[T]) GetMany(query bson.M, options *mongo.FindOptions) (models []T, err error) {
	var result []T
	err = svc.col.Find(query, options).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (svc *ModelService[T]) GetManyContext(ctx context.Context, query bson.M, opts *mongo.FindOptions) (models []T, err error) {
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

func (svc *ModelService[T]) DeleteById(id primitive.ObjectID) (err error) {
	return svc.col.DeleteId(id)
}

func (svc *ModelService[T]) DeleteByIdContext(ctx context.Context, id primitive.ObjectID) (err error) {
	_, err = svc.col.GetCollection().DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (svc *ModelService[T]) DeleteOne(query bson.M) (err error) {
	_, err = svc.col.GetCollection().DeleteOne(svc.col.GetContext(), query)
	return err
}

func (svc *ModelService[T]) DeleteOneContext(ctx context.Context, query bson.M) (err error) {
	_, err = svc.col.GetCollection().DeleteOne(ctx, query)
	return err
}

func (svc *ModelService[T]) DeleteMany(query bson.M) (err error) {
	_, err = svc.col.GetCollection().DeleteMany(svc.col.GetContext(), query, nil)
	return err
}

func (svc *ModelService[T]) DeleteManyContext(ctx context.Context, query bson.M) (err error) {
	_, err = svc.col.GetCollection().DeleteMany(ctx, query, nil)
	return err
}

func (svc *ModelService[T]) UpdateById(id primitive.ObjectID, update bson.M) (err error) {
	return svc.col.UpdateId(id, update)
}

func (svc *ModelService[T]) UpdateByIdContext(ctx context.Context, id primitive.ObjectID, update bson.M) (err error) {
	_, err = svc.col.GetCollection().UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (svc *ModelService[T]) UpdateOne(query bson.M, update bson.M) (err error) {
	_, err = svc.col.GetCollection().UpdateOne(svc.col.GetContext(), query, update)
	return err
}

func (svc *ModelService[T]) UpdateOneContext(ctx context.Context, query bson.M, update bson.M) (err error) {
	_, err = svc.col.GetCollection().UpdateOne(ctx, query, update)
	return err
}

func (svc *ModelService[T]) UpdateMany(query bson.M, update bson.M) (err error) {
	_, err = svc.col.GetCollection().UpdateMany(svc.col.GetContext(), query, update)
	return err
}

func (svc *ModelService[T]) UpdateManyContext(ctx context.Context, query bson.M, update bson.M) (err error) {
	_, err = svc.col.GetCollection().UpdateMany(ctx, query, update)
	return err
}

func (svc *ModelService[T]) ReplaceById(id primitive.ObjectID, model T) (err error) {
	_, err = svc.col.GetCollection().ReplaceOne(svc.col.GetContext(), bson.M{"_id": id}, model)
	return err
}

func (svc *ModelService[T]) ReplaceByIdContext(ctx context.Context, id primitive.ObjectID, model T) (err error) {
	_, err = svc.col.GetCollection().ReplaceOne(ctx, bson.M{"_id": id}, model)
	return err
}

func (svc *ModelService[T]) ReplaceOne(query bson.M, model T) (err error) {
	_, err = svc.col.GetCollection().ReplaceOne(svc.col.GetContext(), query, model)
	return err
}

func (svc *ModelService[T]) ReplaceOneContext(ctx context.Context, query bson.M, model T) (err error) {
	_, err = svc.col.GetCollection().ReplaceOne(ctx, query, model)
	return err
}

func (svc *ModelService[T]) InsertOne(model T) (id primitive.ObjectID, err error) {
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

func (svc *ModelService[T]) InsertOneContext(ctx context.Context, model T) (id primitive.ObjectID, err error) {
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

func (svc *ModelService[T]) InsertMany(models []T) (ids []primitive.ObjectID, err error) {
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

func (svc *ModelService[T]) InsertManyContext(ctx context.Context, models []T) (ids []primitive.ObjectID, err error) {
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

func (svc *ModelService[T]) Count(query bson.M) (total int, err error) {
	return svc.col.Count(query)
}

func (svc *ModelService[T]) GetCol() (col *mongo.Col) {
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

// NewModelService return singleton instance of ModelService
func NewModelService[T any]() *ModelService[T] {
	typeName := fmt.Sprintf("%T", *new(T))

	mu.Lock()
	defer mu.Unlock()

	if _, exists := onceMap[typeName]; !exists {
		onceMap[typeName] = new(sync.Once)
	}

	var instance *ModelService[T]

	onceMap[typeName].Do(func() {
		collectionName := getCollectionName[T]()
		collection := mongo.GetMongoCol(collectionName)
		instance = &ModelService[T]{col: collection}
		instanceMap[typeName] = instance
	})

	return instanceMap[typeName].(*ModelService[T])
}

func NewModelServiceWithColName[T any](colName string) *ModelService[T] {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := onceColNameMap[colName]; !exists {
		onceColNameMap[colName] = new(sync.Once)
	}

	var instance *ModelService[T]

	onceColNameMap[colName].Do(func() {
		collection := mongo.GetMongoCol(colName)
		instance = &ModelService[T]{col: collection}
		instanceMap[colName] = instance
	})

	return instanceMap[colName].(*ModelService[T])
}

package client

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab/core/grpc/client"
	"github.com/crawlab-team/crawlab/core/interfaces"
	nodeconfig "github.com/crawlab-team/crawlab/core/node/config"
	"github.com/crawlab-team/crawlab/db/mongo"
	grpc "github.com/crawlab-team/crawlab/grpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"sync"
)

var (
	instanceMap = make(map[string]interface{})
	onceMap     = make(map[string]*sync.Once)
	mu          sync.Mutex
)

type ModelServiceV2[T any] struct {
	cfg       interfaces.NodeConfigService
	c         *client.GrpcClientV2
	modelType string
}

func (svc *ModelServiceV2[T]) GetById(id primitive.ObjectID) (model *T, err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	res, err := svc.c.ModelBaseServiceV2Client.GetById(ctx, &grpc.ModelServiceV2GetByIdRequest{
		NodeKey:   svc.cfg.GetNodeKey(),
		ModelType: svc.modelType,
		Id:        id.Hex(),
	})
	if err != nil {
		return nil, err
	}
	return svc.deserializeOne(res)
}

func (svc *ModelServiceV2[T]) GetOne(query bson.M, options *mongo.FindOptions) (model *T, err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	queryData, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	findOptionsData, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	res, err := svc.c.ModelBaseServiceV2Client.GetOne(ctx, &grpc.ModelServiceV2GetOneRequest{
		NodeKey:     svc.cfg.GetNodeKey(),
		ModelType:   svc.modelType,
		Query:       queryData,
		FindOptions: findOptionsData,
	})
	if err != nil {
		return nil, err
	}
	return svc.deserializeOne(res)
}

func (svc *ModelServiceV2[T]) GetMany(query bson.M, options *mongo.FindOptions) (models []T, err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	queryData, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	findOptionsData, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	res, err := svc.c.ModelBaseServiceV2Client.GetMany(ctx, &grpc.ModelServiceV2GetManyRequest{
		NodeKey:     svc.cfg.GetNodeKey(),
		ModelType:   svc.modelType,
		Query:       queryData,
		FindOptions: findOptionsData,
	})
	if err != nil {
		return nil, err
	}
	return svc.deserializeMany(res)
}

func (svc *ModelServiceV2[T]) DeleteById(id primitive.ObjectID) (err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	_, err = svc.c.ModelBaseServiceV2Client.DeleteById(ctx, &grpc.ModelServiceV2DeleteByIdRequest{
		NodeKey:   svc.cfg.GetNodeKey(),
		ModelType: svc.modelType,
		Id:        id.Hex(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (svc *ModelServiceV2[T]) DeleteOne(query bson.M) (err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	queryData, err := json.Marshal(query)
	if err != nil {
		return err
	}
	_, err = svc.c.ModelBaseServiceV2Client.DeleteOne(ctx, &grpc.ModelServiceV2DeleteOneRequest{
		NodeKey:   svc.cfg.GetNodeKey(),
		ModelType: svc.modelType,
		Query:     queryData,
	})
	if err != nil {
		return err
	}
	return nil
}

func (svc *ModelServiceV2[T]) DeleteMany(query bson.M) (err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	queryData, err := json.Marshal(query)
	if err != nil {
		return err
	}
	_, err = svc.c.ModelBaseServiceV2Client.DeleteMany(ctx, &grpc.ModelServiceV2DeleteManyRequest{
		NodeKey:   svc.cfg.GetNodeKey(),
		ModelType: svc.modelType,
		Query:     queryData,
	})
	if err != nil {
		return err
	}
	return nil
}

func (svc *ModelServiceV2[T]) UpdateById(id primitive.ObjectID, update bson.M) (err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	updateData, err := json.Marshal(update)
	if err != nil {
		return err
	}
	_, err = svc.c.ModelBaseServiceV2Client.UpdateById(ctx, &grpc.ModelServiceV2UpdateByIdRequest{
		NodeKey:   svc.cfg.GetNodeKey(),
		ModelType: svc.modelType,
		Id:        id.Hex(),
		Update:    updateData,
	})
	if err != nil {
		return err
	}
	return nil
}

func (svc *ModelServiceV2[T]) UpdateOne(query bson.M, update bson.M) (err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	queryData, err := json.Marshal(query)
	if err != nil {
		return err
	}
	updateData, err := json.Marshal(update)
	if err != nil {
		return err
	}
	_, err = svc.c.ModelBaseServiceV2Client.UpdateOne(ctx, &grpc.ModelServiceV2UpdateOneRequest{
		NodeKey:   svc.cfg.GetNodeKey(),
		ModelType: svc.modelType,
		Query:     queryData,
		Update:    updateData,
	})
	if err != nil {
		return err
	}
	return nil
}

func (svc *ModelServiceV2[T]) UpdateMany(query bson.M, update bson.M) (err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	queryData, err := json.Marshal(query)
	if err != nil {
		return err
	}
	updateData, err := json.Marshal(update)
	if err != nil {
		return err
	}
	_, err = svc.c.ModelBaseServiceV2Client.UpdateMany(ctx, &grpc.ModelServiceV2UpdateManyRequest{
		NodeKey:   svc.cfg.GetNodeKey(),
		ModelType: svc.modelType,
		Query:     queryData,
		Update:    updateData,
	})
	return nil
}

func (svc *ModelServiceV2[T]) ReplaceById(id primitive.ObjectID, model T) (err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	modelData, err := json.Marshal(model)
	if err != nil {
		return err
	}
	_, err = svc.c.ModelBaseServiceV2Client.ReplaceById(ctx, &grpc.ModelServiceV2ReplaceByIdRequest{
		NodeKey:   svc.cfg.GetNodeKey(),
		ModelType: svc.modelType,
		Id:        id.Hex(),
		Model:     modelData,
	})
	if err != nil {
		return err
	}
	return nil
}

func (svc *ModelServiceV2[T]) ReplaceOne(query bson.M, model T) (err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	queryData, err := json.Marshal(query)
	if err != nil {
		return err
	}
	modelData, err := json.Marshal(model)
	if err != nil {
		return err
	}
	_, err = svc.c.ModelBaseServiceV2Client.ReplaceOne(ctx, &grpc.ModelServiceV2ReplaceOneRequest{
		NodeKey:   svc.cfg.GetNodeKey(),
		ModelType: svc.modelType,
		Query:     queryData,
		Model:     modelData,
	})
	if err != nil {
		return err
	}
	return nil
}

func (svc *ModelServiceV2[T]) InsertOne(model T) (id primitive.ObjectID, err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	modelData, err := json.Marshal(model)
	if err != nil {
		return primitive.NilObjectID, err
	}
	res, err := svc.c.ModelBaseServiceV2Client.InsertOne(ctx, &grpc.ModelServiceV2InsertOneRequest{
		NodeKey:   svc.cfg.GetNodeKey(),
		ModelType: svc.modelType,
		Model:     modelData,
	})
	if err != nil {
		return primitive.NilObjectID, err
	}
	return deserialize[primitive.ObjectID](res)
	//idStr, err := deserialize[string](res)
	//if err != nil {
	//	return primitive.NilObjectID, err
	//}
	//return primitive.ObjectIDFromHex(idStr)
}

func (svc *ModelServiceV2[T]) InsertMany(models []T) (ids []primitive.ObjectID, err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	modelsData, err := json.Marshal(models)
	if err != nil {
		return nil, err
	}
	res, err := svc.c.ModelBaseServiceV2Client.InsertMany(ctx, &grpc.ModelServiceV2InsertManyRequest{
		NodeKey:   svc.cfg.GetNodeKey(),
		ModelType: svc.modelType,
		Models:    modelsData,
	})
	if err != nil {
		return nil, err
	}
	return deserialize[[]primitive.ObjectID](res)
}

func (svc *ModelServiceV2[T]) Count(query bson.M) (total int, err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	queryData, err := json.Marshal(query)
	if err != nil {
		return 0, err
	}
	res, err := svc.c.ModelBaseServiceV2Client.Count(ctx, &grpc.ModelServiceV2CountRequest{
		NodeKey:   svc.cfg.GetNodeKey(),
		ModelType: svc.modelType,
		Query:     queryData,
	})
	if err != nil {
		return 0, err
	}
	return deserialize[int](res)
}

func (svc *ModelServiceV2[T]) GetCol() (col *mongo.Col) {
	return nil
}

func (svc *ModelServiceV2[T]) deserializeOne(res *grpc.Response) (result *T, err error) {
	r, err := deserialize[T](res)
	if err != nil {
		return nil, err
	}
	return &r, err
}

func (svc *ModelServiceV2[T]) deserializeMany(res *grpc.Response) (results []T, err error) {
	return deserialize[[]T](res)
}

func deserialize[T any](res *grpc.Response) (result T, err error) {
	err = json.Unmarshal(res.Data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func NewModelServiceV2[T any]() *ModelServiceV2[T] {
	mu.Lock()
	defer mu.Unlock()

	var v T
	t := reflect.TypeOf(v)
	typeName := t.Name()

	if _, exists := onceMap[typeName]; !exists {
		onceMap[typeName] = &sync.Once{}
	}

	var instance *ModelServiceV2[T]

	c, err := client.GetGrpcClientV2()
	if err != nil {
		panic(err)
	}
	if !c.IsStarted() {
		err = c.Start()
		if err != nil {
			panic(err)
		}
	}

	onceMap[typeName].Do(func() {
		instance = &ModelServiceV2[T]{
			cfg:       nodeconfig.GetNodeConfigService(),
			c:         c,
			modelType: typeName,
		}
		instanceMap[typeName] = instance
	})

	return instanceMap[typeName].(*ModelServiceV2[T])
}

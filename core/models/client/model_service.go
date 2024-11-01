package client

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab/core/grpc/client"
	"github.com/crawlab-team/crawlab/core/interfaces"
	nodeconfig "github.com/crawlab-team/crawlab/core/node/config"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/crawlab-team/crawlab/grpc"
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

type ModelService[T any] struct {
	cfg       interfaces.NodeConfigService
	c         *client.GrpcClient
	modelType string
}

func (svc *ModelService[T]) GetById(id primitive.ObjectID) (model *T, err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	res, err := svc.c.ModelBaseServiceClient.GetById(ctx, &grpc.ModelServiceGetByIdRequest{
		NodeKey:   svc.cfg.GetNodeKey(),
		ModelType: svc.modelType,
		Id:        id.Hex(),
	})
	if err != nil {
		return nil, err
	}
	return svc.deserializeOne(res)
}

func (svc *ModelService[T]) GetOne(query bson.M, options *mongo.FindOptions) (model *T, err error) {
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
	res, err := svc.c.ModelBaseServiceClient.GetOne(ctx, &grpc.ModelServiceGetOneRequest{
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

func (svc *ModelService[T]) GetMany(query bson.M, options *mongo.FindOptions) (models []T, err error) {
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
	res, err := svc.c.ModelBaseServiceClient.GetMany(ctx, &grpc.ModelServiceGetManyRequest{
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

func (svc *ModelService[T]) DeleteById(id primitive.ObjectID) (err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	_, err = svc.c.ModelBaseServiceClient.DeleteById(ctx, &grpc.ModelServiceDeleteByIdRequest{
		NodeKey:   svc.cfg.GetNodeKey(),
		ModelType: svc.modelType,
		Id:        id.Hex(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (svc *ModelService[T]) DeleteOne(query bson.M) (err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	queryData, err := json.Marshal(query)
	if err != nil {
		return err
	}
	_, err = svc.c.ModelBaseServiceClient.DeleteOne(ctx, &grpc.ModelServiceDeleteOneRequest{
		NodeKey:   svc.cfg.GetNodeKey(),
		ModelType: svc.modelType,
		Query:     queryData,
	})
	if err != nil {
		return err
	}
	return nil
}

func (svc *ModelService[T]) DeleteMany(query bson.M) (err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	queryData, err := json.Marshal(query)
	if err != nil {
		return err
	}
	_, err = svc.c.ModelBaseServiceClient.DeleteMany(ctx, &grpc.ModelServiceDeleteManyRequest{
		NodeKey:   svc.cfg.GetNodeKey(),
		ModelType: svc.modelType,
		Query:     queryData,
	})
	if err != nil {
		return err
	}
	return nil
}

func (svc *ModelService[T]) UpdateById(id primitive.ObjectID, update bson.M) (err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	updateData, err := json.Marshal(update)
	if err != nil {
		return err
	}
	_, err = svc.c.ModelBaseServiceClient.UpdateById(ctx, &grpc.ModelServiceUpdateByIdRequest{
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

func (svc *ModelService[T]) UpdateOne(query bson.M, update bson.M) (err error) {
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
	_, err = svc.c.ModelBaseServiceClient.UpdateOne(ctx, &grpc.ModelServiceUpdateOneRequest{
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

func (svc *ModelService[T]) UpdateMany(query bson.M, update bson.M) (err error) {
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
	_, err = svc.c.ModelBaseServiceClient.UpdateMany(ctx, &grpc.ModelServiceUpdateManyRequest{
		NodeKey:   svc.cfg.GetNodeKey(),
		ModelType: svc.modelType,
		Query:     queryData,
		Update:    updateData,
	})
	return nil
}

func (svc *ModelService[T]) ReplaceById(id primitive.ObjectID, model T) (err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	modelData, err := json.Marshal(model)
	if err != nil {
		return err
	}
	_, err = svc.c.ModelBaseServiceClient.ReplaceById(ctx, &grpc.ModelServiceReplaceByIdRequest{
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

func (svc *ModelService[T]) ReplaceOne(query bson.M, model T) (err error) {
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
	_, err = svc.c.ModelBaseServiceClient.ReplaceOne(ctx, &grpc.ModelServiceReplaceOneRequest{
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

func (svc *ModelService[T]) InsertOne(model T) (id primitive.ObjectID, err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	modelData, err := json.Marshal(model)
	if err != nil {
		return primitive.NilObjectID, err
	}
	res, err := svc.c.ModelBaseServiceClient.InsertOne(ctx, &grpc.ModelServiceInsertOneRequest{
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

func (svc *ModelService[T]) InsertMany(models []T) (ids []primitive.ObjectID, err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	modelsData, err := json.Marshal(models)
	if err != nil {
		return nil, err
	}
	res, err := svc.c.ModelBaseServiceClient.InsertMany(ctx, &grpc.ModelServiceInsertManyRequest{
		NodeKey:   svc.cfg.GetNodeKey(),
		ModelType: svc.modelType,
		Models:    modelsData,
	})
	if err != nil {
		return nil, err
	}
	return deserialize[[]primitive.ObjectID](res)
}

func (svc *ModelService[T]) Count(query bson.M) (total int, err error) {
	ctx, cancel := svc.c.Context()
	defer cancel()
	queryData, err := json.Marshal(query)
	if err != nil {
		return 0, err
	}
	res, err := svc.c.ModelBaseServiceClient.Count(ctx, &grpc.ModelServiceCountRequest{
		NodeKey:   svc.cfg.GetNodeKey(),
		ModelType: svc.modelType,
		Query:     queryData,
	})
	if err != nil {
		return 0, err
	}
	return deserialize[int](res)
}

func (svc *ModelService[T]) GetCol() (col *mongo.Col) {
	return nil
}

func (svc *ModelService[T]) deserializeOne(res *grpc.Response) (result *T, err error) {
	r, err := deserialize[T](res)
	if err != nil {
		return nil, err
	}
	return &r, err
}

func (svc *ModelService[T]) deserializeMany(res *grpc.Response) (results []T, err error) {
	return deserialize[[]T](res)
}

func deserialize[T any](res *grpc.Response) (result T, err error) {
	err = json.Unmarshal(res.Data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func NewModelService[T any]() *ModelService[T] {
	mu.Lock()
	defer mu.Unlock()

	var v T
	t := reflect.TypeOf(v)
	typeName := t.Name()

	if _, exists := onceMap[typeName]; !exists {
		onceMap[typeName] = new(sync.Once)
	}

	var instance *ModelService[T]

	c := client.GetGrpcClient()
	if !c.IsStarted() {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}

	onceMap[typeName].Do(func() {
		instance = &ModelService[T]{
			cfg:       nodeconfig.GetNodeConfigService(),
			c:         c,
			modelType: typeName,
		}
		instanceMap[typeName] = instance
	})

	return instanceMap[typeName].(*ModelService[T])
}

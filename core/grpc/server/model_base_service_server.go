package server

import (
	"context"
	"encoding/json"
	"reflect"

	models2 "github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/crawlab-team/crawlab/grpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	typeNameColNameMap  = make(map[string]string)
	typeOneNameModelMap = make(map[string]any)
	typeOneInstances    = models2.GetModelInstances()
)

func init() {
	for _, v := range typeOneInstances {
		t := reflect.TypeOf(v)
		typeName := t.Name()
		colName := service.GetCollectionNameByInstance(v)
		typeNameColNameMap[typeName] = colName
		typeOneNameModelMap[typeName] = v
	}
}

func GetOneInstanceModel(typeName string) any {
	return typeOneNameModelMap[typeName]
}

type ModelBaseServiceServer struct {
	grpc.UnimplementedModelBaseServiceServer
}

func (svr ModelBaseServiceServer) GetById(_ context.Context, req *grpc.ModelServiceGetByIdRequest) (res *grpc.Response, err error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return HandleError(err)
	}
	modelSvc := service.NewModelServiceWithColName[bson.M](typeNameColNameMap[req.ModelType])
	data, err := modelSvc.GetById(id)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccessWithData(data)
}

func (svr ModelBaseServiceServer) GetOne(_ context.Context, req *grpc.ModelServiceGetOneRequest) (res *grpc.Response, err error) {
	var query bson.M
	err = json.Unmarshal(req.Query, &query)
	if err != nil {
		return HandleError(err)
	}
	var options mongo.FindOptions
	err = json.Unmarshal(req.FindOptions, &options)
	if err != nil {
		return HandleError(err)
	}
	modelSvc := service.NewModelServiceWithColName[bson.M](typeNameColNameMap[req.ModelType])
	data, err := modelSvc.GetOne(query, &options)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccessWithData(data)
}

func (svr ModelBaseServiceServer) GetMany(_ context.Context, req *grpc.ModelServiceGetManyRequest) (res *grpc.Response, err error) {
	var query bson.M
	err = json.Unmarshal(req.Query, &query)
	if err != nil {
		return HandleError(err)
	}
	var options mongo.FindOptions
	err = json.Unmarshal(req.FindOptions, &options)
	if err != nil {
		return HandleError(err)
	}
	modelSvc := service.NewModelServiceWithColName[bson.M](typeNameColNameMap[req.ModelType])
	data, err := modelSvc.GetMany(query, &options)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccessWithData(data)
}

func (svr ModelBaseServiceServer) DeleteById(_ context.Context, req *grpc.ModelServiceDeleteByIdRequest) (res *grpc.Response, err error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return HandleError(err)
	}
	modelSvc := GetModelService[bson.M](req.ModelType)
	err = modelSvc.DeleteById(id)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccess()
}

func (svr ModelBaseServiceServer) DeleteOne(_ context.Context, req *grpc.ModelServiceDeleteOneRequest) (res *grpc.Response, err error) {
	var query bson.M
	err = json.Unmarshal(req.Query, &query)
	if err != nil {
		return HandleError(err)
	}
	modelSvc := GetModelService[bson.M](req.ModelType)
	err = modelSvc.DeleteOne(query)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccess()
}

func (svr ModelBaseServiceServer) DeleteMany(_ context.Context, req *grpc.ModelServiceDeleteManyRequest) (res *grpc.Response, err error) {
	var query bson.M
	err = json.Unmarshal(req.Query, &query)
	if err != nil {
		return HandleError(err)
	}
	modelSvc := GetModelService[bson.M](req.ModelType)
	err = modelSvc.DeleteMany(query)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccess()
}

func (svr ModelBaseServiceServer) UpdateById(_ context.Context, req *grpc.ModelServiceUpdateByIdRequest) (res *grpc.Response, err error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return HandleError(err)
	}
	var update bson.M
	err = json.Unmarshal(req.Update, &update)
	if err != nil {
		return HandleError(err)
	}
	modelSvc := GetModelService[bson.M](req.ModelType)
	err = modelSvc.UpdateById(id, update)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccess()
}

func (svr ModelBaseServiceServer) UpdateOne(_ context.Context, req *grpc.ModelServiceUpdateOneRequest) (res *grpc.Response, err error) {
	var query bson.M
	err = json.Unmarshal(req.Query, &query)
	if err != nil {
		return HandleError(err)
	}
	var update bson.M
	err = json.Unmarshal(req.Update, &update)
	if err != nil {
		return HandleError(err)
	}
	modelSvc := GetModelService[bson.M](req.ModelType)
	err = modelSvc.UpdateOne(query, update)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccess()
}

func (svr ModelBaseServiceServer) UpdateMany(_ context.Context, req *grpc.ModelServiceUpdateManyRequest) (res *grpc.Response, err error) {
	var query bson.M
	err = json.Unmarshal(req.Query, &query)
	if err != nil {
		return HandleError(err)
	}
	var update bson.M
	err = json.Unmarshal(req.Update, &update)
	if err != nil {
		return HandleError(err)
	}
	modelSvc := GetModelService[bson.M](req.ModelType)
	err = modelSvc.UpdateMany(query, update)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccess()
}

func (svr ModelBaseServiceServer) ReplaceById(_ context.Context, req *grpc.ModelServiceReplaceByIdRequest) (res *grpc.Response, err error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return HandleError(err)
	}
	model := GetOneInstanceModel(req.ModelType)
	modelType := reflect.TypeOf(model)
	modelValuePtr := reflect.New(modelType).Interface()
	err = json.Unmarshal(req.Model, modelValuePtr)
	if err != nil {
		return HandleError(err)
	}
	modelSvc := GetModelService[bson.M](req.ModelType)
	err = modelSvc.GetCol().ReplaceId(id, modelValuePtr)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccess()
}

func (svr ModelBaseServiceServer) ReplaceOne(_ context.Context, req *grpc.ModelServiceReplaceOneRequest) (res *grpc.Response, err error) {
	var query bson.M
	err = json.Unmarshal(req.Query, &query)
	if err != nil {
		return HandleError(err)
	}
	model := GetOneInstanceModel(req.ModelType)
	modelType := reflect.TypeOf(model)
	modelValuePtr := reflect.New(modelType).Interface()
	err = json.Unmarshal(req.Model, &modelValuePtr)
	if err != nil {
		return HandleError(err)
	}
	modelSvc := GetModelService[bson.M](req.ModelType)
	err = modelSvc.GetCol().Replace(query, modelValuePtr)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccess()
}

func (svr ModelBaseServiceServer) InsertOne(_ context.Context, req *grpc.ModelServiceInsertOneRequest) (res *grpc.Response, err error) {
	model := GetOneInstanceModel(req.ModelType)
	modelType := reflect.TypeOf(model)
	modelValuePtr := reflect.New(modelType).Interface()
	err = json.Unmarshal(req.Model, modelValuePtr)
	if err != nil {
		return HandleError(err)
	}
	modelSvc := GetModelService[bson.M](req.ModelType)
	r, err := modelSvc.GetCol().GetCollection().InsertOne(modelSvc.GetCol().GetContext(), modelValuePtr)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccessWithData(r.InsertedID)
}

func (svr ModelBaseServiceServer) InsertMany(_ context.Context, req *grpc.ModelServiceInsertManyRequest) (res *grpc.Response, err error) {
	model := GetOneInstanceModel(req.ModelType)
	modelType := reflect.TypeOf(model)
	modelsSliceType := reflect.SliceOf(modelType)
	modelsSlicePtr := reflect.New(modelsSliceType).Interface()
	err = json.Unmarshal(req.Models, modelsSlicePtr)
	if err != nil {
		return HandleError(err)
	}
	modelsSlice := reflect.ValueOf(modelsSlicePtr).Elem()
	modelsInterface := make([]any, modelsSlice.Len())
	for i := 0; i < modelsSlice.Len(); i++ {
		modelValue := modelsSlice.Index(i)
		if modelValue.FieldByName("Id").Interface().(primitive.ObjectID).IsZero() {
			modelValue.FieldByName("Id").Set(reflect.ValueOf(primitive.NewObjectID()))
		}
		modelsInterface[i] = modelValue.Interface()
	}
	modelSvc := GetModelService[bson.M](req.ModelType)
	r, err := modelSvc.GetCol().GetCollection().InsertMany(modelSvc.GetCol().GetContext(), modelsInterface)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccessWithData(r.InsertedIDs)
}

func (svr ModelBaseServiceServer) Count(_ context.Context, req *grpc.ModelServiceCountRequest) (res *grpc.Response, err error) {
	var query bson.M
	err = json.Unmarshal(req.Query, &query)
	if err != nil {
		return HandleError(err)
	}
	count, err := GetModelService[bson.M](req.ModelType).Count(query)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccessWithData(count)
}

func GetModelService[T any](typeName string) *service.ModelService[T] {
	return service.NewModelServiceWithColName[T](typeNameColNameMap[typeName])
}

func NewModelBaseServiceServer() *ModelBaseServiceServer {
	return &ModelBaseServiceServer{}
}

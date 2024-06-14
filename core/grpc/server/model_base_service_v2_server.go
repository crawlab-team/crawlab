package server

import (
	"context"
	"encoding/json"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/db/mongo"
	grpc "github.com/crawlab-team/crawlab/grpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

var (
	typeNameColNameMap  = make(map[string]string)
	typeOneNameModelMap = make(map[string]any)
	typeOneInstances    = []any{
		*new(models.TestModel),
		*new(models.DataCollectionV2),
		*new(models.DataSourceV2),
		*new(models.DependencySettingV2),
		*new(models.EnvironmentV2),
		*new(models.GitV2),
		*new(models.NodeV2),
		*new(models.PermissionV2),
		*new(models.ProjectV2),
		*new(models.RolePermissionV2),
		*new(models.RoleV2),
		*new(models.ScheduleV2),
		*new(models.SettingV2),
		*new(models.SpiderV2),
		*new(models.TaskQueueItemV2),
		*new(models.TaskStatV2),
		*new(models.TaskV2),
		*new(models.TokenV2),
		*new(models.UserRoleV2),
		*new(models.UserV2),
	}
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

type ModelBaseServiceServerV2 struct {
	grpc.UnimplementedModelBaseServiceV2Server
}

func (svr ModelBaseServiceServerV2) GetById(ctx context.Context, req *grpc.ModelServiceV2GetByIdRequest) (res *grpc.Response, err error) {
	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return HandleError(err)
	}
	modelSvc := service.NewModelServiceV2WithColName[bson.M](typeNameColNameMap[req.ModelType])
	data, err := modelSvc.GetById(id)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccessWithData(data)
}

func (svr ModelBaseServiceServerV2) GetOne(ctx context.Context, req *grpc.ModelServiceV2GetOneRequest) (res *grpc.Response, err error) {
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
	modelSvc := service.NewModelServiceV2WithColName[bson.M](typeNameColNameMap[req.ModelType])
	data, err := modelSvc.GetOne(query, &options)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccessWithData(data)
}

func (svr ModelBaseServiceServerV2) GetMany(ctx context.Context, req *grpc.ModelServiceV2GetManyRequest) (res *grpc.Response, err error) {
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
	modelSvc := service.NewModelServiceV2WithColName[bson.M](typeNameColNameMap[req.ModelType])
	data, err := modelSvc.GetMany(query, &options)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccessWithData(data)
}

func (svr ModelBaseServiceServerV2) DeleteById(ctx context.Context, req *grpc.ModelServiceV2DeleteByIdRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) DeleteOne(ctx context.Context, req *grpc.ModelServiceV2DeleteOneRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) DeleteMany(ctx context.Context, req *grpc.ModelServiceV2DeleteManyRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) UpdateById(ctx context.Context, req *grpc.ModelServiceV2UpdateByIdRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) UpdateOne(ctx context.Context, req *grpc.ModelServiceV2UpdateOneRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) UpdateMany(ctx context.Context, req *grpc.ModelServiceV2UpdateManyRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) ReplaceById(ctx context.Context, req *grpc.ModelServiceV2ReplaceByIdRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) ReplaceOne(ctx context.Context, req *grpc.ModelServiceV2ReplaceOneRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) InsertOne(ctx context.Context, req *grpc.ModelServiceV2InsertOneRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) InsertMany(ctx context.Context, req *grpc.ModelServiceV2InsertManyRequest) (res *grpc.Response, err error) {
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
		modelsInterface[i] = modelsSlice.Index(i).Interface()
	}
	modelSvc := GetModelService[bson.M](req.ModelType)
	r, err := modelSvc.GetCol().GetCollection().InsertMany(modelSvc.GetCol().GetContext(), modelsInterface)
	if err != nil {
		return HandleError(err)
	}
	return HandleSuccessWithData(r.InsertedIDs)
}

func (svr ModelBaseServiceServerV2) Count(ctx context.Context, req *grpc.ModelServiceV2CountRequest) (res *grpc.Response, err error) {
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

func GetModelService[T any](typeName string) *service.ModelServiceV2[T] {
	return service.NewModelServiceV2WithColName[T](typeNameColNameMap[typeName])
}

func NewModelBaseServiceV2Server() *ModelBaseServiceServerV2 {
	return &ModelBaseServiceServerV2{}
}

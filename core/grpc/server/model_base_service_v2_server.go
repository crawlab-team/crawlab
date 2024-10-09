package server

import (
	"context"
	"encoding/json"
	models2 "github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/crawlab-team/crawlab/grpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

var (
	typeNameColNameMap  = make(map[string]string)
	typeOneNameModelMap = make(map[string]any)
	typeOneInstances    = []any{
		*new(models2.TestModelV2),
		*new(models2.DataCollectionV2),
		*new(models2.DatabaseV2),
		*new(models2.DatabaseMetricV2),
		*new(models2.DependencyV2),
		*new(models2.DependencyLogV2),
		*new(models2.DependencySettingV2),
		*new(models2.DependencyTaskV2),
		*new(models2.EnvironmentV2),
		*new(models2.GitV2),
		*new(models2.MetricV2),
		*new(models2.NodeV2),
		*new(models2.NotificationChannelV2),
		*new(models2.NotificationRequestV2),
		*new(models2.NotificationSettingV2),
		*new(models2.PermissionV2),
		*new(models2.ProjectV2),
		*new(models2.RolePermissionV2),
		*new(models2.RoleV2),
		*new(models2.ScheduleV2),
		*new(models2.SettingV2),
		*new(models2.SpiderV2),
		*new(models2.SpiderStatV2),
		*new(models2.TaskQueueItemV2),
		*new(models2.TaskStatV2),
		*new(models2.TaskV2),
		*new(models2.TokenV2),
		*new(models2.UserRoleV2),
		*new(models2.UserV2),
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

func (svr ModelBaseServiceServerV2) GetById(_ context.Context, req *grpc.ModelServiceV2GetByIdRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) GetOne(_ context.Context, req *grpc.ModelServiceV2GetOneRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) GetMany(_ context.Context, req *grpc.ModelServiceV2GetManyRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) DeleteById(_ context.Context, req *grpc.ModelServiceV2DeleteByIdRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) DeleteOne(_ context.Context, req *grpc.ModelServiceV2DeleteOneRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) DeleteMany(_ context.Context, req *grpc.ModelServiceV2DeleteManyRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) UpdateById(_ context.Context, req *grpc.ModelServiceV2UpdateByIdRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) UpdateOne(_ context.Context, req *grpc.ModelServiceV2UpdateOneRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) UpdateMany(_ context.Context, req *grpc.ModelServiceV2UpdateManyRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) ReplaceById(_ context.Context, req *grpc.ModelServiceV2ReplaceByIdRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) ReplaceOne(_ context.Context, req *grpc.ModelServiceV2ReplaceOneRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) InsertOne(_ context.Context, req *grpc.ModelServiceV2InsertOneRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) InsertMany(_ context.Context, req *grpc.ModelServiceV2InsertManyRequest) (res *grpc.Response, err error) {
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

func (svr ModelBaseServiceServerV2) Count(_ context.Context, req *grpc.ModelServiceV2CountRequest) (res *grpc.Response, err error) {
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

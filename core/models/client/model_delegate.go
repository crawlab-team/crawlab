package client

import (
	"encoding/json"
	config2 "github.com/crawlab-team/crawlab/core/config"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/grpc/client"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/utils"
	grpc "github.com/crawlab-team/crawlab/grpc"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/spf13/viper"
)

func NewModelDelegate(doc interfaces.Model, opts ...ModelDelegateOption) interfaces.GrpcClientModelDelegate {
	switch doc.(type) {
	case *models.Artifact:
		return newModelDelegate(interfaces.ModelIdArtifact, doc, opts...)
	case *models.Tag:
		return newModelDelegate(interfaces.ModelIdTag, doc, opts...)
	case *models.Node:
		return newModelDelegate(interfaces.ModelIdNode, doc, opts...)
	case *models.Project:
		return newModelDelegate(interfaces.ModelIdProject, doc, opts...)
	case *models.Spider:
		return newModelDelegate(interfaces.ModelIdSpider, doc, opts...)
	case *models.Task:
		return newModelDelegate(interfaces.ModelIdTask, doc, opts...)
	case *models.Job:
		return newModelDelegate(interfaces.ModelIdJob, doc, opts...)
	case *models.Schedule:
		return newModelDelegate(interfaces.ModelIdSchedule, doc, opts...)
	case *models.User:
		return newModelDelegate(interfaces.ModelIdUser, doc, opts...)
	case *models.Setting:
		return newModelDelegate(interfaces.ModelIdSetting, doc, opts...)
	case *models.Token:
		return newModelDelegate(interfaces.ModelIdToken, doc, opts...)
	case *models.Variable:
		return newModelDelegate(interfaces.ModelIdVariable, doc, opts...)
	case *models.TaskQueueItem:
		return newModelDelegate(interfaces.ModelIdTaskQueue, doc, opts...)
	case *models.TaskStat:
		return newModelDelegate(interfaces.ModelIdTaskStat, doc, opts...)
	case *models.SpiderStat:
		return newModelDelegate(interfaces.ModelIdSpiderStat, doc, opts...)
	case *models.DataSource:
		return newModelDelegate(interfaces.ModelIdDataSource, doc, opts...)
	case *models.DataCollection:
		return newModelDelegate(interfaces.ModelIdDataCollection, doc, opts...)
	case *models.Result:
		return newModelDelegate(interfaces.ModelIdResult, doc, opts...)
	case *models.Password:
		return newModelDelegate(interfaces.ModelIdPassword, doc, opts...)
	case *models.ExtraValue:
		return newModelDelegate(interfaces.ModelIdExtraValue, doc, opts...)
	case *models.Git:
		return newModelDelegate(interfaces.ModelIdGit, doc, opts...)
	case *models.UserRole:
		return newModelDelegate(interfaces.ModelIdUserRole, doc, opts...)
	case *models.Permission:
		return newModelDelegate(interfaces.ModelIdPermission, doc, opts...)
	case *models.RolePermission:
		return newModelDelegate(interfaces.ModelIdRolePermission, doc, opts...)
	case *models.Environment:
		return newModelDelegate(interfaces.ModelIdEnvironment, doc, opts...)
	case *models.DependencySetting:
		return newModelDelegate(interfaces.ModelIdDependencySetting, doc, opts...)
	default:
		_ = trace.TraceError(errors.ErrorModelInvalidType)
		return nil
	}
}

func newModelDelegate(id interfaces.ModelId, doc interfaces.Model, opts ...ModelDelegateOption) interfaces.GrpcClientModelDelegate {
	var err error

	// collection name
	colName := models.GetModelColName(id)

	// model delegate
	d := &ModelDelegate{
		id:      id,
		colName: colName,
		doc:     doc,
		cfgPath: config2.GetConfigPath(),
		a: &models.Artifact{
			Col: colName,
		},
	}

	// config path
	if viper.GetString("config.path") != "" {
		d.cfgPath = viper.GetString("config.path")
	}

	// apply options
	for _, opt := range opts {
		opt(d)
	}

	// grpc client
	d.c, err = client.GetClient()
	if err != nil {
		trace.PrintError(errors.ErrorModelInvalidType)
		return nil
	}
	if !d.c.IsStarted() {
		if err := d.c.Start(); err != nil {
			trace.PrintError(err)
			return nil
		}
	} else if d.c.IsClosed() {
		if err := d.c.Restart(); err != nil {
			trace.PrintError(err)
			return nil
		}
	}

	return d
}

type ModelDelegate struct {
	// settings
	cfgPath string

	// internals
	id      interfaces.ModelId
	colName string
	c       interfaces.GrpcClient
	doc     interfaces.Model
	a       interfaces.ModelArtifact
}

func (d *ModelDelegate) Add() (err error) {
	return d.do(interfaces.ModelDelegateMethodAdd)
}

func (d *ModelDelegate) Save() (err error) {
	return d.do(interfaces.ModelDelegateMethodSave)
}

func (d *ModelDelegate) Delete() (err error) {
	return d.do(interfaces.ModelDelegateMethodDelete)
}

func (d *ModelDelegate) GetArtifact() (res interfaces.ModelArtifact, err error) {
	if err := d.do(interfaces.ModelDelegateMethodGetArtifact); err != nil {
		return nil, err
	}
	return d.a, nil
}

func (d *ModelDelegate) GetModel() (res interfaces.Model) {
	return d.doc
}

func (d *ModelDelegate) Refresh() (err error) {
	return d.refresh()
}

func (d *ModelDelegate) GetConfigPath() (path string) {
	return d.cfgPath
}

func (d *ModelDelegate) SetConfigPath(path string) {
	d.cfgPath = path
}

func (d *ModelDelegate) Close() (err error) {
	return d.c.Stop()
}

func (d *ModelDelegate) ToBytes(m interface{}) (bytes []byte, err error) {
	if m != nil {
		return utils.JsonToBytes(m)
	}
	return json.Marshal(d.doc)
}

func (d *ModelDelegate) do(method interfaces.ModelDelegateMethod) (err error) {
	switch method {
	case interfaces.ModelDelegateMethodAdd:
		err = d.add()
	case interfaces.ModelDelegateMethodSave:
		err = d.save()
	case interfaces.ModelDelegateMethodDelete:
		err = d.delete()
	case interfaces.ModelDelegateMethodGetArtifact, interfaces.ModelDelegateMethodRefresh:
		return d.refresh()
	default:
		return trace.TraceError(errors.ErrorModelInvalidType)
	}

	if err != nil {
		return err
	}

	return nil
}

func (d *ModelDelegate) add() (err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	method := interfaces.ModelDelegateMethod(interfaces.ModelDelegateMethodAdd)
	res, err := d.c.GetModelDelegateClient().Do(ctx, d.c.NewRequest(entity.GrpcDelegateMessage{
		ModelId: d.id,
		Method:  method,
		Data:    d.mustGetData(),
	}))
	if err != nil {
		return trace.TraceError(err)
	}
	if err := d.deserialize(res, method); err != nil {
		return err
	}
	return d.refreshArtifact()
}

func (d *ModelDelegate) save() (err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	method := interfaces.ModelDelegateMethod(interfaces.ModelDelegateMethodSave)
	res, err := d.c.GetModelDelegateClient().Do(ctx, d.c.NewRequest(entity.GrpcDelegateMessage{
		ModelId: d.id,
		Method:  method,
		Data:    d.mustGetData(),
	}))
	if err != nil {
		return trace.TraceError(err)
	}
	if err := d.deserialize(res, method); err != nil {
		return err
	}
	return d.refreshArtifact()
}

func (d *ModelDelegate) delete() (err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	method := interfaces.ModelDelegateMethod(interfaces.ModelDelegateMethodDelete)
	res, err := d.c.GetModelDelegateClient().Do(ctx, d.c.NewRequest(entity.GrpcDelegateMessage{
		ModelId: d.id,
		Method:  method,
		Data:    d.mustGetData(),
	}))
	if err != nil {
		return trace.TraceError(err)
	}
	if err := d.deserialize(res, method); err != nil {
		return err
	}
	return d.refreshArtifact()
}

func (d *ModelDelegate) refresh() (err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	method := interfaces.ModelDelegateMethod(interfaces.ModelDelegateMethodRefresh)
	res, err := d.c.GetModelDelegateClient().Do(ctx, d.c.NewRequest(entity.GrpcDelegateMessage{
		ModelId: d.id,
		Method:  method,
		Data:    d.mustGetData(),
	}))
	if err != nil {
		return trace.TraceError(err)
	}
	if err := d.deserialize(res, method); err != nil {
		return err
	}
	return nil
}

func (d *ModelDelegate) refreshArtifact() (err error) {
	_, err = d.getArtifact()
	return err
}

func (d *ModelDelegate) getArtifact() (res2 interfaces.ModelArtifact, err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	method := interfaces.ModelDelegateMethod(interfaces.ModelDelegateMethodGetArtifact)
	res, err := d.c.GetModelDelegateClient().Do(ctx, d.c.NewRequest(entity.GrpcDelegateMessage{
		ModelId: d.id,
		Method:  method,
		Data:    d.mustGetData(),
	}))
	if err != nil {
		return nil, err
	}
	if err := d.deserialize(res, method); err != nil {
		return nil, err
	}
	return d.a, nil
}

func (d *ModelDelegate) mustGetData() (data []byte) {
	data, err := d.getData()
	if err != nil {
		panic(err)
	}
	return data
}

func (d *ModelDelegate) getData() (data []byte, err error) {
	return json.Marshal(d.doc)
}

func (d *ModelDelegate) deserialize(res *grpc.Response, method interfaces.ModelDelegateMethod) (err error) {
	if method == interfaces.ModelDelegateMethodGetArtifact {
		res, err := NewBasicBinder(interfaces.ModelIdArtifact, res).Bind()
		if err != nil {
			return err
		}
		a, ok := res.(interfaces.ModelArtifact)
		if !ok {
			return trace.TraceError(errors.ErrorModelInvalidType)
		}
		d.a = a
	} else {
		d.doc, err = NewBasicBinder(d.id, res).Bind()
		if err != nil {
			return err
		}
	}
	return nil
}

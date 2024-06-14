package server

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/grpc"
)

func NewModelDelegateBinder(req *grpc.Request) (b *ModelDelegateBinder) {
	return &ModelDelegateBinder{
		req: req,
		msg: &entity.GrpcDelegateMessage{},
	}
}

type ModelDelegateBinder struct {
	req *grpc.Request
	msg interfaces.GrpcModelDelegateMessage
}

func (b *ModelDelegateBinder) Bind() (res interface{}, err error) {
	if err := b.bindDelegateMessage(); err != nil {
		return nil, err
	}

	m := models.NewModelMap()

	switch b.msg.GetModelId() {
	case interfaces.ModelIdArtifact:
		return b.process(&m.Artifact, interfaces.ModelIdTag)
	case interfaces.ModelIdTag:
		return b.process(&m.Tag, interfaces.ModelIdTag)
	case interfaces.ModelIdNode:
		return b.process(&m.Node, interfaces.ModelIdTag)
	case interfaces.ModelIdProject:
		return b.process(&m.Project, interfaces.ModelIdTag)
	case interfaces.ModelIdSpider:
		return b.process(&m.Spider, interfaces.ModelIdTag)
	case interfaces.ModelIdTask:
		return b.process(&m.Task)
	case interfaces.ModelIdJob:
		return b.process(&m.Job)
	case interfaces.ModelIdSchedule:
		return b.process(&m.Schedule)
	case interfaces.ModelIdUser:
		return b.process(&m.User)
	case interfaces.ModelIdSetting:
		return b.process(&m.Setting)
	case interfaces.ModelIdToken:
		return b.process(&m.Token)
	case interfaces.ModelIdVariable:
		return b.process(&m.Variable)
	case interfaces.ModelIdTaskQueue:
		return b.process(&m.TaskQueueItem)
	case interfaces.ModelIdTaskStat:
		return b.process(&m.TaskStat)
	case interfaces.ModelIdSpiderStat:
		return b.process(&m.SpiderStat)
	case interfaces.ModelIdDataSource:
		return b.process(&m.DataSource)
	case interfaces.ModelIdDataCollection:
		return b.process(&m.DataCollection)
	case interfaces.ModelIdResult:
		return b.process(&m.Result)
	case interfaces.ModelIdPassword:
		return b.process(&m.Password)
	case interfaces.ModelIdExtraValue:
		return b.process(&m.ExtraValue)
	case interfaces.ModelIdGit:
		return b.process(&m.Git)
	case interfaces.ModelIdRole:
		return b.process(&m.Role)
	case interfaces.ModelIdUserRole:
		return b.process(&m.UserRole)
	case interfaces.ModelIdPermission:
		return b.process(&m.Permission)
	case interfaces.ModelIdRolePermission:
		return b.process(&m.RolePermission)
	case interfaces.ModelIdEnvironment:
		return b.process(&m.Environment)
	case interfaces.ModelIdDependencySetting:
		return b.process(&m.DependencySetting)
	default:
		return nil, errors.ErrorModelInvalidModelId
	}
}

func (b *ModelDelegateBinder) MustBind() (res interface{}) {
	res, err := b.Bind()
	if err != nil {
		panic(err)
	}
	return res
}

func (b *ModelDelegateBinder) BindWithDelegateMessage() (res interface{}, msg interfaces.GrpcModelDelegateMessage, err error) {
	if err := json.Unmarshal(b.req.Data, b.msg); err != nil {
		return nil, nil, err
	}
	res, err = b.Bind()
	if err != nil {
		return nil, nil, err
	}
	return res, b.msg, nil
}

func (b *ModelDelegateBinder) process(d interface{}, fieldIds ...interfaces.ModelId) (res interface{}, err error) {
	if err := json.Unmarshal(b.msg.GetData(), d); err != nil {
		return nil, err
	}
	//return models.AssignFields(d, fieldIds...) // TODO: do we need to assign fields?
	return d, nil
}

func (b *ModelDelegateBinder) bindDelegateMessage() (err error) {
	return json.Unmarshal(b.req.Data, b.msg)
}

package client

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	grpc "github.com/crawlab-team/crawlab/grpc"
	"github.com/crawlab-team/crawlab/trace"
)

func NewListBinder(id interfaces.ModelId, res *grpc.Response) (b interfaces.GrpcModelListBinder) {
	return &ListBinder{
		id:  id,
		res: res,
	}
}

type ListBinder struct {
	id  interfaces.ModelId
	res *grpc.Response
}

func (b *ListBinder) Bind() (l interfaces.List, err error) {
	m := models.NewModelListMap()

	switch b.id {
	case interfaces.ModelIdArtifact:
		return b.Process(&m.Artifacts)
	case interfaces.ModelIdTag:
		return b.Process(&m.Tags)
	case interfaces.ModelIdNode:
		return b.Process(&m.Nodes)
	case interfaces.ModelIdProject:
		return b.Process(&m.Projects)
	case interfaces.ModelIdSpider:
		return b.Process(&m.Spiders)
	case interfaces.ModelIdTask:
		return b.Process(&m.Tasks)
	case interfaces.ModelIdJob:
		return b.Process(&m.Jobs)
	case interfaces.ModelIdSchedule:
		return b.Process(&m.Schedules)
	case interfaces.ModelIdUser:
		return b.Process(&m.Users)
	case interfaces.ModelIdSetting:
		return b.Process(&m.Settings)
	case interfaces.ModelIdToken:
		return b.Process(&m.Tokens)
	case interfaces.ModelIdVariable:
		return b.Process(&m.Variables)
	case interfaces.ModelIdTaskQueue:
		return b.Process(&m.TaskQueueItems)
	case interfaces.ModelIdTaskStat:
		return b.Process(&m.TaskStats)
	case interfaces.ModelIdSpiderStat:
		return b.Process(&m.SpiderStats)
	case interfaces.ModelIdDataSource:
		return b.Process(&m.DataSources)
	case interfaces.ModelIdDataCollection:
		return b.Process(&m.DataCollections)
	case interfaces.ModelIdResult:
		return b.Process(&m.Results)
	case interfaces.ModelIdPassword:
		return b.Process(&m.Passwords)
	case interfaces.ModelIdExtraValue:
		return b.Process(&m.ExtraValues)
	case interfaces.ModelIdGit:
		return b.Process(&m.Gits)
	case interfaces.ModelIdRole:
		return b.Process(&m.Roles)
	case interfaces.ModelIdUserRole:
		return b.Process(&m.UserRoles)
	case interfaces.ModelIdPermission:
		return b.Process(&m.PermissionList)
	case interfaces.ModelIdRolePermission:
		return b.Process(&m.RolePermissionList)
	case interfaces.ModelIdEnvironment:
		return b.Process(&m.Environments)
	case interfaces.ModelIdDependencySetting:
		return b.Process(&m.DependencySettings)
	default:
		return l, errors.ErrorModelInvalidModelId
	}
}

func (b *ListBinder) MustBind() (res interface{}) {
	res, err := b.Bind()
	if err != nil {
		panic(err)
	}
	return res
}

func (b *ListBinder) Process(d interface{}) (l interfaces.List, err error) {
	if err := json.Unmarshal(b.res.Data, d); err != nil {
		return l, trace.TraceError(err)
	}
	return d.(interfaces.List), nil
}

package controllers

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/gin-gonic/gin"
)

func NewJsonBinder(id ControllerId) (b *JsonBinder) {
	return &JsonBinder{
		id: id,
	}
}

type JsonBinder struct {
	id ControllerId
}

func (b *JsonBinder) Bind(c *gin.Context) (res interfaces.Model, err error) {
	// declare
	m := models.NewModelMap()

	switch b.id {
	case ControllerIdNode:
		err = c.ShouldBindJSON(&m.Node)
		return &m.Node, err
	case ControllerIdProject:
		err = c.ShouldBindJSON(&m.Project)
		return &m.Project, err
	case ControllerIdSpider:
		err = c.ShouldBindJSON(&m.Spider)
		return &m.Spider, err
	case ControllerIdTask:
		err = c.ShouldBindJSON(&m.Task)
		return &m.Task, err
	case ControllerIdJob:
		err = c.ShouldBindJSON(&m.Job)
		return &m.Job, err
	case ControllerIdSchedule:
		err = c.ShouldBindJSON(&m.Schedule)
		return &m.Schedule, err
	case ControllerIdUser:
		err = c.ShouldBindJSON(&m.User)
		return &m.User, nil
	case ControllerIdSetting:
		err = c.ShouldBindJSON(&m.Setting)
		return &m.Setting, nil
	case ControllerIdToken:
		err = c.ShouldBindJSON(&m.Token)
		return &m.Token, nil
	case ControllerIdVariable:
		err = c.ShouldBindJSON(&m.Variable)
		return &m.Variable, nil
	case ControllerIdTag:
		err = c.ShouldBindJSON(&m.Tag)
		return &m.Tag, nil
	case ControllerIdDataSource:
		err = c.ShouldBindJSON(&m.DataSource)
		return &m.DataSource, nil
	case ControllerIdDataCollection:
		err = c.ShouldBindJSON(&m.DataCollection)
		return &m.DataCollection, nil
	case ControllerIdGit:
		err = c.ShouldBindJSON(&m.Git)
		return &m.Git, nil
	case ControllerIdRole:
		err = c.ShouldBindJSON(&m.Role)
		return &m.Role, nil
	case ControllerIdPermission:
		err = c.ShouldBindJSON(&m.Permission)
		return &m.Permission, nil
	case ControllerIdEnvironment:
		err = c.ShouldBindJSON(&m.Environment)
		return &m.Environment, nil
	default:
		return nil, errors.ErrorControllerInvalidControllerId
	}
}

func (b *JsonBinder) BindList(c *gin.Context) (res interface{}, err error) {
	// declare
	m := models.NewModelListMap()

	// bind
	switch b.id {
	case ControllerIdNode:
		err = c.ShouldBindJSON(&m.Nodes)
		return m.Nodes, err
	case ControllerIdProject:
		err = c.ShouldBindJSON(&m.Projects)
		return m.Projects, err
	case ControllerIdSpider:
		err = c.ShouldBindJSON(&m.Spiders)
		return m.Spiders, err
	case ControllerIdTask:
		err = c.ShouldBindJSON(&m.Tasks)
		return m.Tasks, err
	case ControllerIdJob:
		err = c.ShouldBindJSON(&m.Jobs)
		return m.Jobs, err
	case ControllerIdSchedule:
		err = c.ShouldBindJSON(&m.Schedules)
		return m.Schedules, err
	case ControllerIdUser:
		err = c.ShouldBindJSON(&m.Users)
		return m.Users, nil
	case ControllerIdSetting:
		err = c.ShouldBindJSON(&m.Settings)
		return m.Settings, nil
	case ControllerIdToken:
		err = c.ShouldBindJSON(&m.Tokens)
		return m.Tokens, nil
	case ControllerIdVariable:
		err = c.ShouldBindJSON(&m.Variables)
		return m.Variables, nil
	case ControllerIdTag:
		err = c.ShouldBindJSON(&m.Tags)
		return m.Tags, nil
	case ControllerIdDataSource:
		err = c.ShouldBindJSON(&m.DataSources)
		return m.DataSources, nil
	case ControllerIdDataCollection:
		err = c.ShouldBindJSON(&m.DataCollections)
		return m.DataCollections, nil
	case ControllerIdGit:
		err = c.ShouldBindJSON(&m.Gits)
		return m.Gits, nil
	case ControllerIdRole:
		err = c.ShouldBindJSON(&m.Roles)
		return m.Roles, nil
	case ControllerIdEnvironment:
		err = c.ShouldBindJSON(&m.Environments)
		return m.Environments, nil
	default:
		return nil, errors.ErrorControllerInvalidControllerId
	}
}

func (b *JsonBinder) BindBatchRequestPayload(c *gin.Context) (payload entity.BatchRequestPayload, err error) {
	if err := c.ShouldBindJSON(&payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func (b *JsonBinder) BindBatchRequestPayloadWithStringData(c *gin.Context) (payload entity.BatchRequestPayloadWithStringData, res interfaces.Model, err error) {
	// declare
	m := models.NewModelMap()

	// bind
	if err := c.ShouldBindJSON(&payload); err != nil {
		return payload, nil, err
	}

	// validate
	if len(payload.Ids) == 0 ||
		len(payload.Fields) == 0 {
		return payload, nil, errors.ErrorControllerRequestPayloadInvalid
	}

	// unmarshall
	switch b.id {
	case ControllerIdNode:
		err = json.Unmarshal([]byte(payload.Data), &m.Node)
		return payload, &m.Node, err
	case ControllerIdProject:
		err = json.Unmarshal([]byte(payload.Data), &m.Project)
		return payload, &m.Project, err
	case ControllerIdSpider:
		err = json.Unmarshal([]byte(payload.Data), &m.Spider)
		return payload, &m.Spider, err
	case ControllerIdTask:
		err = json.Unmarshal([]byte(payload.Data), &m.Task)
		return payload, &m.Task, err
	case ControllerIdJob:
		err = json.Unmarshal([]byte(payload.Data), &m.Job)
		return payload, &m.Job, err
	case ControllerIdSchedule:
		err = json.Unmarshal([]byte(payload.Data), &m.Schedule)
		return payload, &m.Schedule, err
	case ControllerIdUser:
		err = json.Unmarshal([]byte(payload.Data), &m.User)
		return payload, &m.User, err
	case ControllerIdSetting:
		err = json.Unmarshal([]byte(payload.Data), &m.Setting)
		return payload, &m.Setting, err
	case ControllerIdToken:
		err = json.Unmarshal([]byte(payload.Data), &m.Token)
		return payload, &m.Token, err
	case ControllerIdVariable:
		err = json.Unmarshal([]byte(payload.Data), &m.Variable)
		return payload, &m.Variable, err
	case ControllerIdDataSource:
		err = json.Unmarshal([]byte(payload.Data), &m.DataSource)
		return payload, &m.DataSource, err
	case ControllerIdDataCollection:
		err = json.Unmarshal([]byte(payload.Data), &m.DataCollection)
		return payload, &m.DataCollection, err
	case ControllerIdEnvironment:
		err = json.Unmarshal([]byte(payload.Data), &m.Environment)
		return payload, &m.Environment, err
	default:
		return payload, nil, errors.ErrorControllerInvalidControllerId
	}
}

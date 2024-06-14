package binders

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
)

func NewColNameBinder(id interfaces.ModelId) (b *ColNameBinder) {
	return &ColNameBinder{id: id}
}

type ColNameBinder struct {
	id interfaces.ModelId
}

func (b *ColNameBinder) Bind() (res interface{}, err error) {
	switch b.id {
	// system models
	case interfaces.ModelIdArtifact:
		return interfaces.ModelColNameArtifact, nil
	case interfaces.ModelIdTag:
		return interfaces.ModelColNameTag, nil

	// operation models
	case interfaces.ModelIdNode:
		return interfaces.ModelColNameNode, nil
	case interfaces.ModelIdProject:
		return interfaces.ModelColNameProject, nil
	case interfaces.ModelIdSpider:
		return interfaces.ModelColNameSpider, nil
	case interfaces.ModelIdTask:
		return interfaces.ModelColNameTask, nil
	case interfaces.ModelIdJob:
		return interfaces.ModelColNameJob, nil
	case interfaces.ModelIdSchedule:
		return interfaces.ModelColNameSchedule, nil
	case interfaces.ModelIdUser:
		return interfaces.ModelColNameUser, nil
	case interfaces.ModelIdSetting:
		return interfaces.ModelColNameSetting, nil
	case interfaces.ModelIdToken:
		return interfaces.ModelColNameToken, nil
	case interfaces.ModelIdVariable:
		return interfaces.ModelColNameVariable, nil
	case interfaces.ModelIdTaskQueue:
		return interfaces.ModelColNameTaskQueue, nil
	case interfaces.ModelIdTaskStat:
		return interfaces.ModelColNameTaskStat, nil
	case interfaces.ModelIdSpiderStat:
		return interfaces.ModelColNameSpiderStat, nil
	case interfaces.ModelIdDataSource:
		return interfaces.ModelColNameDataSource, nil
	case interfaces.ModelIdDataCollection:
		return interfaces.ModelColNameDataCollection, nil
	case interfaces.ModelIdPassword:
		return interfaces.ModelColNamePasswords, nil
	case interfaces.ModelIdExtraValue:
		return interfaces.ModelColNameExtraValues, nil
	case interfaces.ModelIdGit:
		return interfaces.ModelColNameGit, nil
	case interfaces.ModelIdRole:
		return interfaces.ModelColNameRole, nil
	case interfaces.ModelIdUserRole:
		return interfaces.ModelColNameUserRole, nil
	case interfaces.ModelIdPermission:
		return interfaces.ModelColNamePermission, nil
	case interfaces.ModelIdRolePermission:
		return interfaces.ModelColNameRolePermission, nil
	case interfaces.ModelIdEnvironment:
		return interfaces.ModelColNameEnvironment, nil
	case interfaces.ModelIdDependencySetting:
		return interfaces.ModelColNameDependencySetting, nil

	// invalid
	default:
		return res, errors.ErrorModelNotImplemented
	}
}

func (b *ColNameBinder) MustBind() (res interface{}) {
	res, err := b.Bind()
	if err != nil {
		panic(err)
	}
	return res
}

func (b *ColNameBinder) BindString() (res string, err error) {
	res_, err := b.Bind()
	if err != nil {
		return "", err
	}
	res = res_.(string)
	return res, nil
}

func (b *ColNameBinder) MustBindString() (res string) {
	return b.MustBind().(string)
}

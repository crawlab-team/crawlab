package test

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/client"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/node/test"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestModelDelegate_Do(t *testing.T) {
	var err error

	T, _ = NewTest()
	T.Setup(t)

	// add
	modelDelegateAdd(t)
	project, err := test.T.ModelSvc.GetProject(bson.M{"name": "test-project"}, nil)
	require.Nil(t, err)

	// get artifact
	a := modelDelegateGetArtifact(t)
	require.Equal(t, project.GetId(), a.GetId())

	// save
	modelDelegateSave(t)
	project, err = test.T.ModelSvc.GetProject(bson.M{"name": "test-new-project"}, nil)
	require.Nil(t, err)
	require.Equal(t, "test-new-project", project.Name)
	require.Equal(t, "test-new-description", project.Description)

	// delete
	modelDelegateDelete(t)
	_, err = test.T.ModelSvc.GetProject(bson.M{"name": "test-new-project"}, nil)
	require.Equal(t, mongo2.ErrNoDocuments, err)
}

func TestModelDelegate_Do_All(t *testing.T) {
	T, _ = NewTest()
	T.Setup(t)

	// add
	modelDelegateAddAll(t)
	modelDelegateValidateAddAll(t)
}

func modelDelegateAdd(t *testing.T) {
	// modelDelegateAdd
	project := &models.Project{
		Name:        "test-project",
		Description: "test-description",
	}
	projectD := client.NewModelDelegate(project, client.WithDelegateConfigPath(T.Client.GetConfigPath()))
	err := projectD.Add()
	require.Nil(t, err)
}

func modelDelegateGetArtifact(t *testing.T) interfaces.ModelArtifact {
	project, err := test.T.ModelSvc.GetProject(bson.M{"name": "test-project"}, nil)
	require.Nil(t, err)

	projectD := client.NewModelDelegate(project, client.WithDelegateConfigPath(T.Client.GetConfigPath()))
	a, err := projectD.GetArtifact()
	require.Nil(t, err)
	return a
}

func modelDelegateSave(t *testing.T) {
	project, err := test.T.ModelSvc.GetProject(bson.M{"name": "test-project"}, nil)
	require.Nil(t, err)

	project.Name = "test-new-project"
	project.Description = "test-new-description"

	projectD := client.NewModelDelegate(project, client.WithDelegateConfigPath(T.Client.GetConfigPath()))
	err = projectD.Save()
	require.Nil(t, err)
}

func modelDelegateDelete(t *testing.T) {
	project, err := test.T.ModelSvc.GetProject(bson.M{"name": "test-new-project"}, nil)
	require.Nil(t, err)

	projectD := client.NewModelDelegate(project, client.WithDelegateConfigPath(T.Client.GetConfigPath()))
	err = projectD.Delete()
	require.Nil(t, err)
}

func modelDelegateAddAll(t *testing.T) {
	var err error
	cfgOpt := client.WithDelegateConfigPath(T.Client.GetConfigPath())
	m := models.NewModelMap()
	err = client.NewModelDelegate(&m.Tag, cfgOpt).Add()
	require.Nil(t, err)
	err = client.NewModelDelegate(&m.Node, cfgOpt).Add()
	require.Nil(t, err)
	err = client.NewModelDelegate(&m.Project, cfgOpt).Add()
	require.Nil(t, err)
	err = client.NewModelDelegate(&m.Spider, cfgOpt).Add()
	require.Nil(t, err)
	err = client.NewModelDelegate(&m.Task, cfgOpt).Add()
	require.Nil(t, err)
	err = client.NewModelDelegate(&m.Job, cfgOpt).Add()
	require.Nil(t, err)
	err = client.NewModelDelegate(&m.Schedule, cfgOpt).Add()
	require.Nil(t, err)
	err = client.NewModelDelegate(&m.User, cfgOpt).Add()
	require.Nil(t, err)
	err = client.NewModelDelegate(&m.Setting, cfgOpt).Add()
	require.Nil(t, err)
	err = client.NewModelDelegate(&m.Token, cfgOpt).Add()
	require.Nil(t, err)
	err = client.NewModelDelegate(&m.Variable, cfgOpt).Add()
	require.Nil(t, err)
	err = client.NewModelDelegate(&m.User, cfgOpt).Add()
	require.Nil(t, err)
}

func modelDelegateValidateAddAll(t *testing.T) {
	var err error
	_, err = test.T.ModelSvc.GetBaseService(interfaces.ModelIdNode).Get(bson.M{}, nil)
	require.Nil(t, err)
	_, err = test.T.ModelSvc.GetBaseService(interfaces.ModelIdNode).Get(bson.M{}, nil)
	require.Nil(t, err)
	_, err = test.T.ModelSvc.GetBaseService(interfaces.ModelIdNode).Get(bson.M{}, nil)
	require.Nil(t, err)
	_, err = test.T.ModelSvc.GetBaseService(interfaces.ModelIdNode).Get(bson.M{}, nil)
	require.Nil(t, err)
	_, err = test.T.ModelSvc.GetBaseService(interfaces.ModelIdNode).Get(bson.M{}, nil)
	require.Nil(t, err)
	_, err = test.T.ModelSvc.GetBaseService(interfaces.ModelIdNode).Get(bson.M{}, nil)
	require.Nil(t, err)
	_, err = test.T.ModelSvc.GetBaseService(interfaces.ModelIdNode).Get(bson.M{}, nil)
	require.Nil(t, err)
	_, err = test.T.ModelSvc.GetBaseService(interfaces.ModelIdNode).Get(bson.M{}, nil)
	require.Nil(t, err)
	_, err = test.T.ModelSvc.GetBaseService(interfaces.ModelIdNode).Get(bson.M{}, nil)
	require.Nil(t, err)
	_, err = test.T.ModelSvc.GetBaseService(interfaces.ModelIdNode).Get(bson.M{}, nil)
	require.Nil(t, err)
	_, err = test.T.ModelSvc.GetBaseService(interfaces.ModelIdNode).Get(bson.M{}, nil)
	require.Nil(t, err)
}

package test

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/node/test"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestModelBaseService_GetById(t *testing.T) {
	var err error

	T, _ = NewTest()
	T.Setup(t)

	// add
	modelDelegateAdd(t)
	p, err := test.T.ModelSvc.GetProject(bson.M{"name": "test-project"}, nil)
	require.Nil(t, err)

	// get by id
	ctx, cancel := T.Client.Context()
	defer cancel()
	req, err := T.Client.NewModelBaseServiceRequest(interfaces.ModelIdProject, &entity.GrpcBaseServiceParams{Id: p.Id})
	require.Nil(t, err)
	res, err := T.Client.GetModelBaseServiceClient().GetById(ctx, req)
	require.Nil(t, err)
	var p2 models.Project
	err = json.Unmarshal(res.Data, &p2)
	require.Nil(t, err)

	// validate
	require.Equal(t, p.Id, p2.Id)
	require.Equal(t, p.Name, p2.Name)
	require.Equal(t, p.Description, p2.Description)
}

func TestModelBaseService_Get(t *testing.T) {
	var err error

	T, _ = NewTest()
	T.Setup(t)

	// add
	modelDelegateAdd(t)
	p, err := test.T.ModelSvc.GetProject(bson.M{"name": "test-project"}, nil)
	require.Nil(t, err)

	// get
	ctx, cancel := T.Client.Context()
	defer cancel()
	req, err := T.Client.NewModelBaseServiceRequest(interfaces.ModelIdProject, &entity.GrpcBaseServiceParams{Query: bson.M{"name": "test-project"}})
	require.Nil(t, err)
	res, err := T.Client.GetModelBaseServiceClient().Get(ctx, req)
	require.Nil(t, err)
	var p2 models.Project
	err = json.Unmarshal(res.Data, &p2)
	require.Nil(t, err)

	// validate
	require.Equal(t, p.Id, p2.Id)
	require.Equal(t, p.Name, p2.Name)
	require.Equal(t, p.Description, p2.Description)
}

func TestModelBaseService_GetList(t *testing.T) {
	var err error

	T, _ = NewTest()
	T.Setup(t)

	// add
	n := 10
	for i := 0; i < n; i++ {
		modelDelegateAdd(t)
	}

	// get list
	ctx, cancel := T.Client.Context()
	defer cancel()
	req, err := T.Client.NewModelBaseServiceRequest(interfaces.ModelIdProject, &entity.GrpcBaseServiceParams{Query: bson.M{"name": "test-project"}})
	require.Nil(t, err)
	res, err := T.Client.GetModelBaseServiceClient().GetList(ctx, req)
	require.Nil(t, err)
	var list []models.Project
	err = json.Unmarshal(res.Data, &list)
	require.Nil(t, err)

	// validate
	require.Equal(t, n, len(list))
}

func TestModelBaseService_DeleteById(t *testing.T) {
	var err error

	T, _ = NewTest()
	T.Setup(t)

	// add
	modelDelegateAdd(t)
	p, err := test.T.ModelSvc.GetProject(bson.M{"name": "test-project"}, nil)
	require.Nil(t, err)

	// delete by id
	ctx, cancel := T.Client.Context()
	defer cancel()
	req, err := T.Client.NewModelBaseServiceRequest(interfaces.ModelIdProject, &entity.GrpcBaseServiceParams{Id: p.Id})
	require.Nil(t, err)
	_, err = T.Client.GetModelBaseServiceClient().DeleteById(ctx, req)
	require.Nil(t, err)

	// validate
	p, err = test.T.ModelSvc.GetProjectById(p.Id)
	require.Equal(t, mongo2.ErrNoDocuments, err)
}

func TestModelBaseService_Delete(t *testing.T) {
	var err error

	T, _ = NewTest()
	T.Setup(t)

	// add
	modelDelegateAdd(t)
	p, err := test.T.ModelSvc.GetProject(bson.M{"name": "test-project"}, nil)
	require.Nil(t, err)

	// delete by id
	ctx, cancel := T.Client.Context()
	defer cancel()
	req, err := T.Client.NewModelBaseServiceRequest(interfaces.ModelIdProject, &entity.GrpcBaseServiceParams{Query: bson.M{"name": "test-project"}})
	require.Nil(t, err)
	_, err = T.Client.GetModelBaseServiceClient().Delete(ctx, req)
	require.Nil(t, err)

	// validate
	p, err = test.T.ModelSvc.GetProjectById(p.Id)
	require.Equal(t, mongo2.ErrNoDocuments, err)
}

func TestModelBaseService_DeleteList(t *testing.T) {
	var err error

	T, _ = NewTest()
	T.Setup(t)

	// add
	n := 10
	for i := 0; i < n; i++ {
		modelDelegateAdd(t)
	}

	// delete by id
	ctx, cancel := T.Client.Context()
	defer cancel()
	req, err := T.Client.NewModelBaseServiceRequest(interfaces.ModelIdProject, &entity.GrpcBaseServiceParams{Query: bson.M{"name": "test-project"}})
	require.Nil(t, err)
	_, err = T.Client.GetModelBaseServiceClient().DeleteList(ctx, req)
	require.Nil(t, err)

	// validate
	list, err := test.T.ModelSvc.GetProjectList(bson.M{"name": "test-project"}, nil)
	require.Nil(t, err)
	require.Equal(t, 0, len(list))
}

func TestModelBaseService_ForceDeleteList(t *testing.T) {
	var err error

	T, _ = NewTest()
	T.Setup(t)

	// add
	n := 10
	for i := 0; i < n; i++ {
		modelDelegateAdd(t)
	}

	// delete by id
	ctx, cancel := T.Client.Context()
	defer cancel()
	req, err := T.Client.NewModelBaseServiceRequest(interfaces.ModelIdProject, &entity.GrpcBaseServiceParams{Query: bson.M{"name": "test-project"}})
	require.Nil(t, err)
	_, err = T.Client.GetModelBaseServiceClient().ForceDeleteList(ctx, req)
	require.Nil(t, err)

	// validate
	list, err := test.T.ModelSvc.GetProjectList(bson.M{"name": "test-project"}, nil)
	require.Nil(t, err)
	require.Equal(t, 0, len(list))
}

func TestModelBaseService_UpdateById(t *testing.T) {
	var err error

	T, _ = NewTest()
	T.Setup(t)

	// add
	modelDelegateAdd(t)
	p, err := test.T.ModelSvc.GetProject(bson.M{"name": "test-project"}, nil)
	require.Nil(t, err)

	// update by id
	ctx, cancel := T.Client.Context()
	defer cancel()
	update := bson.M{
		"name": "test-new-project",
	}
	req, err := T.Client.NewModelBaseServiceRequest(interfaces.ModelIdProject, &entity.GrpcBaseServiceParams{Id: p.Id, Update: update})
	require.Nil(t, err)
	_, err = T.Client.GetModelBaseServiceClient().UpdateById(ctx, req)
	require.Nil(t, err)

	// validate
	p2, err := test.T.ModelSvc.GetProjectById(p.Id)
	require.Nil(t, err)
	require.Equal(t, "test-new-project", p2.Name)
}

func TestModelBaseService_Update(t *testing.T) {
	var err error

	T, _ = NewTest()
	T.Setup(t)

	// add
	n := 10
	for i := 0; i < n; i++ {
		modelDelegateAdd(t)
	}

	// update
	ctx, cancel := T.Client.Context()
	defer cancel()
	update := bson.M{
		"name": "test-new-project",
	}
	req, err := T.Client.NewModelBaseServiceRequest(interfaces.ModelIdProject, &entity.GrpcBaseServiceParams{Query: bson.M{"name": "test-project"}, Update: update})
	require.Nil(t, err)
	_, err = T.Client.GetModelBaseServiceClient().Update(ctx, req)
	require.Nil(t, err)

	// validate
	list, err := test.T.ModelSvc.GetProjectList(bson.M{"name": "test-project"}, nil)
	require.Nil(t, err)
	require.Equal(t, 0, len(list))
	list, err = test.T.ModelSvc.GetProjectList(bson.M{"name": "test-new-project"}, nil)
	require.Nil(t, err)
	require.Equal(t, n, len(list))
}

func TestModelBaseService_Insert(t *testing.T) {
	var err error

	T, _ = NewTest()
	T.Setup(t)

	// insert
	var docs []interface{}
	n := 10
	for i := 0; i < n; i++ {
		docs = append(docs, models.Project{
			Id:   primitive.NewObjectID(),
			Name: "test-project",
		})
	}
	ctx, cancel := T.Client.Context()
	defer cancel()
	req, err := T.Client.NewModelBaseServiceRequest(interfaces.ModelIdProject, &entity.GrpcBaseServiceParams{Docs: docs})
	require.Nil(t, err)
	_, err = T.Client.GetModelBaseServiceClient().Insert(ctx, req)
	require.Nil(t, err)

	// validate
	list, err := test.T.ModelSvc.GetProjectList(bson.M{"name": "test-project"}, nil)
	require.Nil(t, err)
	require.Equal(t, n, len(list))
}

func TestModelBaseService_Count(t *testing.T) {
	var err error

	T, _ = NewTest()
	T.Setup(t)

	// add
	n := 10
	for i := 0; i < n; i++ {
		modelDelegateAdd(t)
	}

	// count
	ctx, cancel := T.Client.Context()
	defer cancel()
	req, err := T.Client.NewModelBaseServiceRequest(interfaces.ModelIdProject, &entity.GrpcBaseServiceParams{Query: bson.M{"name": "test-project"}})
	require.Nil(t, err)
	res, err := T.Client.GetModelBaseServiceClient().Count(ctx, req)
	require.Nil(t, err)

	// validate
	var total int
	err = json.Unmarshal(res.Data, &total)
	require.Nil(t, err)
	require.Equal(t, n, total)
}

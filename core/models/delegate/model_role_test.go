package delegate_test

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	models2 "github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRole_Add(t *testing.T) {
	SetupTest(t)

	p := &models2.Role{}

	err := delegate.NewModelDelegate(p).Add()
	require.Nil(t, err)
	require.NotNil(t, p.Id)

	a, err := delegate.NewModelDelegate(p).GetArtifact()
	require.Nil(t, err)
	require.Equal(t, p.Id, a.GetId())
	require.NotNil(t, a.GetSys().GetCreateTs())
	require.NotNil(t, a.GetSys().GetUpdateTs())
}

func TestRole_Save(t *testing.T) {
	SetupTest(t)

	p := &models2.Role{}

	err := delegate.NewModelDelegate(p).Add()
	require.Nil(t, err)

	name := "test_role"
	p.Name = name
	err = delegate.NewModelDelegate(p).Save()
	require.Nil(t, err)

	err = mongo.GetMongoCol(interfaces.ModelColNameRole).FindId(p.Id).One(&p)
	require.Nil(t, err)
	require.Equal(t, name, p.Name)
}

func TestRole_Delete(t *testing.T) {
	SetupTest(t)

	p := &models2.Role{
		Name: "test_role",
	}

	err := delegate.NewModelDelegate(p).Add()
	require.Nil(t, err)

	err = delegate.NewModelDelegate(p).Delete()
	require.Nil(t, err)

	var a models2.Artifact
	col := mongo.GetMongoCol(interfaces.ModelColNameArtifact)
	err = col.FindId(p.Id).One(&a)
	require.Nil(t, err)
	require.NotNil(t, a.Obj)
	require.True(t, a.Del)
}

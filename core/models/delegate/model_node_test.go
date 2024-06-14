package delegate_test

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	models2 "github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNode_Add(t *testing.T) {
	SetupTest(t)

	n := &models2.Node{}

	err := delegate.NewModelDelegate(n).Add()
	require.Nil(t, err)
	require.NotNil(t, n.Id)

	// validate artifact
	a, err := delegate.NewModelDelegate(n).GetArtifact()
	require.Nil(t, err)
	require.Equal(t, n.Id, a.GetId())
	require.NotNil(t, a.GetSys().GetCreateTs())
	require.NotNil(t, a.GetSys().GetUpdateTs())
}

func TestNode_Save(t *testing.T) {
	SetupTest(t)

	n := &models2.Node{}

	err := delegate.NewModelDelegate(n).Add()

	name := "test_node"
	n.Name = name
	err = delegate.NewModelDelegate(n).Save()
	require.Nil(t, err)

	err = mongo.GetMongoCol(interfaces.ModelColNameNode).FindId(n.Id).One(&n)
	require.Nil(t, err)
	require.Equal(t, name, n.Name)
}

func TestNode_Delete(t *testing.T) {
	SetupTest(t)

	n := &models2.Node{
		Name: "test_node",
	}

	err := delegate.NewModelDelegate(n).Add()
	require.Nil(t, err)

	err = delegate.NewModelDelegate(n).Delete()
	require.Nil(t, err)

	var a models2.Artifact
	col := mongo.GetMongoCol(interfaces.ModelColNameArtifact)
	err = col.FindId(n.Id).One(&a)
	require.Nil(t, err)
	require.NotNil(t, a.Obj)
	require.True(t, a.Del)
}

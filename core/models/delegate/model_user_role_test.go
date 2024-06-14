package delegate_test

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/common"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	models2 "github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func init() {
	viper.Set("mongo.db", "crawlab_test")
	common.CreateIndexes()
}

func TestUserRole_Add(t *testing.T) {
	SetupTest(t)

	p := &models2.UserRole{}

	err := delegate.NewModelDelegate(p).Add()
	require.Nil(t, err)
	require.NotNil(t, p.Id)

	a, err := delegate.NewModelDelegate(p).GetArtifact()
	require.Nil(t, err)
	require.Equal(t, p.Id, a.GetId())
	require.NotNil(t, a.GetSys().GetCreateTs())
	require.NotNil(t, a.GetSys().GetUpdateTs())
}

func TestUserRole_Save(t *testing.T) {
	SetupTest(t)

	p := &models2.UserRole{
		UserId: primitive.NewObjectID(),
		RoleId: primitive.NewObjectID(),
	}

	err := delegate.NewModelDelegate(p).Add()
	require.Nil(t, err)

	uid := primitive.NewObjectID()
	rid := primitive.NewObjectID()
	p.UserId = uid
	p.RoleId = rid
	err = delegate.NewModelDelegate(p).Save()
	require.Nil(t, err)

	err = mongo.GetMongoCol(interfaces.ModelColNameUserRole).FindId(p.Id).One(&p)
	require.Nil(t, err)
	require.Equal(t, uid, p.UserId)
	require.Equal(t, rid, p.RoleId)
}

func TestUserRole_Delete(t *testing.T) {
	SetupTest(t)

	p := &models2.UserRole{}

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

func TestUserRole_AddDuplicates(t *testing.T) {
	SetupTest(t)

	uid := primitive.NewObjectID()
	rid := primitive.NewObjectID()
	p := &models2.UserRole{
		UserId: uid,
		RoleId: rid,
	}
	p2 := &models2.UserRole{
		UserId: uid,
		RoleId: rid,
	}

	err := delegate.NewModelDelegate(p).Add()
	require.Nil(t, err)
	err = delegate.NewModelDelegate(p2).Add()
	require.NotNil(t, err)
}

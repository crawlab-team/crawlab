package test

import (
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"testing"
)

func TestTaskController_Delete(t *testing.T) {
	T.Setup(t)
	e := T.NewExpect(t)

	task := models.Task{
		Id: primitive.NewObjectID(),
	}

	// add task
	err := delegate.NewModelDelegate(&task).Add()
	require.Nil(t, err)

	// add task stat
	err = delegate.NewModelDelegate(&models.TaskStat{
		Id: task.Id,
	}).Add()
	require.Nil(t, err)

	// delete
	T.WithAuth(e.DELETE("/tasks/" + task.Id.Hex())).
		Expect().Status(http.StatusOK)

	// get
	T.WithAuth(e.GET("/tasks/" + task.Id.Hex())).
		Expect().Status(http.StatusNotFound)

	// task stats
	modelTaskStatSvc := service.NewBaseService(interfaces.ModelIdTaskStat)
	taskStatCount, err := modelTaskStatSvc.Count(bson.M{
		"_id": task.Id,
	})
	require.Nil(t, err)
	require.Zero(t, taskStatCount)
}

func TestTaskController_DeleteList(t *testing.T) {
	T.Setup(t)
	e := T.NewExpect(t)

	tasks := []models.Task{
		{
			Id: primitive.NewObjectID(),
		},
		{
			Id: primitive.NewObjectID(),
		},
	}

	// add spiders
	var taskIds []primitive.ObjectID
	for _, task := range tasks {
		// add task
		err := delegate.NewModelDelegate(&task).Add()
		require.Nil(t, err)

		// add task stat
		err = delegate.NewModelDelegate(&models.TaskStat{
			Id: task.Id,
		}).Add()
		require.Nil(t, err)

		taskIds = append(taskIds, task.Id)
	}

	// delete tasks
	T.WithAuth(e.DELETE("/tasks")).
		WithJSON(entity.BatchRequestPayload{
			Ids: taskIds,
		}).Expect().Status(http.StatusOK)

	// get tasks
	for _, task := range tasks {
		// get
		T.WithAuth(e.GET("/tasks/" + task.Id.Hex())).
			Expect().Status(http.StatusNotFound)
	}

	// task stats
	modelTaskStatSvc := service.NewBaseService(interfaces.ModelIdTaskStat)
	taskStatCount, err := modelTaskStatSvc.Count(bson.M{
		"_id": bson.M{
			"$in": taskIds,
		},
	})
	require.Nil(t, err)
	require.Zero(t, taskStatCount)
}

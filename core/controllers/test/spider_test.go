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

func TestSpiderController_Delete(t *testing.T) {
	T.Setup(t)
	e := T.NewExpect(t)

	s := models.Spider{
		Name:        "test spider",
		Description: "this is a test spider",
		ColName:     "test col name",
	}

	// add spider
	res := T.WithAuth(e.POST("/spiders")).
		WithJSON(s).
		Expect().Status(http.StatusOK).
		JSON().Object()
	res.Path("$.data._id").NotNull()
	id := res.Path("$.data._id").String().Raw()
	oid, err := primitive.ObjectIDFromHex(id)
	require.Nil(t, err)
	require.False(t, oid.IsZero())

	// add tasks
	var taskIds []primitive.ObjectID
	tasks := []models.Task{
		{
			Id:       primitive.NewObjectID(),
			SpiderId: oid,
		},
		{
			Id:       primitive.NewObjectID(),
			SpiderId: oid,
		},
	}
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

	// delete
	T.WithAuth(e.DELETE("/spiders/" + id)).
		Expect().Status(http.StatusOK)

	// get
	T.WithAuth(e.GET("/spiders/" + id)).
		Expect().Status(http.StatusNotFound)

	// get tasks
	for _, task := range tasks {
		T.WithAuth(e.GET("/tasks/" + task.Id.Hex())).
			Expect().Status(http.StatusNotFound)
	}

	// spider stat
	modelSpiderStatSvc := service.NewBaseService(interfaces.ModelIdSpiderStat)
	spiderStatCount, err := modelSpiderStatSvc.Count(bson.M{
		"_id": oid,
	})
	require.Nil(t, err)
	require.Zero(t, spiderStatCount)

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

func TestSpiderController_DeleteList(t *testing.T) {
	T.Setup(t)
	e := T.NewExpect(t)

	spiders := []models.Spider{
		{
			Id:          primitive.NewObjectID(),
			Name:        "test spider 1",
			Description: "this is a test spider 1",
			ColName:     "test col name 1",
		},
		{
			Id:          primitive.NewObjectID(),
			Name:        "test spider 2",
			Description: "this is a test spider 2",
			ColName:     "test col name 2",
		},
	}

	// add spiders
	for _, spider := range spiders {
		T.WithAuth(e.POST("/spiders")).
			WithJSON(spider).
			Expect().Status(http.StatusOK)
	}

	var spiderIds []primitive.ObjectID
	var taskIds []primitive.ObjectID
	for _, spider := range spiders {
		// task id
		taskId := primitive.NewObjectID()

		// add task
		err := delegate.NewModelDelegate(&models.Task{
			Id:       taskId,
			SpiderId: spider.Id,
		}).Add()
		require.Nil(t, err)

		// add task stats
		err = delegate.NewModelDelegate(&models.TaskStat{
			Id: taskId,
		}).Add()
		require.Nil(t, err)

		spiderIds = append(spiderIds, spider.Id)
		taskIds = append(taskIds, taskId)
	}

	// delete spiders
	T.WithAuth(e.DELETE("/spiders")).
		WithJSON(entity.BatchRequestPayload{
			Ids: spiderIds,
		}).Expect().Status(http.StatusOK)

	// get spiders
	for _, spider := range spiders {
		// get
		T.WithAuth(e.GET("/spiders/" + spider.Id.Hex())).
			Expect().Status(http.StatusNotFound)
	}

	// get tasks
	for _, taskId := range taskIds {
		T.WithAuth(e.GET("/tasks/" + taskId.Hex())).
			Expect().Status(http.StatusNotFound)
	}

	// spider stat
	modelSpiderStatSvc := service.NewBaseService(interfaces.ModelIdSpiderStat)
	spiderStatCount, err := modelSpiderStatSvc.Count(bson.M{
		"_id": bson.M{
			"$in": spiderIds,
		},
	})
	require.Nil(t, err)
	require.Zero(t, spiderStatCount)

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

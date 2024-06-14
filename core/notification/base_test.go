package notification

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/gavv/httpexpect/v2"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http/httptest"
	"testing"
)

func init() {
	viper.Set("mongo.db", "crawlab_test")
	var err error
	T, err = NewTest()
	if err != nil {
		panic(err)
	}
}

type Test struct {
	svc *Service
	svr *httptest.Server

	// test data
	TestNode     interfaces.Node
	TestSpider   interfaces.Spider
	TestTask     interfaces.Task
	TestTaskStat interfaces.TaskStat
}

func (t *Test) Setup(t2 *testing.T) {
	_ = t.svc.Start()
	t2.Cleanup(t.Cleanup)
}

func (t *Test) Cleanup() {
	_ = t.svc.Stop()
}

func (t *Test) NewExpect(t2 *testing.T) (e *httpexpect.Expect) {
	e = httpexpect.New(t2, t.svr.URL)
	return e
}

var T *Test

func NewTest() (res *Test, err error) {
	// test
	t := &Test{
		svc: NewService(),
	}

	// test node
	t.TestNode = &models.Node{Id: primitive.NewObjectID(), Name: "test-node"}
	_ = delegate.NewModelDelegate(t.TestNode).Add()

	// test spider
	t.TestSpider = &models.Spider{Id: primitive.NewObjectID(), Name: "test-spider"}
	_ = delegate.NewModelDelegate(t.TestSpider).Add()

	// test task
	t.TestTask = &models.Task{Id: primitive.NewObjectID(), SpiderId: t.TestSpider.GetId(), NodeId: t.TestNode.GetId()}
	_ = delegate.NewModelDelegate(t.TestTask).Add()

	// test task stat
	t.TestTaskStat = &models.TaskStat{Id: t.TestTask.GetId()}
	_ = delegate.NewModelDelegate(t.TestTaskStat).Add()

	return t, nil
}

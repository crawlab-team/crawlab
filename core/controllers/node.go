package controllers

import (
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var NodeController *nodeController

type nodeController struct {
	ListControllerDelegate
}

func (ctr *nodeController) Post(c *gin.Context) {
	var n models.Node
	if err := c.ShouldBindJSON(&n); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	if err := ctr._post(c, &n); err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}

	HandleSuccess(c)
}

func (ctr *nodeController) PostList(c *gin.Context) {
	// bind
	var docs []models.Node
	if err := c.ShouldBindJSON(&docs); err != nil {
		HandleErrorBadRequest(c, err)
		return
	}

	// success ids
	var ids []primitive.ObjectID

	// iterate nodes
	for _, n := range docs {
		if err := ctr._post(c, &n); err != nil {
			trace.PrintError(err)
			continue
		}
		ids = append(ids, n.Id)
	}

	// success
	HandleSuccessWithData(c, docs)
}

func (ctr *nodeController) _post(c *gin.Context, n *models.Node) (err error) {
	// set default key
	if n.Key == "" {
		id, err := uuid.NewUUID()
		if err != nil {
			return trace.TraceError(err)
		}
		n.Key = id.String()
	}

	// set default status
	if n.Status == "" {
		n.Status = constants.NodeStatusUnregistered
	}

	// add
	if err := delegate.NewModelDelegate(n, GetUserFromContext(c)).Add(); err != nil {
		return trace.TraceError(err)
	}

	return nil
}

func newNodeController() *nodeController {
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}

	ctr := NewListControllerDelegate(ControllerIdNode, modelSvc.GetBaseService(interfaces.ModelIdNode))

	return &nodeController{
		ListControllerDelegate: *ctr,
	}
}

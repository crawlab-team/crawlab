package controllers

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ProjectController *projectController

type projectController struct {
	ListControllerDelegate
}

func (ctr *projectController) GetList(c *gin.Context) {
	// get all if query field "all" is set true
	all := MustGetFilterAll(c)
	if all {
		ctr.getAll(c)
		return
	}

	// get list
	list, total, err := ctr.getList(c)
	if err != nil {
		return
	}
	data := list.GetModels()

	// check empty list
	if len(list.GetModels()) == 0 {
		HandleSuccessWithListData(c, nil, 0)
		return
	}

	// project ids
	var ids []primitive.ObjectID

	// count cache
	cache := map[primitive.ObjectID]int{}

	// iterate
	for _, d := range data {
		p, ok := d.(*models.Project)
		if !ok {
			HandleErrorInternalServerError(c, errors.ErrorControllerInvalidType)
			return
		}
		ids = append(ids, p.Id)
		cache[p.Id] = 0
	}

	// spiders
	modelSvc, err := service.NewService()
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	spiders, err := modelSvc.GetSpiderList(bson.M{
		"project_id": bson.M{
			"$in": ids,
		},
	}, nil)
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	for _, s := range spiders {
		_, ok := cache[s.ProjectId]
		if !ok {
			HandleErrorInternalServerError(c, errors.ErrorControllerMissingInCache)
			return
		}
		cache[s.ProjectId]++
	}

	// assign
	var projects []models.Project
	for _, d := range data {
		p := d.(*models.Project)
		p.Spiders = cache[p.Id]
		projects = append(projects, *p)
	}

	HandleSuccessWithListData(c, projects, total)
}

func newProjectController() *projectController {
	modelSvc, err := service.GetService()
	if err != nil {
		panic(err)
	}

	ctr := NewListControllerDelegate(ControllerIdProject, modelSvc.GetBaseService(interfaces.ModelIdProject))

	return &projectController{
		ListControllerDelegate: *ctr,
	}
}

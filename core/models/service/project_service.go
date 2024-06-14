package service

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	models2 "github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeProject(d interface{}, err error) (res *models2.Project, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.Project)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetProjectById(id primitive.ObjectID) (res *models2.Project, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdProject).GetById(id)
	return convertTypeProject(d, err)
}

func (svc *Service) GetProject(query bson.M, opts *mongo.FindOptions) (res *models2.Project, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdProject).Get(query, opts)
	return convertTypeProject(d, err)
}

func (svc *Service) GetProjectList(query bson.M, opts *mongo.FindOptions) (res []models2.Project, err error) {
	l, err := svc.GetBaseService(interfaces.ModelIdProject).GetList(query, opts)
	for _, doc := range l.GetModels() {
		d := doc.(*models2.Project)
		res = append(res, *d)
	}
	return res, nil
}

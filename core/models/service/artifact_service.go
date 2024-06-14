package service

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	models2 "github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeArtifact(d interface{}, err error) (res *models2.Artifact, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.Artifact)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetArtifactById(id primitive.ObjectID) (res *models2.Artifact, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdArtifact).GetById(id)
	return convertTypeArtifact(d, err)
}

func (svc *Service) GetArtifact(query bson.M, opts *mongo.FindOptions) (res *models2.Artifact, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdArtifact).Get(query, opts)
	return convertTypeArtifact(d, err)
}

func (svc *Service) GetArtifactList(query bson.M, opts *mongo.FindOptions) (res []models2.Artifact, err error) {
	l, err := svc.GetBaseService(interfaces.ModelIdArtifact).GetList(query, opts)
	for _, doc := range l.GetModels() {
		d := doc.(*models2.Artifact)
		res = append(res, *d)
	}
	return res, nil
}

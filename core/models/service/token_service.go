package service

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	models2 "github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeToken(d interface{}, err error) (res *models2.Token, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.Token)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetTokenById(id primitive.ObjectID) (res *models2.Token, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdToken).GetById(id)
	return convertTypeToken(d, err)
}

func (svc *Service) GetToken(query bson.M, opts *mongo.FindOptions) (res *models2.Token, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdToken).Get(query, opts)
	return convertTypeToken(d, err)
}

func (svc *Service) GetTokenList(query bson.M, opts *mongo.FindOptions) (res []models2.Token, err error) {
	l, err := svc.GetBaseService(interfaces.ModelIdToken).GetList(query, opts)
	for _, doc := range l.GetModels() {
		d := doc.(*models2.Token)
		res = append(res, *d)
	}
	return res, nil
}

package service

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	models2 "github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertTypeUser(d interface{}, err error) (res *models2.User, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.User)
	if !ok {
		return nil, errors.ErrorModelInvalidType
	}
	return res, nil
}

func (svc *Service) GetUserById(id primitive.ObjectID) (res *models2.User, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdUser).GetById(id)
	return convertTypeUser(d, err)
}

func (svc *Service) GetUser(query bson.M, opts *mongo.FindOptions) (res *models2.User, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdUser).Get(query, opts)
	return convertTypeUser(d, err)
}

func (svc *Service) GetUserList(query bson.M, opts *mongo.FindOptions) (res []models2.User, err error) {
	l, err := svc.GetBaseService(interfaces.ModelIdUser).GetList(query, opts)
	for _, doc := range l.GetModels() {
		d := doc.(*models2.User)
		res = append(res, *d)
	}
	return res, nil
}

func (svc *Service) GetUserByUsername(username string, opts *mongo.FindOptions) (res *models2.User, err error) {
	query := bson.M{"username": username}
	return svc.GetUser(query, opts)
}

func (svc *Service) GetUserByUsernameWithPassword(username string, opts *mongo.FindOptions) (res *models2.User, err error) {
	u, err := svc.GetUserByUsername(username, opts)
	if err != nil {
		return nil, err
	}
	p, err := svc.GetPasswordById(u.Id)
	if err != nil {
		return nil, err
	}
	u.Password = p.Password
	return u, nil
}

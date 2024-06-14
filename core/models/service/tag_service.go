package service

import (
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	models2 "github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
)

func convertTypeTag(d interface{}, err error) (res *models2.Tag, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*models2.Tag)
	if !ok {
		return nil, trace.TraceError(errors.ErrorModelInvalidType)
	}
	return res, nil
}

func (svc *Service) GetTagById(id primitive.ObjectID) (res *models2.Tag, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdTag).GetById(id)
	return convertTypeTag(d, err)
}

func (svc *Service) GetTag(query bson.M, opts *mongo.FindOptions) (res *models2.Tag, err error) {
	d, err := svc.GetBaseService(interfaces.ModelIdTag).Get(query, opts)
	return convertTypeTag(d, err)
}

func (svc *Service) GetTagList(query bson.M, opts *mongo.FindOptions) (res []models2.Tag, err error) {
	l, err := svc.GetBaseService(interfaces.ModelIdTag).GetList(query, opts)
	for _, doc := range l.GetModels() {
		d := doc.(*models2.Tag)
		res = append(res, *d)
	}
	return res, nil
}

func (svc *Service) GetTagIds(colName string, tags []interfaces.Tag) (tagIds []primitive.ObjectID, err error) {
	// iterate tag names
	for _, tag := range tags {
		// count of tags with the name
		tagDb, err := svc.GetTag(bson.M{"name": tag.GetName(), "col": colName}, nil)
		if err == nil {
			// tag exists
			tag = tagDb
		} else if err == mongo2.ErrNoDocuments {
			// add new tag if not exists
			colorHex := tag.GetColor()
			if colorHex == "" {
				color, _ := svc.colorSvc.GetRandom()
				colorHex = color.GetHex()
			}
			tag = &models2.Tag{
				Id:    primitive.NewObjectID(),
				Name:  tag.GetName(),
				Color: colorHex,
				Col:   colName,
			}
			if err := delegate.NewModelDelegate(tag).Add(); err != nil {
				return tagIds, trace.TraceError(err)
			}
		}

		// add to tag ids
		tagIds = append(tagIds, tag.GetId())
	}

	return tagIds, nil
}

func (svc *Service) UpdateTagsById(colName string, id primitive.ObjectID, tags []interfaces.Tag) (tagIds []primitive.ObjectID, err error) {
	// get tag ids to update
	tagIds, err = svc.GetTagIds(colName, tags)
	if err != nil {
		return tagIds, trace.TraceError(err)
	}

	// update in db
	a, err := svc.GetArtifactById(id)
	if err != nil {
		return tagIds, trace.TraceError(err)
	}
	a.TagIds = tagIds
	if err := mongo.GetMongoCol(interfaces.ModelColNameArtifact).ReplaceId(id, a); err != nil {
		return tagIds, err
	}
	return tagIds, nil
}

func (svc *Service) UpdateTags(colName string, query bson.M, tags []interfaces.Tag) (tagIds []primitive.ObjectID, err error) {
	// tag ids to update
	tagIds, err = svc.GetTagIds(colName, tags)
	if err != nil {
		return tagIds, trace.TraceError(err)
	}

	// update
	update := bson.M{
		"_tid": tagIds,
	}

	// fields
	fields := []string{"_tid"}

	// update in db
	if err := svc.GetBaseService(interfaces.ModelIdTag).Update(query, update, fields); err != nil {
		return tagIds, trace.TraceError(err)
	}

	return tagIds, nil
}

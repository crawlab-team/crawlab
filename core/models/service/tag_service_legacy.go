package service

//
//import (
//	"github.com/crawlab-team/crawlab/core/interfaces"
//	"github.com/crawlab-team/crawlab/db/mongo"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/bson/primitive"
//	mongo2 "go.mongodb.org/mongo-driver/mongo"
//)
//
//type TagServiceInterface interface {
//	getTagIds(colName string, tags []Tag) (tagIds []primitive.ObjectID, err error)
//	GetModelById(id primitive.ObjectID) (res Tag, err error)
//	GetModel(query bson.M, opts *mongo.FindOptions) (res Tag, err error)
//	GetModelList(query bson.M, opts *mongo.FindOptions) (res []Tag, err error)
//	UpdateTagsById(colName string, id primitive.ObjectID, tags []Tag) (tagIds []primitive.ObjectID, err error)
//	UpdateTags(colName string, query bson.M, tags []Tag) (tagIds []primitive.ObjectID, err error)
//}
//
//type tagService struct {
//	*baseService
//}
//
//func (svc *tagService) getTagIds(colName string, tags []Tag) (tagIds []primitive.ObjectID, err error) {
//	// iterate tag names
//	for _, tag := range tags {
//		// count of tags with the name
//		tagDb, err := MustGetRootService().GetTag(bson.M{"name": tag.Name, "col": colName}, nil)
//		if err == nil {
//			// tag exists
//			tag = tagDb
//		} else if err == mongo2.ErrNoDocuments {
//			// add new tag if not exists
//			colorHex := tag.Color
//			if colorHex == "" {
//				color, _ := ColorService.GetRandom()
//				colorHex = color.Hex
//			}
//			tag = Tag{
//				Name:  tag.Name,
//				Color: colorHex,
//				Col:   colName,
//			}
//			if err := tag.Add(); err != nil {
//				return tagIds, err
//			}
//		}
//
//		// add to tag ids
//		tagIds = append(tagIds, tag.Id)
//	}
//
//	return tagIds, nil
//}
//
//func (svc *tagService) GetModelById(id primitive.ObjectID) (res Tag, err error) {
//	err = svc.findId(id).One(&res)
//	return res, err
//}
//
//func (svc *tagService) GetModel(query bson.M, opts *mongo.FindOptions) (res Tag, err error) {
//	err = svc.find(query, opts).One(&res)
//	return res, err
//}
//
//func (svc *tagService) GetModelList(query bson.M, opts *mongo.FindOptions) (res []Tag, err error) {
//	err = svc.find(query, opts).All(&res)
//	return res, err
//}
//
//func (svc *tagService) UpdateTagsById(colName string, id primitive.ObjectID, tags []Tag) (tagIds []primitive.ObjectID, err error) {
//	// get tag ids to update
//	tagIds, err = svc.getTagIds(colName, tags)
//	if err != nil {
//		return tagIds, err
//	}
//
//	// update in db
//	a, err := MustGetRootService().GetArtifactById(id)
//	if err != nil {
//		return tagIds, err
//	}
//	a.TagIds = tagIds
//	if err := mongo.GetMongoCol(interfaces.ModelColNameArtifact).ReplaceId(id, a); err != nil {
//		return tagIds, err
//	}
//	return tagIds, nil
//}
//
//func (svc *tagService) UpdateTags(colName string, query bson.M, tags []Tag) (tagIds []primitive.ObjectID, err error) {
//	// tag ids to update
//	tagIds, err = svc.getTagIds(colName, tags)
//	if err != nil {
//		return tagIds, err
//	}
//
//	// update
//	update := bson.M{
//		"_tid": tagIds,
//	}
//
//	// fields
//	fields := []string{"_tid"}
//
//	// update in db
//	if err := ArtifactService.Update(query, update, fields); err != nil {
//		return tagIds, err
//	}
//
//	return tagIds, nil
//}
//
//func NewTagService() (svc *tagService) {
//	return &tagService{svc.GetBaseService(interfaces.ModelIdTag)}
//}
//
//var TagService *tagService

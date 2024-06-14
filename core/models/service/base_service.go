package service

import (
	"encoding/json"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/delegate"
	models2 "github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strings"
	"sync"
	"time"
)

type BaseService struct {
	id  interfaces.ModelId
	col *mongo.Col
}

func (svc *BaseService) GetModelId() (id interfaces.ModelId) {
	return svc.id
}

func (svc *BaseService) SetModelId(id interfaces.ModelId) {
	svc.id = id
}

func (svc *BaseService) GetCol() (col *mongo.Col) {
	return svc.col
}

func (svc *BaseService) SetCol(col *mongo.Col) {
	svc.col = col
}

func (svc *BaseService) GetById(id primitive.ObjectID) (res interfaces.Model, err error) {
	// find result
	fr := svc.findId(id)

	// bind
	return NewBasicBinder(svc.id, fr).Bind()
}

func (svc *BaseService) Get(query bson.M, opts *mongo.FindOptions) (res interfaces.Model, err error) {
	// find result
	fr := svc.find(query, opts)

	// bind
	return NewBasicBinder(svc.id, fr).Bind()
}

func (svc *BaseService) GetList(query bson.M, opts *mongo.FindOptions) (l interfaces.List, err error) {
	// find result
	tic := time.Now()
	log.Debugf("baseService.GetMany -> svc.find:start")
	log.Debugf("baseService.GetMany -> svc.id: %v", svc.id)
	log.Debugf("baseService.GetMany -> svc.col.GetName(): %v", svc.col.GetName())
	log.Debugf("baseService.GetMany -> query: %v", query)
	log.Debugf("baseService.GetMany -> opts: %v", opts)
	fr := svc.find(query, opts)
	log.Debugf("baseService.GetMany -> svc.find:end. elapsed: %d ms", time.Now().Sub(tic).Milliseconds())

	// bind
	return NewListBinder(svc.id, fr).Bind()
}

func (svc *BaseService) DeleteById(id primitive.ObjectID, args ...interface{}) (err error) {
	return svc.deleteId(id, args...)
}

func (svc *BaseService) Delete(query bson.M, args ...interface{}) (err error) {
	return svc.delete(query)
}

func (svc *BaseService) DeleteList(query bson.M, args ...interface{}) (err error) {
	return svc.deleteList(query)
}

func (svc *BaseService) ForceDeleteList(query bson.M, args ...interface{}) (err error) {
	return svc.forceDeleteList(query)
}

func (svc *BaseService) UpdateById(id primitive.ObjectID, update bson.M, args ...interface{}) (err error) {
	return svc.updateId(id, update)
}

func (svc *BaseService) Update(query bson.M, update bson.M, fields []string, args ...interface{}) (err error) {
	return svc.update(query, update, fields)
}

func (svc *BaseService) UpdateDoc(query bson.M, doc interfaces.Model, fields []string, args ...interface{}) (err error) {
	return svc.update(query, doc, fields)
}

func (svc *BaseService) Insert(u interfaces.User, docs ...interface{}) (err error) {
	log.Debugf("baseService.Insert -> svc.col.GetName(): %v", svc.col.GetName())
	log.Debugf("baseService.Insert -> docs: %v", docs)
	return svc.insert(u, docs...)
}

func (svc *BaseService) Count(query bson.M) (total int, err error) {
	return svc.count(query)
}

func (svc *BaseService) findId(id primitive.ObjectID) (fr *mongo.FindResult) {
	if svc.col == nil {
		return mongo.NewFindResultWithError(constants.ErrMissingCol)
	}
	return svc.col.FindId(id)
}

func (svc *BaseService) find(query bson.M, opts *mongo.FindOptions) (fr *mongo.FindResult) {
	if svc.col == nil {
		return mongo.NewFindResultWithError(constants.ErrMissingCol)
	}
	return svc.col.Find(query, opts)
}

func (svc *BaseService) deleteId(id primitive.ObjectID, args ...interface{}) (err error) {
	if svc.col == nil {
		return trace.TraceError(constants.ErrMissingCol)
	}
	fr := svc.findId(id)
	doc, err := NewBasicBinder(svc.id, fr).Bind()
	if err != nil {
		return err
	}
	return delegate.NewModelDelegate(doc, svc._getUserFromArgs(args...)).Delete()
}

func (svc *BaseService) delete(query bson.M, args ...interface{}) (err error) {
	if svc.col == nil {
		return trace.TraceError(constants.ErrMissingCol)
	}
	var doc models2.BaseModel
	if err := svc.find(query, nil).One(&doc); err != nil {
		return err
	}
	return svc.deleteId(doc.Id, svc._getUserFromArgs(args...))
}

func (svc *BaseService) deleteList(query bson.M, args ...interface{}) (err error) {
	if svc.col == nil {
		return trace.TraceError(constants.ErrMissingCol)
	}
	fr := svc.find(query, nil)
	list, err := NewListBinder(svc.id, fr).Bind()
	if err != nil {
		return err
	}
	for _, doc := range list.GetModels() {
		if err := delegate.NewModelDelegate(doc, svc._getUserFromArgs(args...)).Delete(); err != nil {
			return err
		}
	}
	return nil
}

func (svc *BaseService) forceDeleteList(query bson.M, args ...interface{}) (err error) {
	return svc.col.Delete(query)
}

func (svc *BaseService) count(query bson.M) (total int, err error) {
	if svc.col == nil {
		return total, trace.TraceError(constants.ErrMissingCol)
	}
	return svc.col.Count(query)
}

func (svc *BaseService) update(query bson.M, update interface{}, fields []string, args ...interface{}) (err error) {
	update, err = svc._getUpdateBsonM(update, fields)
	if err != nil {
		return err
	}
	return svc._update(query, update, svc._getUserFromArgs(args...))
}

func (svc *BaseService) updateId(id primitive.ObjectID, update interface{}, args ...interface{}) (err error) {
	update, err = svc._getUpdateBsonM(update, nil)
	if err != nil {
		return err
	}
	return svc._updateById(id, update, svc._getUserFromArgs(args...))
}

func (svc *BaseService) insert(u interfaces.User, docs ...interface{}) (err error) {
	// validate col
	if svc.col == nil {
		return trace.TraceError(constants.ErrMissingCol)
	}

	// iterate docs
	for i, doc := range docs {
		switch doc.(type) {
		case map[string]interface{}:
			// doc type: map[string]interface{}, need to handle _id
			d := doc.(map[string]interface{})
			vId, ok := d["_id"]
			if !ok {
				// _id not exists
				d["_id"] = primitive.NewObjectID()
			} else {
				// _id exists
				switch vId.(type) {
				case string:
					// _id type: string
					sId, ok := vId.(string)
					if ok {
						d["_id"], err = primitive.ObjectIDFromHex(sId)
						if err != nil {
							return trace.TraceError(err)
						}
					}
				case primitive.ObjectID:
					// _id type: primitive.ObjectID
					// do nothing
				default:
					return trace.TraceError(errors.ErrorModelInvalidType)
				}
			}
		}
		docs[i] = doc
	}

	// perform insert
	ids, err := svc.col.InsertMany(docs)
	if err != nil {
		return err
	}

	// upsert artifacts
	query := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}
	fr := svc.col.Find(query, nil)
	list, err := NewListBinder(svc.id, fr).Bind()
	for _, doc := range list.GetModels() {
		// upsert artifact when performing model delegate save
		if err := delegate.NewModelDelegate(doc, u).Save(); err != nil {
			return err
		}
	}

	return nil
}

func (svc *BaseService) _update(query bson.M, update interface{}, args ...interface{}) (err error) {
	// ids of query
	var ids []primitive.ObjectID
	list, err := NewListBinder(svc.id, svc.find(query, nil)).Bind()
	if err != nil {
		return err
	}
	for _, doc := range list.GetModels() {
		ids = append(ids, doc.GetId())
	}

	// update model objects
	if err := svc.col.Update(query, update); err != nil {
		return err
	}

	// update artifacts
	u := svc._getUserFromArgs(args...)
	return mongo.GetMongoCol(interfaces.ModelColNameArtifact).Update(query, svc._getUpdateArtifactUpdate(u))
}

func (svc *BaseService) _updateById(id primitive.ObjectID, update interface{}, args ...interface{}) (err error) {
	// update model object
	if err := svc.col.UpdateId(id, update); err != nil {
		return err
	}

	// update artifact
	u := svc._getUserFromArgs(args...)
	return mongo.GetMongoCol(interfaces.ModelColNameArtifact).UpdateId(id, svc._getUpdateArtifactUpdate(u))
}

func (svc *BaseService) _getUpdateBsonM(update interface{}, fields []string) (res bson.M, err error) {
	switch update.(type) {
	case interfaces.Model:
		// convert to bson.M
		var updateBsonM bson.M
		bytes, err := json.Marshal(&update)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(bytes, &updateBsonM); err != nil {
			return nil, err
		}
		return svc._getUpdateBsonM(updateBsonM, fields)

	case bson.M:
		// convert to bson.M
		updateBsonM := update.(bson.M)

		// filter fields if not nil
		if fields != nil {
			// fields map
			fieldsMap := map[string]bool{}
			for _, f := range fields {
				fieldsMap[f] = true
			}

			// remove unselected fields
			for k := range updateBsonM {
				if _, ok := fieldsMap[k]; !ok {
					delete(updateBsonM, k)
				}
			}
		}

		// normalize update bson.M
		if !svc._containsDollar(updateBsonM) {
			if _, ok := updateBsonM["$set"]; !ok {
				updateBsonM = bson.M{
					"$set": updateBsonM,
				}
			}
		}

		return updateBsonM, nil
	}

	v := reflect.ValueOf(update)
	switch v.Kind() {
	case reflect.Struct:
		if v.CanAddr() {
			update = v.Addr().Interface()
			return svc._getUpdateBsonM(update, fields)
		}
		return nil, errors.ErrorModelInvalidType
	default:
		return nil, errors.ErrorModelInvalidType
	}
}

func (svc *BaseService) _getUpdateArtifactUpdate(u interfaces.User) (res bson.M) {
	var uid primitive.ObjectID
	if u != nil {
		uid = u.GetId()
	}
	return bson.M{
		"$set": bson.M{
			"_sys.update_ts":  time.Now(),
			"_sys.update_uid": uid,
		},
	}
}

func (svc *BaseService) _getUserFromArgs(args ...interface{}) (u interfaces.User) {
	return utils.GetUserFromArgs(args...)
}

func (svc *BaseService) _containsDollar(updateBsonM bson.M) (ok bool) {
	for k := range updateBsonM {
		if strings.HasPrefix(k, "$") {
			return true
		}
	}
	return false
}

func NewBaseService(id interfaces.ModelId, opts ...BaseServiceOption) (svc2 interfaces.ModelBaseService) {
	// service
	svc := &BaseService{
		id: id,
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// get collection name if not set
	if svc.GetCol() == nil {
		colName := models2.GetModelColName(id)
		svc.SetCol(mongo.GetMongoCol(colName))
	}

	return svc
}

var store = sync.Map{}

func GetBaseService(id interfaces.ModelId) (svc interfaces.ModelBaseService) {
	res, ok := store.Load(id)
	if ok {
		svc, ok = res.(interfaces.ModelBaseService)
		if ok {
			return svc
		}
	}
	svc = NewBaseService(id)
	store.Store(id, svc)
	return svc
}

func GetBaseServiceByColName(id interfaces.ModelId, colName string) (svc interfaces.ModelBaseService) {
	res, ok := store.Load(colName)
	if ok {
		svc, ok = res.(interfaces.ModelBaseService)
		if ok {
			return svc
		}
	}
	col := mongo.GetMongoCol(colName)
	svc = NewBaseService(id, WithBaseServiceCol(col))
	store.Store(colName, svc)
	return svc
}

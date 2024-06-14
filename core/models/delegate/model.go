package delegate

import (
	"encoding/json"
	errors2 "github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/event"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/db/errors"
	"github.com/crawlab-team/crawlab/db/mongo"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"time"
)

func NewModelDelegate(doc interfaces.Model, args ...interface{}) interfaces.ModelDelegate {
	switch doc.(type) {
	case *models.Artifact:
		return newModelDelegate(interfaces.ModelIdArtifact, doc, args...)
	case *models.Tag:
		return newModelDelegate(interfaces.ModelIdTag, doc, args...)
	case *models.Node:
		return newModelDelegate(interfaces.ModelIdNode, doc, args...)
	case *models.Project:
		return newModelDelegate(interfaces.ModelIdProject, doc, args...)
	case *models.Spider:
		return newModelDelegate(interfaces.ModelIdSpider, doc, args...)
	case *models.Task:
		return newModelDelegate(interfaces.ModelIdTask, doc, args...)
	case *models.Job:
		return newModelDelegate(interfaces.ModelIdJob, doc, args...)
	case *models.Schedule:
		return newModelDelegate(interfaces.ModelIdSchedule, doc, args...)
	case *models.User:
		return newModelDelegate(interfaces.ModelIdUser, doc, args...)
	case *models.Setting:
		return newModelDelegate(interfaces.ModelIdSetting, doc, args...)
	case *models.Token:
		return newModelDelegate(interfaces.ModelIdToken, doc, args...)
	case *models.Variable:
		return newModelDelegate(interfaces.ModelIdVariable, doc, args...)
	case *models.TaskQueueItem:
		return newModelDelegate(interfaces.ModelIdTaskQueue, doc, args...)
	case *models.TaskStat:
		return newModelDelegate(interfaces.ModelIdTaskStat, doc, args...)
	case *models.SpiderStat:
		return newModelDelegate(interfaces.ModelIdSpiderStat, doc, args...)
	case *models.DataSource:
		return newModelDelegate(interfaces.ModelIdDataSource, doc, args...)
	case *models.DataCollection:
		return newModelDelegate(interfaces.ModelIdDataCollection, doc, args...)
	case *models.Result:
		return newModelDelegate(interfaces.ModelIdResult, doc, args...)
	case *models.Password:
		return newModelDelegate(interfaces.ModelIdPassword, doc, args...)
	case *models.ExtraValue:
		return newModelDelegate(interfaces.ModelIdExtraValue, doc, args...)
	case *models.Git:
		return newModelDelegate(interfaces.ModelIdGit, doc, args...)
	case *models.Role:
		return newModelDelegate(interfaces.ModelIdRole, doc, args...)
	case *models.UserRole:
		return newModelDelegate(interfaces.ModelIdUserRole, doc, args...)
	case *models.Permission:
		return newModelDelegate(interfaces.ModelIdPermission, doc, args...)
	case *models.RolePermission:
		return newModelDelegate(interfaces.ModelIdRolePermission, doc, args...)
	case *models.Environment:
		return newModelDelegate(interfaces.ModelIdEnvironment, doc, args...)
	case *models.DependencySetting:
		return newModelDelegate(interfaces.ModelIdDependencySetting, doc, args...)
	default:
		_ = trace.TraceError(errors2.ErrorModelInvalidType)
		return nil
	}
}

func newModelDelegate(id interfaces.ModelId, doc interfaces.Model, args ...interface{}) interfaces.ModelDelegate {
	// user
	u := utils.GetUserFromArgs(args...)

	// collection name
	colName := models.GetModelColName(id)

	// model delegate
	d := &ModelDelegate{
		id:      id,
		colName: colName,
		doc:     doc,
		a: &models.Artifact{
			Col: colName,
		},
		u: u,
	}

	return d
}

type ModelDelegate struct {
	id      interfaces.ModelId
	colName string
	doc     interfaces.Model         // doc to delegate
	cd      bson.M                   // current doc
	od      bson.M                   // original doc
	a       interfaces.ModelArtifact // artifact
	u       interfaces.User          // user
}

// Add model
func (d *ModelDelegate) Add() (err error) {
	return d.do(interfaces.ModelDelegateMethodAdd)
}

// Save model
func (d *ModelDelegate) Save() (err error) {
	return d.do(interfaces.ModelDelegateMethodSave)
}

// Delete model
func (d *ModelDelegate) Delete() (err error) {
	return d.do(interfaces.ModelDelegateMethodDelete)
}

// GetArtifact refresh artifact and return it
func (d *ModelDelegate) GetArtifact() (res interfaces.ModelArtifact, err error) {
	if err := d.do(interfaces.ModelDelegateMethodGetArtifact); err != nil {
		return nil, err
	}
	return d.a, nil
}

// Refresh model
func (d *ModelDelegate) Refresh() (err error) {
	return d.refresh()
}

// GetModel return model
func (d *ModelDelegate) GetModel() (res interfaces.Model) {
	return d.doc
}

func (d *ModelDelegate) ToBytes(m interface{}) (bytes []byte, err error) {
	if m != nil {
		return utils.JsonToBytes(m)
	}
	return json.Marshal(d.doc)
}

// do action given the model delegate method
func (d *ModelDelegate) do(method interfaces.ModelDelegateMethod) (err error) {
	switch method {
	case interfaces.ModelDelegateMethodAdd:
		err = d.add()
	case interfaces.ModelDelegateMethodSave:
		err = d.save()
	case interfaces.ModelDelegateMethodDelete:
		err = d.delete()
	case interfaces.ModelDelegateMethodGetArtifact, interfaces.ModelDelegateMethodRefresh:
		err = d.refresh()
	default:
		return trace.TraceError(errors2.ErrorModelInvalidType)
	}

	if err != nil {
		return err
	}

	// trigger event
	eventName := GetEventName(d, method)
	go event.SendEvent(eventName, d.doc)

	return nil
}

// add model
func (d *ModelDelegate) add() (err error) {
	if d.doc == nil {
		return trace.TraceError(errors.ErrMissingValue)
	}
	if d.doc.GetId().IsZero() {
		d.doc.SetId(primitive.NewObjectID())
	}
	col := mongo.GetMongoCol(d.colName)
	if _, err = col.Insert(d.doc); err != nil {
		return trace.TraceError(err)
	}
	if err := d.upsertArtifact(); err != nil {
		return trace.TraceError(err)
	}
	return d.refresh()
}

// save model
func (d *ModelDelegate) save() (err error) {
	// validate
	if d.doc == nil || d.doc.GetId().IsZero() {
		return trace.TraceError(errors.ErrMissingValue)
	}

	// collection
	col := mongo.GetMongoCol(d.colName)

	// current doc
	docData, err := bson.Marshal(d.doc)
	if err != nil {
		trace.PrintError(err)
	} else {
		if err := bson.Unmarshal(docData, &d.cd); err != nil {
			trace.PrintError(err)
		}
	}

	// original doc
	if err := col.FindId(d.doc.GetId()).One(&d.od); err != nil {
		trace.PrintError(err)
	}

	// replace
	if err := col.ReplaceId(d.doc.GetId(), d.doc); err != nil {
		return trace.TraceError(err)
	}

	// upsert artifact
	if err := d.upsertArtifact(); err != nil {
		return trace.TraceError(err)
	}

	return d.refresh()
}

// delete model
func (d *ModelDelegate) delete() (err error) {
	if d.doc.GetId().IsZero() {
		return trace.TraceError(errors2.ErrorModelMissingId)
	}
	col := mongo.GetMongoCol(d.colName)
	if err := col.FindId(d.doc.GetId()).One(d.doc); err != nil {
		return trace.TraceError(err)
	}
	if err := col.DeleteId(d.doc.GetId()); err != nil {
		return trace.TraceError(err)
	}
	return d.deleteArtifact()
}

// refresh model and artifact
func (d *ModelDelegate) refresh() (err error) {
	if d.doc.GetId().IsZero() {
		return trace.TraceError(errors2.ErrorModelMissingId)
	}
	col := mongo.GetMongoCol(d.colName)
	fr := col.FindId(d.doc.GetId())
	if err := fr.One(d.doc); err != nil {
		return trace.TraceError(err)
	}
	return d.refreshArtifact()
}

// refresh artifact
func (d *ModelDelegate) refreshArtifact() (err error) {
	if d.doc.GetId().IsZero() {
		return trace.TraceError(errors2.ErrorModelMissingId)
	}
	col := mongo.GetMongoCol(interfaces.ModelColNameArtifact)
	if err := col.FindId(d.doc.GetId()).One(d.a); err != nil {
		return trace.TraceError(err)
	}
	return nil
}

// upsertArtifact
func (d *ModelDelegate) upsertArtifact() (err error) {
	// skip
	if d._skip() {
		return nil
	}

	// validate id
	if d.doc.GetId().IsZero() {
		return trace.TraceError(errors.ErrMissingValue)
	}

	// mongo collection
	col := mongo.GetMongoCol(interfaces.ModelColNameArtifact)

	// assign id to artifact
	d.a.SetId(d.doc.GetId())

	// attempt to find artifact
	if err := col.FindId(d.doc.GetId()).One(d.a); err != nil {
		if err == mongo2.ErrNoDocuments {
			// new artifact
			d.a.GetSys().SetCreateTs(time.Now())
			d.a.GetSys().SetUpdateTs(time.Now())
			if d.u != nil && !reflect.ValueOf(d.u).IsZero() {
				d.a.GetSys().SetCreateUid(d.u.GetId())
				d.a.GetSys().SetUpdateUid(d.u.GetId())
			}
			_, err = col.Insert(d.a)
			if err != nil {
				return trace.TraceError(err)
			}
			return nil
		} else {
			// error
			return trace.TraceError(err)
		}
	}

	// existing artifact
	d.a.GetSys().SetUpdateTs(time.Now())
	if d.u != nil {
		d.a.GetSys().SetUpdateUid(d.u.GetId())
	}

	// save new artifact
	return col.ReplaceId(d.a.GetId(), d.a)
}

// deleteArtifact
func (d *ModelDelegate) deleteArtifact() (err error) {
	// skip
	if d._skip() {
		return nil
	}

	if d.doc.GetId().IsZero() {
		return trace.TraceError(errors.ErrMissingValue)
	}
	col := mongo.GetMongoCol(interfaces.ModelColNameArtifact)
	d.a.SetId(d.doc.GetId())
	d.a.SetObj(d.doc)
	d.a.SetDel(true)
	d.a.GetSys().SetDeleteTs(time.Now())
	if d.u != nil {
		d.a.GetSys().SetDeleteUid(d.u.GetId())
	}
	return col.ReplaceId(d.doc.GetId(), d.a)
}

func (d *ModelDelegate) hasChange() (ok bool) {
	return !utils.BsonMEqual(d.cd, d.od)
}

func (d *ModelDelegate) _skip() (ok bool) {
	switch d.id {
	case
		interfaces.ModelIdArtifact,
		interfaces.ModelIdTaskQueue,
		interfaces.ModelIdTaskStat,
		interfaces.ModelIdSpiderStat,
		interfaces.ModelIdResult,
		interfaces.ModelIdPassword:
		return true
	default:
		return false
	}
}

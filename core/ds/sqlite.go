package ds

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	utils2 "github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SqliteService struct {
	SqlService
}

func NewDataSourceSqliteService(colId primitive.ObjectID, dsId primitive.ObjectID) (svc2 interfaces.ResultService, err error) {
	// service
	svc := &SqliteService{}

	// dependency injection
	svc.modelSvc, err = service.GetService()
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// data source
	if dsId.IsZero() {
		svc.ds = &models.DataSource{}
	} else {
		svc.ds, err = svc.modelSvc.GetDataSourceById(dsId)
		if err != nil {
			return nil, trace.TraceError(err)
		}
	}

	// data collection
	svc.dc, err = svc.modelSvc.GetDataCollectionById(colId)
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// session
	svc.s, err = utils2.GetSqliteSession(svc.ds)
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// collection
	svc.col = svc.s.Collection(svc.dc.Name)

	return svc, nil
}

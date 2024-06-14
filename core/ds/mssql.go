package ds

import (
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/utils"
	utils2 "github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/trace"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MssqlService struct {
	SqlService
}

func NewDataSourceMssqlService(colId primitive.ObjectID, dsId primitive.ObjectID) (svc2 interfaces.ResultService, err error) {
	// service
	svc := &MssqlService{}

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

	// data source defaults
	if svc.ds.Host == "" {
		svc.ds.Host = constants.DefaultHost
	}
	if svc.ds.Port == "" {
		svc.ds.Port = constants.DefaultMssqlPort
	}

	// data source password
	pwd, err := svc.modelSvc.GetPasswordById(svc.ds.Id)
	if err == nil {
		svc.ds.Password, err = utils.DecryptAES(pwd.Password)
		if err != nil {
			return nil, err
		}
	}

	// data collection
	svc.dc, err = svc.modelSvc.GetDataCollectionById(colId)
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// session
	svc.s, err = utils2.GetMssqlSession(svc.ds)
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// collection
	svc.col = svc.s.Collection(svc.dc.Name)

	return svc, nil
}

package ds

import (
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/utils"
	utils2 "github.com/crawlab-team/crawlab/core/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MysqlService struct {
	SqlService
}

func NewDataSourceMysqlService(colId primitive.ObjectID, dsId primitive.ObjectID) (svc2 interfaces.ResultService, err error) {
	// service
	svc := &MysqlService{}

	// dependency injection
	svc.modelSvc, err = service.GetService()
	if err != nil {
		return nil, err
	}

	// data source
	if dsId.IsZero() {
		svc.ds = &models.DataSource{}
	} else {
		svc.ds, err = svc.modelSvc.GetDataSourceById(dsId)
		if err != nil {
			return nil, err
		}
	}

	// data source defaults
	if svc.ds.Host == "" {
		svc.ds.Host = constants.DefaultHost
	}
	if svc.ds.Port == 0 {
		svc.ds.Port = constants.DefaultMysqlPort
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
		return nil, err
	}

	// session
	svc.s, err = utils2.GetMysqlSession(svc.ds)
	if err != nil {
		return nil, err
	}

	// collection
	svc.col = svc.s.Collection(svc.dc.Name)

	return svc, nil
}

package ds

import (
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/db/generic"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type KafkaService struct {
	// dependencies
	modelSvc service.ModelService

	// internals
	dc *models.DataCollection // models.DataCollection
	ds *models.DataSource     // models.DataSource
	c  *kafka.Conn            // kafka.Conn
	rb backoff.BackOff
	t  time.Time
}

func (svc *KafkaService) Insert(records ...interface{}) (err error) {
	var messages []kafka.Message
	for _, r := range records {
		switch r.(type) {
		case entity.Result:
			d := r.(entity.Result)
			messages = append(messages, kafka.Message{
				Topic: svc.ds.Database,
				Key:   []byte(d.GetTaskId().Hex()),
				Value: d.Bytes(),
			})
		}
	}
	_, err = svc.c.WriteMessages(messages...)
	if err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (svc *KafkaService) List(query generic.ListQuery, opts *generic.ListOptions) (results []interface{}, err error) {
	// N/A
	return nil, nil
}

func (svc *KafkaService) Count(query generic.ListQuery) (n int, err error) {
	// N/A
	return 0, nil
}

func NewDataSourceKafkaService(colId primitive.ObjectID, dsId primitive.ObjectID) (svc2 interfaces.ResultService, err error) {
	// service
	svc := &KafkaService{}

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
		svc.ds.Port = constants.DefaultKafkaPort
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

	return svc, nil
}

func (svc *KafkaService) Index(fields []string) {
	// TODO: implement me
}

func (svc *KafkaService) SetTime(t time.Time) {
	svc.t = t
}

func (svc *KafkaService) GetTime() (t time.Time) {
	return svc.t
}

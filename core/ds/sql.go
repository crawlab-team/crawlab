package ds

import (
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/crawlab-team/crawlab/core/models/service"
	utils2 "github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/db/generic"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/upper/db/v4"
	"time"
)

type SqlService struct {
	// dependencies
	modelSvc service.ModelService

	// internals
	ds  *models.DataSource
	dc  *models.DataCollection
	s   db.Session
	col db.Collection
	t   time.Time
}

func (svc *SqlService) Insert(records ...interface{}) (err error) {
	for _, d := range records {
		var r entity.Result
		switch d.(type) {
		case entity.Result:
			r = d.(entity.Result)
		default:
			continue
		}
		_r := r.Flatten()
		if _, err = svc.col.Insert(_r); err != nil {
			trace.PrintError(err)
			continue
		}
	}
	return nil
}

func (svc *SqlService) List(query generic.ListQuery, opts *generic.ListOptions) (results []interface{}, err error) {
	var docs []entity.Result
	if err := svc.col.Find(utils2.GetSqlQuery(query)).
		Offset(opts.Skip).
		Limit(opts.Limit).All(&docs); err != nil {
		return nil, trace.TraceError(err)
	}
	for i := range docs {
		d := docs[i].ToJSON()
		results = append(results, &d)
	}
	return results, nil
}

func (svc *SqlService) Count(query generic.ListQuery) (n int, err error) {
	nInt64, err := svc.col.Find(utils2.GetSqlQuery(query)).Count()
	if err != nil {
		return n, err
	}
	return int(nInt64), nil
}

func (svc *SqlService) Index(fields []string) {
	// TODO: implement me
}

func (svc *SqlService) SetTime(t time.Time) {
	svc.t = t
}

func (svc *SqlService) GetTime() (t time.Time) {
	return svc.t
}

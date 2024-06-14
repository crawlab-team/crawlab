package client

import (
	"encoding/json"
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab/core/container"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/crawlab-team/crawlab/db/mongo"
	grpc "github.com/crawlab-team/crawlab/grpc"
	"github.com/crawlab-team/crawlab/trace"
	errors2 "github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BaseServiceDelegate struct {
	// settings
	cfgPath string

	// internals
	id interfaces.ModelId
	c  interfaces.GrpcClient
}

func (d *BaseServiceDelegate) GetModelId() (id interfaces.ModelId) {
	return d.id
}

func (d *BaseServiceDelegate) SetModelId(id interfaces.ModelId) {
	d.id = id
}

func (d *BaseServiceDelegate) GetConfigPath() (path string) {
	return d.cfgPath
}

func (d *BaseServiceDelegate) SetConfigPath(path string) {
	d.cfgPath = path
}

func (d *BaseServiceDelegate) GetById(id primitive.ObjectID) (doc interfaces.Model, err error) {
	log.Debugf("[BaseServiceDelegate] get by id[%s]", id.Hex())
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Id: id})
	c := d.getModelBaseServiceClient()
	if c == nil {
		return nil, trace.TraceError(errors.ErrorModelNilPointer)
	}
	log.Debugf("[BaseServiceDelegate] get by id[%s] req: %v", id.Hex(), req)
	res, err := c.GetById(ctx, req)
	if err != nil {
		return nil, trace.TraceError(err)
	}
	log.Debugf("[BaseServiceDelegate] get by id[%s] res: %v", id.Hex(), res)
	return NewBasicBinder(d.id, res).Bind()
}

func (d *BaseServiceDelegate) Get(query bson.M, opts *mongo.FindOptions) (doc interfaces.Model, err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Query: query, FindOptions: opts})
	res, err := d.getModelBaseServiceClient().Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return NewBasicBinder(d.id, res).Bind()
}

func (d *BaseServiceDelegate) GetList(query bson.M, opts *mongo.FindOptions) (l interfaces.List, err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Query: query, FindOptions: opts})
	res, err := d.getModelBaseServiceClient().GetList(ctx, req)
	if err != nil {
		return l, err
	}
	return NewListBinder(d.id, res).Bind()
}

func (d *BaseServiceDelegate) DeleteById(id primitive.ObjectID, args ...interface{}) (err error) {
	u := utils.GetUserFromArgs(args...)
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Id: id, User: u})
	_, err = d.getModelBaseServiceClient().DeleteById(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (d *BaseServiceDelegate) Delete(query bson.M, args ...interface{}) (err error) {
	u := utils.GetUserFromArgs(args...)
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Query: query, User: u})
	_, err = d.getModelBaseServiceClient().Delete(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (d *BaseServiceDelegate) DeleteList(query bson.M, args ...interface{}) (err error) {
	u := utils.GetUserFromArgs(args...)
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Query: query, User: u})
	_, err = d.getModelBaseServiceClient().DeleteList(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (d *BaseServiceDelegate) ForceDeleteList(query bson.M, args ...interface{}) (err error) {
	u := utils.GetUserFromArgs(args...)
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Query: query, User: u})
	_, err = d.getModelBaseServiceClient().ForceDeleteList(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (d *BaseServiceDelegate) UpdateById(id primitive.ObjectID, update bson.M, args ...interface{}) (err error) {
	u := utils.GetUserFromArgs(args...)
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Id: id, Update: update, User: u})
	_, err = d.getModelBaseServiceClient().UpdateById(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (d *BaseServiceDelegate) Update(query bson.M, update bson.M, fields []string, args ...interface{}) (err error) {
	u := utils.GetUserFromArgs(args...)
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Query: query, Update: update, Fields: fields, User: u})
	_, err = d.getModelBaseServiceClient().Update(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (d *BaseServiceDelegate) UpdateDoc(query bson.M, doc interfaces.Model, fields []string, args ...interface{}) (err error) {
	u := utils.GetUserFromArgs(args...)
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Query: query, Doc: doc, Fields: fields, User: u})
	_, err = d.getModelBaseServiceClient().UpdateDoc(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (d *BaseServiceDelegate) Insert(u interfaces.User, docs ...interface{}) (err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Docs: docs, User: u})
	_, err = d.getModelBaseServiceClient().Insert(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (d *BaseServiceDelegate) Count(query bson.M) (total int, err error) {
	ctx, cancel := d.c.Context()
	defer cancel()
	req := d.mustNewRequest(&entity.GrpcBaseServiceParams{Query: query})
	res, err := d.getModelBaseServiceClient().Count(ctx, req)
	if err != nil {
		return total, err
	}
	if err := json.Unmarshal(res.Data, &total); err != nil {
		return total, err
	}
	return total, nil
}

func (d *BaseServiceDelegate) newRequest(params interfaces.GrpcBaseServiceParams) (req *grpc.Request, err error) {
	return d.c.NewModelBaseServiceRequest(d.id, params)
}

func (d *BaseServiceDelegate) mustNewRequest(params *entity.GrpcBaseServiceParams) (req *grpc.Request) {
	req, err := d.newRequest(params)
	if err != nil {
		panic(err)
	}
	return req
}

func (d *BaseServiceDelegate) getModelBaseServiceClient() (c grpc.ModelBaseServiceClient) {
	if err := backoff.Retry(func() (err error) {
		c = d.c.GetModelBaseServiceClient()
		if c == nil {
			err = errors2.New("unable to get model base service client")
			log.Debugf("[BaseServiceDelegate] err: %v", err)
			return err
		}
		return nil
	}, backoff.NewConstantBackOff(1*time.Second)); err != nil {
		trace.PrintError(err)
	}
	return c
}

func NewBaseServiceDelegate(opts ...ModelBaseServiceDelegateOption) (svc2 interfaces.GrpcClientModelBaseService, err error) {
	// base service
	svc := &BaseServiceDelegate{}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// config path
	if viper.GetString("config.path") != "" {
		svc.cfgPath = viper.GetString("config.path")
	}

	// dependency injection
	if err := container.GetContainer().Invoke(func(client interfaces.GrpcClient) {
		svc.c = client
	}); err != nil {
		return nil, err
	}

	return svc, nil
}

func ProvideBaseServiceDelegate(id interfaces.ModelId) func() (svc interfaces.GrpcClientModelBaseService, err error) {
	return func() (svc interfaces.GrpcClientModelBaseService, err error) {
		return NewBaseServiceDelegate(WithBaseServiceModelId(id))
	}
}

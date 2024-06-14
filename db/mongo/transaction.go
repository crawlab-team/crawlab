package mongo

import (
	"context"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/mongo"
)

func RunTransaction(fn func(mongo.SessionContext) error) (err error) {
	return RunTransactionWithContext(context.Background(), fn)
}

func RunTransactionWithContext(ctx context.Context, fn func(mongo.SessionContext) error) (err error) {
	// default client
	c, err := GetMongoClient()
	if err != nil {
		return err
	}

	// start session
	s, err := c.StartSession()
	if err != nil {
		return trace.TraceError(err)
	}

	// start transaction
	if err := s.StartTransaction(); err != nil {
		return trace.TraceError(err)
	}

	// perform operation
	if err := mongo.WithSession(ctx, s, func(sc mongo.SessionContext) error {
		if err := fn(sc); err != nil {
			return trace.TraceError(err)
		}
		if err = s.CommitTransaction(sc); err != nil {
			return trace.TraceError(err)
		}
		return nil
	}); err != nil {
		return trace.TraceError(err)
	}

	return nil
}

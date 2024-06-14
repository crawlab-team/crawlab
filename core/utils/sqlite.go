package utils

import (
	"context"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/sqlite"
	"time"
)

func GetSqliteSession(ds *models.DataSource) (s db.Session, err error) {
	return getSqliteSession(context.Background(), ds)
}

func GetSqliteSessionWithTimeout(ds *models.DataSource, timeout time.Duration) (s db.Session, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return getSqliteSession(ctx, ds)
}

func getSqliteSession(ctx context.Context, ds *models.DataSource) (s db.Session, err error) {
	// connect settings
	settings := sqlite.ConnectionURL{
		Database: ds.Database,
		Options:  nil,
	}

	// session
	done := make(chan struct{})
	go func() {
		s, err = sqlite.Open(settings)
		close(done)
	}()

	// wait for done
	select {
	case <-ctx.Done():
		if ctx.Err() != nil {
			err = ctx.Err()
		}
	case <-done:
	}

	return s, err
}

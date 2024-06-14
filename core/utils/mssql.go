package utils

import (
	"context"
	"fmt"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mssql"
	"time"
)

func GetMssqlSession(ds *models.DataSource) (s db.Session, err error) {
	return getMssqlSession(context.Background(), ds)
}

func GetMssqlSessionWithTimeout(ds *models.DataSource, timeout time.Duration) (s db.Session, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return getMssqlSession(ctx, ds)
}

func getMssqlSession(ctx context.Context, ds *models.DataSource) (s db.Session, err error) {
	// normalize settings
	host := ds.Host
	port := ds.Port
	if ds.Host == "" {
		host = constants.DefaultHost
	}
	if ds.Port == "" {
		port = constants.DefaultMssqlPort
	}

	// connect settings
	settings := mssql.ConnectionURL{
		User:     ds.Username,
		Password: ds.Password,
		Database: ds.Database,
		Host:     fmt.Sprintf("%s:%s", host, port),
		Options:  nil,
	}

	// session
	done := make(chan struct{})
	go func() {
		s, err = mssql.Open(settings)
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

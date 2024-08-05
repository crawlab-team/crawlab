package utils

import (
	"context"
	"fmt"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/models/models"
	models2 "github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"time"
)

func GetPostgresqlSession(ds *models.DataSource) (s db.Session, err error) {
	return getPostgresqlSession(context.Background(), ds)
}

func GetPostgresqlSessionWithTimeout(ds *models.DataSource, timeout time.Duration) (s db.Session, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return getPostgresqlSession(ctx, ds)
}

func getPostgresqlSession(ctx context.Context, ds *models.DataSource) (s db.Session, err error) {
	// normalize settings
	host := ds.Host
	port := ds.Port
	if ds.Host == "" {
		host = constants.DefaultHost
	}
	if ds.Port == 0 {
		port = constants.DefaultPostgresqlPort
	}

	// connect settings
	settings := postgresql.ConnectionURL{
		User:     ds.Username,
		Password: ds.Password,
		Database: ds.Database,
		Host:     fmt.Sprintf("%s:%s", host, port),
		Options:  nil,
	}

	// session
	done := make(chan struct{})
	go func() {
		s, err = postgresql.Open(settings)
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

func GetPostgresqlSessionWithTimeoutV2(ds *models2.DatabaseV2, timeout time.Duration) (s db.Session, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return getPostgresqlSessionV2(ctx, ds)
}

func getPostgresqlSessionV2(ctx context.Context, ds *models2.DatabaseV2) (s db.Session, err error) {
	// normalize settings
	host := ds.Host
	port := ds.Port
	if ds.Host == "" {
		host = constants.DefaultHost
	}
	if ds.Port == 0 {
		port = constants.DefaultPostgresqlPort
	}

	// connect settings
	settings := postgresql.ConnectionURL{
		User:     ds.Username,
		Password: ds.Password,
		Database: ds.Database,
		Host:     fmt.Sprintf("%s:%s", host, port),
		Options:  nil,
	}

	// session
	done := make(chan struct{})
	go func() {
		s, err = postgresql.Open(settings)
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

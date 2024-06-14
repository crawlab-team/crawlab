package utils

import (
	"context"
	"fmt"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/models/models"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
	"time"
)

func GetMysqlSession(ds *models.DataSource) (s db.Session, err error) {
	return getMysqlSession(context.Background(), ds)
}

func GetMysqlSessionWithTimeout(ds *models.DataSource, timeout time.Duration) (s db.Session, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return getMysqlSession(ctx, ds)
}

func getMysqlSession(ctx context.Context, ds *models.DataSource) (s db.Session, err error) {
	// normalize settings
	host := ds.Host
	port := ds.Port
	if ds.Host == "" {
		host = constants.DefaultHost
	}
	if ds.Port == "" {
		port = constants.DefaultMysqlPort
	}

	// connect settings
	settings := mysql.ConnectionURL{
		User:     ds.Username,
		Password: ds.Password,
		Database: ds.Database,
		Host:     fmt.Sprintf("%s:%s", host, port),
		Options:  nil,
	}

	// session
	done := make(chan struct{})
	go func() {
		s, err = mysql.Open(settings)
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

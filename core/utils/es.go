package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/models/models"
	models2 "github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/db/generic"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func GetElasticsearchClient(ds *models.DataSource) (c *elasticsearch.Client, err error) {
	return getElasticsearchClient(context.Background(), ds)
}

func GetElasticsearchClientWithTimeout(ds *models.DataSource, timeout time.Duration) (c *elasticsearch.Client, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return getElasticsearchClient(ctx, ds)
}

func getElasticsearchClient(ctx context.Context, ds *models.DataSource) (c *elasticsearch.Client, err error) {
	// normalize settings
	host := ds.Host
	port := ds.Port
	if ds.Host == "" {
		host = constants.DefaultHost
	}
	if ds.Port == 0 {
		port = constants.DefaultElasticsearchPort
	}

	// es hosts
	addresses := []string{
		fmt.Sprintf("http://%s:%s", host, port),
	}

	// retry backoff
	rb := backoff.NewExponentialBackOff()

	// es client options
	cfg := elasticsearch.Config{
		Addresses: addresses,
		Username:  ds.Username,
		Password:  ds.Password,
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				rb.Reset()
			}
			return rb.NextBackOff()
		},
	}

	// es client
	done := make(chan struct{})
	go func() {
		c, err = elasticsearch.NewClient(cfg)
		if err != nil {
			return
		}
		var res *esapi.Response
		res, err = c.Info()
		fmt.Println(res)
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

	return c, err
}

func GetElasticsearchClientWithTimeoutV2(ds *models2.DatabaseV2, timeout time.Duration) (c *elasticsearch.Client, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return getElasticsearchClientV2(ctx, ds)
}

func getElasticsearchClientV2(ctx context.Context, ds *models2.DatabaseV2) (c *elasticsearch.Client, err error) {
	// normalize settings
	host := ds.Host
	port := ds.Port
	if ds.Host == "" {
		host = constants.DefaultHost
	}
	if ds.Port == 0 {
		port = constants.DefaultElasticsearchPort
	}

	// es hosts
	addresses := []string{
		fmt.Sprintf("http://%s:%s", host, port),
	}

	// retry backoff
	rb := backoff.NewExponentialBackOff()

	// es client options
	cfg := elasticsearch.Config{
		Addresses: addresses,
		Username:  ds.Username,
		Password:  ds.Password,
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				rb.Reset()
			}
			return rb.NextBackOff()
		},
	}

	// es client
	done := make(chan struct{})
	go func() {
		c, err = elasticsearch.NewClient(cfg)
		if err != nil {
			return
		}
		var res *esapi.Response
		res, err = c.Info()
		fmt.Println(res)
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

	return c, err
}

func GetElasticsearchQuery(query generic.ListQuery) (buf *bytes.Buffer) {
	q := map[string]interface{}{}
	if len(query) > 0 {
		match := getElasticsearchQueryMatch(query)
		q["query"] = map[string]interface{}{
			"match": match,
		}
	}
	buf = &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(q); err != nil {
		trace.PrintError(err)
	}
	return buf
}

func GetElasticsearchQueryWithOptions(query generic.ListQuery, opts *generic.ListOptions) (buf *bytes.Buffer) {
	q := map[string]interface{}{
		"size": opts.Limit,
		"from": opts.Skip,
		// TODO: sort
	}
	if len(query) > 0 {
		match := getElasticsearchQueryMatch(query)
		q["query"] = map[string]interface{}{
			"match": match,
		}
	}
	buf = &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(q); err != nil {
		trace.PrintError(err)
	}
	return buf
}

func getElasticsearchQueryMatch(query generic.ListQuery) (match map[string]interface{}) {
	match = map[string]interface{}{}
	for _, c := range query {
		switch c.Value.(type) {
		case primitive.ObjectID:
			c.Value = c.Value.(primitive.ObjectID).Hex()
		}
		switch c.Op {
		case generic.OpEqual:
			match[c.Key] = c.Value
		default:
			match[c.Key] = map[string]interface{}{
				c.Op: c.Value,
			}
		}
	}
	return match
}

package utils

import (
	"context"
	"fmt"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/models/models"
	models2 "github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/segmentio/kafka-go"
	"time"
)

func GetKafkaConnection(ds *models.DataSource) (c *kafka.Conn, err error) {
	return getKafkaConnection(context.Background(), ds)
}

func GetKafkaConnectionWithTimeout(ds *models.DataSource, timeout time.Duration) (c *kafka.Conn, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return getKafkaConnection(ctx, ds)
}

func getKafkaConnection(ctx context.Context, ds *models.DataSource) (c *kafka.Conn, err error) {
	// normalize settings
	host := ds.Host
	port := ds.Port
	if ds.Host == "" {
		host = constants.DefaultHost
	}
	if ds.Port == 0 {
		port = constants.DefaultKafkaPort
	}

	// kafka connection address
	network := "tcp"
	address := fmt.Sprintf("%s:%s", host, port)
	topic := ds.Database
	partition := 0 // TODO: parameterize

	// kafka connection
	return kafka.DialLeader(ctx, network, address, topic, partition)
}

func GetKafkaConnectionWithTimeoutV2(ds *models2.DatabaseV2, timeout time.Duration) (c *kafka.Conn, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return getKafkaConnectionV2(ctx, ds)
}

func getKafkaConnectionV2(ctx context.Context, ds *models2.DatabaseV2) (c *kafka.Conn, err error) {
	// normalize settings
	host := ds.Host
	port := ds.Port
	if ds.Host == "" {
		host = constants.DefaultHost
	}
	if ds.Port == 0 {
		port = constants.DefaultKafkaPort
	}

	// kafka connection address
	network := "tcp"
	address := fmt.Sprintf("%s:%s", host, port)
	topic := ds.Database
	partition := 0 // TODO: parameterize

	// kafka connection
	return kafka.DialLeader(ctx, network, address, topic, partition)
}

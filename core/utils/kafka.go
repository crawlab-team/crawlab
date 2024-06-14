package utils

import (
	"context"
	"fmt"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/models/models"
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
	if ds.Port == "" {
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

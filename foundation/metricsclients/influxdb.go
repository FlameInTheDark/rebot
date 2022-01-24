package metricsclients

import (
	"context"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/pkg/errors"
	"time"
)

type InfluxMetrics struct {
	client influxdb2.Client
	writer api.WriteAPI
}

func NewInfluxClient(server, token, org, bucket string) (*InfluxMetrics, error) {
	client := influxdb2.NewClient(server, token)
	ctx, closefn := context.WithTimeout(context.Background(), time.Second)
	defer closefn()
	res, err := client.Ping(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to InfluxDB endpoint")
	}
	if !res {
		return nil, errors.New("wrong response of InfluxDB endpoint")
	}
	writer := client.WriteAPI(org, bucket)
	return &InfluxMetrics{
		client: client,
		writer: writer,
	}, nil
}

func (im *InfluxMetrics) Client() influxdb2.Client {
	return im.client
}

func (im *InfluxMetrics) Writer() api.WriteAPI {
	return im.writer
}

func (im *InfluxMetrics) Close() {
	im.client.Close()
}

package metrics

import (
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/FlameInTheDark/rebot/foundation/metricsclients"
)

// Metrics is a metrics for the weather service
type Metrics struct {
	client *metricsclients.InfluxMetrics
}

// NewMetrics creates a new metrics
func NewMetrics(client *metricsclients.InfluxMetrics) *Metrics {
	return &Metrics{client: client}
}

// CommandUsed create point for the command usage
func (m *Metrics) CommandUsed(command string) {
	p := influxdb2.NewPointWithMeasurement(fmt.Sprintf("command_%s", command)).
		AddField("usage", 1).
		AddTag("success", "true")
	m.client.Writer().WritePoint(p)
}

// CommandFailed create point for the command failure
func (m *Metrics) CommandFailed(command string) {
	p := influxdb2.NewPointWithMeasurement(fmt.Sprintf("command_%s", command)).
		AddField("usage", 1).
		AddTag("success", "false")
	m.client.Writer().WritePoint(p)
}

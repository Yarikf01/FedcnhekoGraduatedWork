package business

import (
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

//go:generate mockery --name MetricWriter --outpkg businessmocks --output ./businessmocks --dir .
type MetricWriter interface {
	WriteStatPoint(tags map[string]string, fields map[string]interface{})
}

type MetricConfig struct {
	ServerURL string
	Token     string
	Org       string
	Bucket    string
}

func NewMetricWriter(cfg MetricConfig) (MetricWriter, func()) {
	client := influxdb2.NewClientWithOptions(cfg.ServerURL, cfg.Token, influxdb2.DefaultOptions())
	return &writer{
		writeAPI: client.WriteAPI(cfg.Org, cfg.Bucket),
	}, client.Close
}

//impl

const statsMeasurement = "stat"

type writer struct {
	writeAPI api.WriteAPI
}

func (w *writer) WriteStatPoint(tags map[string]string, fields map[string]interface{}) {
	w.writeAPI.WritePoint(influxdb2.NewPoint(statsMeasurement, tags, fields, time.Now()))
}

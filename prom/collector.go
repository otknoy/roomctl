package prom

import (
	"context"
	"log"
	"roomctl/switchbot"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "roomctl"
)

var (
	temperatureGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "temperature",
		Help:      "temperature",
	})
	humidityGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "humidity",
		Help:      "humidity",
	})
)

var _ prometheus.Collector = (*SwitchBotSensorCollector)(nil)

type SwitchBotSensorCollector struct {
	Client switchbot.Client
}

func (c *SwitchBotSensorCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- temperatureGauge.Desc()
	ch <- humidityGauge.Desc()
}

func (c *SwitchBotSensorCollector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()

	temp, hum, err := c.Client.GetMetrics(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(temp, hum)

	temperatureGauge.Set(float64(temp))
	humidityGauge.Set(float64(hum))

	ch <- temperatureGauge
	ch <- humidityGauge
}

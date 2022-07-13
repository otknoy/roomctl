package collector

import (
	"context"
	"log"
	"roomctl/switchbot"

	"github.com/prometheus/client_golang/prometheus"
)

func NewSwitchBotSensorCollector(token, deviceId string) prometheus.Collector {
	labels := prometheus.Labels{"device": "switchbot"}

	return &switchBotSensorCollector{
		client: &switchbot.ClientImpl{
			Token:    token,
			DeviceId: deviceId,
		},
		temperatureGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   namespace,
			Name:        "temperature",
			Help:        "temperature",
			ConstLabels: labels,
		}),
		humidityGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   namespace,
			Name:        "humidity",
			Help:        "humidity",
			ConstLabels: labels,
		}),
	}
}

type switchBotSensorCollector struct {
	client switchbot.Client

	temperatureGauge prometheus.Gauge
	humidityGauge    prometheus.Gauge
}

func (c *switchBotSensorCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.temperatureGauge.Desc()
	ch <- c.humidityGauge.Desc()
}

func (c *switchBotSensorCollector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()

	temp, hum, err := c.client.GetMetrics(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(temp, hum)

	c.temperatureGauge.Set(float64(temp))
	c.humidityGauge.Set(float64(hum))

	ch <- c.temperatureGauge
	ch <- c.humidityGauge
}

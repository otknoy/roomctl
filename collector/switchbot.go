package collector

import (
	"context"
	"log"
	"roomctl/switchbot"

	"github.com/prometheus/client_golang/prometheus"
)

func NewSwitchBotSensorCollectors(token, secret string, deviceIds, deviceNames []string) []prometheus.Collector {
	l := make([]prometheus.Collector, len(deviceIds))
	for i := range l {
		l[i] = NewSwitchBotSensorCollector(token, secret, deviceIds[i], deviceNames[i])
	}

	return l
}

func NewSwitchBotSensorCollector(token, secret, deviceId, deviceName string) prometheus.Collector {
	labels := prometheus.Labels{
		"device":      "switchbot",
		"device_id":   deviceId,
		"device_name": deviceName,
	}

	return &switchBotSensorCollector{
		client: switchbot.NewCacheClient(&switchbot.ClientImpl{
			Token:    token,
			Secret:   secret,
			DeviceId: deviceId,
		}),
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

	c.temperatureGauge.Set(float64(temp))
	c.humidityGauge.Set(float64(hum))

	ch <- c.temperatureGauge
	ch <- c.humidityGauge
}

package collector

import (
	"context"
	"log"
	"roomctl/mflight"

	"github.com/prometheus/client_golang/prometheus"
)

func NewMfLightSensorCollector(url, mobileId string) prometheus.Collector {
	labels := prometheus.Labels{"device": "mflight"}

	return &mfLightSensorCollector{
		client: &mflight.ClientImpl{
			BaseUrl:  url,
			MobileId: mobileId,
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
		illuminanceGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   namespace,
			Name:        "illuminance",
			Help:        "illuminance",
			ConstLabels: labels,
		}),
	}
}

type mfLightSensorCollector struct {
	client mflight.Client

	temperatureGauge prometheus.Gauge
	humidityGauge    prometheus.Gauge
	illuminanceGauge prometheus.Gauge
}

func (c *mfLightSensorCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.temperatureGauge.Desc()
	ch <- c.humidityGauge.Desc()
	ch <- c.illuminanceGauge.Desc()
}

func (c *mfLightSensorCollector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()

	temp, hum, illu, err := c.client.GetMetrics(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(temp, hum, illu)
	log.Println(float64(temp), float64(hum), float64(illu))

	c.temperatureGauge.Set(float64(temp))
	c.humidityGauge.Set(float64(hum))
	c.illuminanceGauge.Set(float64(illu))

	ch <- c.temperatureGauge
	ch <- c.humidityGauge
	ch <- c.illuminanceGauge
}

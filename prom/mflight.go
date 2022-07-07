package prom

import (
	"github.com/prometheus/client_golang/prometheus"
)

func NewMfLightSensorCollector() prometheus.Collector {
	return &mfLightSensorCollector{
		temperatureGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "temperature",
			Help:      "temperature",
		}),
		humidityGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "humidity",
			Help:      "humidity",
		}),
		illuminanceGauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "illuminance",
			Help:      "multifunction light illuminance",
		}),
	}
}

type mfLightSensorCollector struct {
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
	// todo
	temp := 0.0
	hum := 0.0
	illu := 0.0

	c.temperatureGauge.Set(float64(temp))
	c.humidityGauge.Set(float64(hum))
	c.illuminanceGauge.Set(float64(illu))

	ch <- c.temperatureGauge
	ch <- c.humidityGauge
}

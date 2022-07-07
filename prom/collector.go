package prom

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "roomctl"
)

func NewCollector(token, deviceId string) prometheus.Collector {
	return &collector{
		NewSwitchBotSensorCollector(token, deviceId),
	}
}

type collector struct {
	switchBotSensorCollector prometheus.Collector
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	c.switchBotSensorCollector.Describe(ch)
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	c.switchBotSensorCollector.Collect(ch)
}

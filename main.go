package main

import (
	"fmt"
	"log"
	"net/http"
	"roomctl/collector"
	"roomctl/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	c, err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	http.Handle("/metrics", promhttp.Handler())

	prometheus.MustRegister(collector.NewSwitchBotSensorCollector(
		c.SwitchBot.Token,
		c.SwitchBot.DeviceId,
	))
	prometheus.MustRegister(collector.NewMfLightSensorCollector(
		c.MfLight.URL,
		c.MfLight.MobileId,
	))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", c.Port), nil); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

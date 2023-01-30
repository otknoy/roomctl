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

	for _, col := range collector.NewSwitchBotSensorCollectors(
		c.SwitchBot.Token,
		c.SwitchBot.DeviceId,
		c.SwitchBot.DeviceName,
	) {
		prometheus.MustRegister(col)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", c.Port), nil); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

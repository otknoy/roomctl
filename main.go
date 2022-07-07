package main

import (
	"log"
	"net/http"
	"roomctl/config"
	"roomctl/prom"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	c, err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	http.Handle("/metrics", promhttp.Handler())

	err = prometheus.Register(prom.NewCollector(
		c.SwitchBot.Token,
		c.SwitchBot.DeviceId,
	))
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":9199", nil); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

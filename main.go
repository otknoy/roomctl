package main

import (
	"log"
	"net/http"
	"roomctl/config"
	"roomctl/prom"
	"roomctl/switchbot"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	c, err := config.Load()
	if err != nil {
		log.Fatalln(err)
	}

	http.Handle("/metrics", promhttp.Handler())

	prometheus.Register(&prom.SwitchBotSensorCollector{
		&switchbot.ClientImpl{
			Token: c.SwitchBot.Token,
		},
	})

	if err := http.ListenAndServe(":9199", nil); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

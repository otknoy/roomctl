package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(":9199", nil); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

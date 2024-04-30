package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"roomctl/collector"
	"roomctl/config"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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

	done := make(chan error, 1)
	go func() {
		done <- http.ListenAndServe(fmt.Sprintf(":%d", c.Port), nil)
	}()

	select {
	case err := <-done:
		log.Println(err)
	case <-ctx.Done():
		log.Println("done")
	}

	// if err := http.ListenAndServe(fmt.Sprintf(":%d", c.Port), nil); err != http.ErrServerClosed {
	// 	log.Fatal(err)
	// }
}

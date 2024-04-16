package main

import (
	"os"
	"os/signal"
	"simple-telemetry-publisher/internal/prometheus"
	"simple-telemetry-publisher/internal/otel"
	"syscall"
)

func main() {
	prometheusMetrics := prometheus.PrometheusMetrics{Port: 9004}
	prometheusMetrics.Init()
	otel := otel.Otel{}
	otel.Init()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(1)
}

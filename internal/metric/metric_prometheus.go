package metric

import (
	"context"
	"math/rand"
	"net/http"
	"simple-telemetry-publisher/internal/model"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMetrics struct {
	model.PrometheusConfig
}

var (
	counterMetric = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "my_counter",
			Help: "This is my counter",
		},
	)

	gaugeMetric = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "random_number",
			Help: "Random number from 0 to 10",
		},
	)
)

func (prometheusMetrics PrometheusMetrics) Start(ctx context.Context) error {

	r := prometheus.NewRegistry()
	r.MustRegister(counterMetric, gaugeMetric)

	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	address := ":" + strconv.Itoa(int(prometheusMetrics.Port))
	srv := &http.Server{Addr: address}
	go srv.ListenAndServe()
	go prometheusMetrics.updateMetrics(ctx)

	<-ctx.Done()
	srv.Shutdown(ctx)

	return nil
}

func (prometheusMetrics PrometheusMetrics) updateMetrics(ctx context.Context) {
	ticker := time.NewTicker(prometheusMetrics.Interval)
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-ticker.C:
			counterMetric.Inc()
			gaugeMetric.Set(rand.Float64() * 10)
		}
	}
}

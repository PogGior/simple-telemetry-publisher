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
	GracefulShutdown time.Duration
	Config model.PrometheusConfig
}

var (
	counterMetric = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "counter",
			Help: "Example of a Counter",
		},
	)

	gaugeMetric = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "gauge",
			Help: "Example of a Gauge",
		},
	)
)

func (p PrometheusMetrics) Start(ctx context.Context) error {

	r := prometheus.NewRegistry()
	r.MustRegister(counterMetric, gaugeMetric)

	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{
		Registry: r,
	}))
	address := ":" + strconv.Itoa(int(p.Config.Port))
	srv := &http.Server{Addr: address} 
	go srv.ListenAndServe()
	go p.produceMetrics(ctx, srv)

	return nil
}

func (p PrometheusMetrics) produceMetrics(ctx context.Context, srv *http.Server) {
	ticker := time.NewTicker(p.Config.Interval)
loop:
	for {
		select {
		case <-ctx.Done():
			timeoutContext, cancel := context.WithTimeout(context.Background(), p.GracefulShutdown)
			defer cancel()
			srv.Shutdown(timeoutContext)
			break loop
		case <-ticker.C:
			counterMetric.Inc()
			gaugeMetric.Set(rand.Float64() * 10)
		}
	}
}

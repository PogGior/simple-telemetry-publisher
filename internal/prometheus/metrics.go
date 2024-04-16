package prometheus

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMetrics struct {
    Port uint16
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

func (prometheusMetrics PrometheusMetrics) Init() {
    // Create a new registry
    r := prometheus.NewRegistry()

    // Register your metrics with your registry
    r.MustRegister(counterMetric, gaugeMetric)

    // Use your registry for the /metrics endpoint
    http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))

    // Start the HTTP server
    address := ":" + strconv.Itoa(int(prometheusMetrics.Port))
    go http.ListenAndServe(address, nil)
    go prometheusMetrics.updateMetrics()
}

func (prometheusMetrics PrometheusMetrics) updateMetrics() {
    for range time.Tick(15 * time.Second) {
        counterMetric.Inc()
        gaugeMetric.Set(rand.Float64() * 10)
    }
}
package otel

import (
	"log"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"context"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	sdk_metric "go.opentelemetry.io/otel/sdk/metric"
	"math/rand"
	"time"
)

type Otel struct {}

func (otel *Otel) Init() {

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("ExampleService"),
		),
	)

	if err != nil {
		log.Fatalf("failed to merge resources: %v", err)
	}

	metricExporter, err := stdoutmetric.New()
	if err != nil {
		panic(err)
	}

	meterProvider := sdk_metric.NewMeterProvider(
		sdk_metric.WithResource(res),
		sdk_metric.WithReader(sdk_metric.NewPeriodicReader(metricExporter,
			sdk_metric.WithInterval(15*time.Second))),
	)

	meter := meterProvider.Meter("example-meter")

	counter, err := meter.Int64Counter("counter",
		metric.WithDescription("Example of a Counter"),
	)
	if err != nil {
		panic(err)
	}

	_, err = meter.Float64ObservableGauge("gauge",
		metric.WithDescription("Example of a Gauge"),
		metric.WithFloat64Callback(func(_ context.Context, o metric.Float64Observer) error {
			o.Observe(rand.Float64() * 10)
			return nil
		}),
	)
	if err != nil {
		panic(err)
	}

	go updateMetrics(counter)

}

func updateMetrics(counter metric.Int64Counter) {
	for range time.Tick(15 * time.Second) {
		ctx := context.Background()
		counter.Add(ctx, 1)
	}
}

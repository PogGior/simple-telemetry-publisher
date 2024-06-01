package metric

import (
	"simple-telemetry-publisher/internal/model"

	"context"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type OtelMetrics struct {
	GracefulShutdown time.Duration
	Config model.OtelMetricConfig
}

func (o OtelMetrics) Start(ctx context.Context) error {

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(o.Config.ServiceName),
		),
	)
	if err != nil {
		return errors.WithStack(err)
	}

	opts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithEndpoint(o.Config.Endpoint),
	}

	metricExporter, err := otlpmetrichttp.New(ctx, opts...)
	if err != nil {
		errors.WithStack(err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter,
			sdkmetric.WithInterval(o.Config.Interval))),
	)

	meter := meterProvider.Meter(o.Config.MeterProviderName)

	counter, err := meter.Int64Counter("counter",
		metric.WithDescription("Example of a Counter"),
	)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = meter.Float64ObservableGauge("gauge",
		metric.WithDescription("Example of a Gauge"),
		metric.WithFloat64Callback(func(_ context.Context, o metric.Float64Observer) error {
			o.Observe(rand.Float64() * 10)
			return nil
		}),
	)
	if err != nil {
		return errors.WithStack(err)
	}

	go o.produceMetrics(ctx, meterProvider, counter)

	return nil
}

func (o OtelMetrics) produceMetrics(ctx context.Context, meterProvider *sdkmetric.MeterProvider, counter metric.Int64Counter) {
	ticker := time.NewTicker(o.Config.Interval)
loop:
	for {
		select {
		case <-ctx.Done():
			contextTimeout, cancel := context.WithTimeout(context.Background(), o.GracefulShutdown)
			defer cancel()
			meterProvider.Shutdown(contextTimeout)
			break loop
		case <-ticker.C:
			counter.Add(ctx, 1)
		}
	}
}

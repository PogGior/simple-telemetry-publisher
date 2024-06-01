package trace

import (
	"simple-telemetry-publisher/internal/model"
	"time"

	"context"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

type OtelTrace struct {
	GracefulShutdown time.Duration
	Config         model.TraceConfig
}

func (o OtelTrace) Start(ctx context.Context) error {

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

	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(o.Config.Endpoint),
		otlptracehttp.WithInsecure(),
	}

	traceExporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return errors.WithStack(err)
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)

	traceProvider.Shutdown(ctx)

	go o.produceTrace(ctx, traceProvider)

	return nil
}

func (o OtelTrace) produceTrace(ctx context.Context, traceProvider *sdktrace.TracerProvider) {
	tracer := traceProvider.Tracer(o.Config.TracerName)
	ticker := time.NewTicker(o.Config.Interval)
loop:
	for {
		select {
		case <-ctx.Done():
			timeoutContext, cancel := context.WithTimeout(context.Background(), o.GracefulShutdown)
			defer cancel()
			traceProvider.Shutdown(timeoutContext)
			break loop
		case <-ticker.C:
			_, span := tracer.Start(ctx, "foo")
			span.End()
		}
	}
}

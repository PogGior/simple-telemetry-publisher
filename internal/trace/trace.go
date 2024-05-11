package trace

import (
	"simple-telemetry-publisher/internal/model"

	"context"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type OtelTrace struct {
	Config model.TraceConfig
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
	}

	traceExporter, err := otlptracehttp.New(ctx,opts...)
	if err != nil {
		return errors.WithStack(err)
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)

	tracer := traceProvider.Tracer(o.Config.TracerName)

	_, span := tracer.Start(ctx, "foo")
	span.End()

	return nil
}



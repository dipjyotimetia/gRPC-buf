package logz

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const (
	otelEndpoint    = "otel-collector:4317"
	initialInterval = 2 * time.Second
	maxInterval     = 5 * time.Second
	maxElapsedTime  = 10 * time.Second
	shutdownTimeout = 2 * time.Second
)

func newResource() *resource.Resource {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("grpc-buf"),
			semconv.ServiceVersion("1.0.0"),
			semconv.DeploymentEnvironment(os.Getenv("ENVIRONMENT")),
			attribute.Bool("sampling.exempt", true),
		),
	)
	if err != nil {
		slog.Error("Failed to create resource", "error", err)
		return resource.Default()
	}
	return r
}

// StartTracer starts the OpenTelemetry tracer with the OTLP exporter.
func StartTracer() (func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	expo, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(otelEndpoint),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
			Enabled:         true,
			InitialInterval: initialInterval,
			MaxInterval:     maxInterval,
			MaxElapsedTime:  maxElapsedTime,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(expo),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(newResource()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := tp.Shutdown(ctx); err != nil {
			slog.Error("Error shutting down tracer provider", "error", err)
		}
	}, nil
}

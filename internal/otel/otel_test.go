package otel

import (
	"context"
	"testing"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.17.0"
	"time"
)

func TestInitOtelProviders(t *testing.T) {
	ctx := context.Background()

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("service-go-kratos-sample"),
		semconv.ServiceVersion("v1.0.0"),
	)

	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithInsecure())
	if err != nil {
		t.Errorf("Error creating metric exporter: %v", err)
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter, metric.WithInterval(time.Second*1))),
		metric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		t.Errorf("Error creating trace exporter: %v", err)
	}

	idg := xray.NewIDGenerator()
	tracerProvider := trace.NewTracerProvider(
		trace.WithIDGenerator(idg),
		trace.WithBatcher(traceExporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)

	// Check that the MeterProvider and TracerProvider are set.
	if otel.GetMeterProvider() != meterProvider {
		t.Errorf("MeterProvider not set correctly")
	}
	if otel.GetTracerProvider() != tracerProvider {
		t.Errorf("TracerProvider not set correctly")
	}
}

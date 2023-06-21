package otel

import (
	"context"
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

func InitOtelProviders() {
	ctx := context.Background()

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("service-go-kratos-sample"),
		semconv.ServiceVersion("v1.0.0"),
	)

	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter, metric.WithInterval(time.Second*1))),
		metric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)
	// The MeterProvider is configured and registered globally. You can now run
	// your code instrumented with the OpenTelemetry API that uses the global
	// MeterProvider without having to pass this MeterProvider instance. Or,
	// you can pass this instance directly to your instrumented code if it
	// accepts a MeterProvider instance.

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	idg := xray.NewIDGenerator()
	tracerProvider := trace.NewTracerProvider(
		trace.WithIDGenerator(idg),
		trace.WithBatcher(traceExporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)
}

// Package telemetry configures OpenTelemetry tracing from the standard OTEL_* environment variables
package telemetry

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"os"

	"go.opentelemetry.io/contrib/propagators/autoprop"
	"go.opentelemetry.io/contrib/propagators/envcar"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

// exportEnabled reports whether the user has asked for spans to be exported.
// The spec defaults to OTLP when OTEL_TRACES_EXPORTER is unset, which for a CLI would mean
// connection errors against a collector that isn't there, so an explicit opt-in is required.
func exportEnabled() bool {
	if os.Getenv("OTEL_SDK_DISABLED") == "true" {
		return false
	}
	for _, k := range []string{"OTEL_TRACES_EXPORTER", "OTEL_EXPORTER_OTLP_ENDPOINT", "OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"} {
		if os.Getenv(k) != "" {
			return true
		}
	}
	return false
}

// newSpanExporter picks the exporter from OTEL_TRACES_EXPORTER and OTEL_EXPORTER_OTLP_[TRACES_]PROTOCOL,
// a nil exporter means "none". Only traces are supported, which keeps metric and log exporters out of the binary.
func newSpanExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	switch exporter := os.Getenv("OTEL_TRACES_EXPORTER"); exporter {
	case "", "otlp":
	case "none":
		return nil, nil
	case "console":
		// stderr so spans don't interleave with reports on stdout
		return stdouttrace.New(stdouttrace.WithWriter(os.Stderr))
	default:
		return nil, fmt.Errorf("unsupported OTEL_TRACES_EXPORTER %q", exporter)
	}

	switch protocol := cmp.Or(os.Getenv("OTEL_EXPORTER_OTLP_TRACES_PROTOCOL"), os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL"), "http/protobuf"); protocol {
	case "http/protobuf":
		return otlptracehttp.New(ctx)
	case "grpc":
		return otlptracegrpc.New(ctx)
	default:
		return nil, fmt.Errorf("unsupported OTLP protocol %q", protocol)
	}
}

// Setup installs the global propagator and, if exporting is configured, a tracer provider.
// The returned context carries any parent trace context found in TRACEPARENT/TRACESTATE/BAGGAGE.
// Without an exporter the parent context is still propagated through to the linters.
func Setup(ctx context.Context, version string) (context.Context, func(context.Context) error, error) {
	prop := autoprop.NewTextMapPropagator()
	otel.SetTextMapPropagator(prop)
	ctx = prop.Extract(ctx, &envcar.Carrier{})

	noop := func(context.Context) error { return nil }
	if !exportEnabled() {
		return ctx, noop, nil
	}

	exp, err := newSpanExporter(ctx)
	if err != nil || exp == nil {
		return ctx, noop, err
	}

	// later options take precedence, so OTEL_SERVICE_NAME / OTEL_RESOURCE_ATTRIBUTES override our defaults
	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName("crie"), semconv.ServiceVersion(version)),
		resource.WithTelemetrySDK(),
		resource.WithFromEnv(),
	)
	if err != nil {
		return ctx, noop, errors.Join(err, exp.Shutdown(ctx))
	}

	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exp), sdktrace.WithResource(res))
	otel.SetTracerProvider(tp)
	return ctx, tp.Shutdown, nil
}

package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func TestExportEnabled(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want bool
	}{
		{"nothing set", nil, false},
		{"exporter set", map[string]string{"OTEL_TRACES_EXPORTER": "console"}, true},
		{"endpoint set", map[string]string{"OTEL_EXPORTER_OTLP_ENDPOINT": "http://localhost:4318"}, true},
		{"traces endpoint set", map[string]string{"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT": "http://localhost:4318"}, true},
		{"sdk disabled wins", map[string]string{"OTEL_TRACES_EXPORTER": "console", "OTEL_SDK_DISABLED": "true"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, k := range []string{"OTEL_TRACES_EXPORTER", "OTEL_EXPORTER_OTLP_ENDPOINT", "OTEL_EXPORTER_OTLP_TRACES_ENDPOINT", "OTEL_SDK_DISABLED"} {
				t.Setenv(k, tt.env[k])
			}
			assert.Equal(t, tt.want, exportEnabled())
		})
	}
}

func TestNewSpanExporter(t *testing.T) {
	tests := []struct {
		name     string
		exporter string
		protocol string
		wantNil  bool
		wantErr  bool
	}{
		{name: "default is otlp http", wantNil: false},
		{name: "otlp grpc", exporter: "otlp", protocol: "grpc"},
		{name: "console", exporter: "console"},
		{name: "none", exporter: "none", wantNil: true},
		{name: "unknown exporter", exporter: "zipkin", wantNil: true, wantErr: true},
		{name: "unknown protocol", exporter: "otlp", protocol: "http/json", wantNil: true, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("OTEL_TRACES_EXPORTER", tt.exporter)
			t.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", tt.protocol)
			t.Setenv("OTEL_EXPORTER_OTLP_TRACES_PROTOCOL", "")

			exp, err := newSpanExporter(t.Context())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantNil, exp == nil)
			if exp != nil {
				assert.NoError(t, exp.Shutdown(t.Context()))
			}
		})
	}
}

func TestSetup_ExtractsParentFromEnv(t *testing.T) {
	t.Setenv("OTEL_TRACES_EXPORTER", "")
	t.Setenv("TRACEPARENT", "00-0102030405060708090a0b0c0d0e0f10-0102030405060708-01")

	ctx, shutdown, err := Setup(t.Context(), "test")
	require.NoError(t, err)
	defer func() { assert.NoError(t, shutdown(t.Context())) }()

	sc := trace.SpanContextFromContext(ctx)
	assert.Equal(t, "0102030405060708090a0b0c0d0e0f10", sc.TraceID().String())
	assert.True(t, sc.IsRemote())
}

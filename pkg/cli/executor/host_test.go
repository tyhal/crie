package executor

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func TestWillHost(t *testing.T) {
	// Only verify error path to avoid environment-dependent binary checks
	assert.Error(t, WillHost("definitely-not-a-real-binary-xyz"))
}

func TestHostExecutor_Setup(t *testing.T) {
	e := &hostExecutor{}
	assert.NoError(t, e.Setup(t.Context(), Instance{}))
}

func TestHostExecutor_Run(t *testing.T) {
	if WillHost("sh") != nil {
		t.Skip("neither 'sh' not available on host; skipping")
	}
	testHelperExecutor(t, NewHost)
}

func TestHostExecutor_Cleanup(t *testing.T) {
	e := &hostExecutor{}
	assert.NoError(t, e.Cleanup(t.Context()))
}

func TestHostExecutor_PropagatesTraceContext(t *testing.T) {
	if WillHost("printenv") != nil {
		t.Skip("'printenv' not available on host; skipping")
	}
	otel.SetTextMapPropagator(propagation.TraceContext{})
	t.Setenv("TRACEPARENT", "00-ffffffffffffffffffffffffffffffff-ffffffffffffffff-01")

	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
		SpanID:     trace.SpanID{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		TraceFlags: trace.FlagsSampled,
	})
	ctx := trace.ContextWithSpanContext(t.Context(), sc)

	e := NewHost()
	// setup has no span, the per-run context passed to Exec is what must be propagated
	assert.NoError(t, e.Setup(t.Context(), Instance{Bin: "printenv", Start: Par{"TRACEPARENT"}, NoFileArg: true}))
	defer func() { _ = e.Cleanup(t.Context()) }()

	var stdout, stderr bytes.Buffer
	assert.NoError(t, e.Exec(ctx, &stdout, &stderr, "."))
	// the span context from ctx must win over the TRACEPARENT crie itself inherited
	assert.Equal(t, "00-0102030405060708090a0b0c0d0e0f10-0102030405060708-01", strings.TrimSpace(stdout.String()))
}

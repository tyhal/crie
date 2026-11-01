package cli

import (
	"context"

	"charm.land/fang/v2"
	log "charm.land/log/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"github.com/tyhal/crie/internal/telemetry"
)

var tracer = otel.Tracer("github.com/tyhal/crie/internal/cli")

// Execute runs the root command inside a root span, with tracing configured from the OTEL_* env
func Execute(ctx context.Context, version string) error {
	ctx, shutdown, err := telemetry.Setup(ctx, version)
	if err != nil {
		log.Warn("telemetry disabled", "err", err)
	}
	defer func() {
		if err := shutdown(context.WithoutCancel(ctx)); err != nil {
			log.Warn("flushing telemetry", "err", err)
		}
	}()

	// renamed to the full command path once cobra has parsed it, see RootCmd.PersistentPreRunE
	ctx, span := tracer.Start(ctx, RootCmd.Name())
	defer span.End()

	err = fang.Execute(ctx, RootCmd, fang.WithVersion(version))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}

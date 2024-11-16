package logging

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/scope3-dio/src/common"
)

func Fatal(ctx context.Context, err error, data Data, m string) {
	log.Fatal().
		Str(common.CtxKeyTraceID, ctx.Value(common.CtxKeyTraceID).(string)).
		Interface("data", data).
		Err(err).
		Msg(m)
}

func Info(ctx context.Context, data Data, m string) {
	log.Info().
		Str(common.CtxKeyTraceID, ctx.Value(common.CtxKeyTraceID).(string)).
		Interface("data", data).
		Msg(m)
}

func Error(ctx context.Context, err error, data Data, m string) {
	log.Error().
		Str(common.CtxKeyTraceID, ctx.Value(common.CtxKeyTraceID).(string)).
		Interface("data", data).
		Err(err).
		Msg(m)
}

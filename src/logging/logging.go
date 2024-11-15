package logging

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/scope3-dio/common"
)

// Data represents logging data in structured logging.
type Data map[string]interface{}

func Fatal(ctx context.Context, err error, data Data, m string) {
	log.Fatal().
		Str(common.TraceIdKey, ctx.Value(common.TraceIdKey).(string)).
		Interface("data", data).
		Err(err).
		Msg(m)
}

func Info(ctx context.Context, data Data, m string) {
	log.Info().
		Str(common.TraceIdKey, ctx.Value(common.TraceIdKey).(string)).
		Interface("data", data).
		Msg(m)
}

func Error(ctx context.Context, err error, data Data, m string) {
	log.Error().
		Str(common.TraceIdKey, ctx.Value(common.TraceIdKey).(string)).
		Interface("data", data).
		Err(err).
		Msg(m)
}

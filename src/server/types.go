package server

import (
	"context"

	"github.com/scope3-dio/common"
)

type StorageClient interface {
	Get(context.Context, []common.PropertyQuery) map[string]string
}

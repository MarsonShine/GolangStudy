package fzlog

import (
	"context"

	"go.uber.org/zap"
)

func WithContext(ctx *context.Context) FzLog {
	logger := zap.L()
	return FzLog{
		logger,
		ctx,
	}
}

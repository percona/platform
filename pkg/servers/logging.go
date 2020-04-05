// Package servers provides common servers starting code for all SaaS components.
package servers

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/percona-platform/platform/pkg/logger"
)

// getCtxForRequest returns derived context with request-scoped logger set, and the logger itself.
func getCtxForRequest(ctx context.Context) (context.Context, *zap.Logger) {
	// UUID version 1: first 8 characters are time-based and lexicography sorted,
	// which is a useful property there
	u, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	l := zap.L().With(zap.String("request", u.String()))
	return logger.GetCtxWithLogger(ctx, l), l
}

// Package servers provides common servers starting code for all SaaS components.
package servers

import (
	"context"

	"go.uber.org/zap"

	"github.com/percona/platform/pkg/logger"
)

// getCtxForRequest returns derived context with request-scoped logger set, and the logger itself.
func getCtxForRequest(ctx context.Context) context.Context {
	return logger.GetContextWithLogger(ctx, zap.L())
}

// Package logger contains logging utilities.
package logger

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// key is unexported to prevent collisions - it is different from any other type in other packages
var key = struct{}{}

func newLogger(requestID string) *zap.Logger {
	return zap.L().With(zap.String("request", requestID))
}

// Get returns logger for given context. Set must be called before this method is called.
func Get(ctx context.Context) *zap.Logger {
	v := ctx.Value(key)
	if v == nil {
		l := newLogger("")
		l.DPanic("context logger not set")
		return l
	}
	return v.(*zap.Logger)
}

// Set returns derived context with set logger with given request ID.
func Set(ctx context.Context, requestID string) context.Context {
	return SetEntry(ctx, newLogger(requestID))
}

// SetEntry returns derived context with set given logger.
func SetEntry(ctx context.Context, l *zap.Logger) context.Context {
	if ctx.Value(key) != nil {
		Get(ctx).DPanic("context logger already set")
		return ctx
	}

	return context.WithValue(ctx, key, l)
}

// MakeRequestID returns a new request ID.
func MakeRequestID() string {
	// UUID version 1: first 8 characters are time-based and lexicography sorted,
	// which is a useful property there
	u, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	return u.String()
}

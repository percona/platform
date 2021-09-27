// Package tracing provides common request tracing utilities for all SaaS components.
package tracing

import (
	"context"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	tracingHeaderPrefix = "X-B3-"
	tracingHeaderName   = "X-B3-TraceId"
)

// OpenTracingHeadersMatcher preserves the OpenTracing headers added by Traefik
// after the HTTP request is received by grpc-gateway and are forwarded as-is
// to the grpc server.
// NOTE: key parameter must be in a Canonical format.
func OpenTracingHeadersMatcher(key string) bool {
	return strings.HasPrefix(key, tracingHeaderPrefix)
}

// GetRequestIDFromGrpcIncomingContext extracts from trace-id value from gRPC incoming metadata.
func GetRequestIDFromGrpcIncomingContext(ctx context.Context) string {
	if headers, ok := metadata.FromIncomingContext(ctx); ok {
		if reqIDs := headers.Get(tracingHeaderName); len(reqIDs) != 0 {
			return reqIDs[0]
		}
	}
	return ""
}

// GetRequestIDFromHTTPRequest extracts from trace-id value from incoming HTTP request.
func GetRequestIDFromHTTPRequest(r *http.Request) string {
	if reqID := r.Header.Get(tracingHeaderName); len(reqID) > 0 {
		return reqID
	}
	return ""
}

// Package tracing provides common request tracing utilities for all SaaS components.
package tracing

import (
	"context"
	"net/http"
	"net/textproto"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	tracingHeaderPrefix = "X-B3-"
	tracingHeaderName   = "X-B3-TraceId"
	portalRequestID     = "X-Percona-Portal-Request-Id"
)

// OpenTracingHeadersMatcher preserves the OpenTracing headers added by Traefik
// after the HTTP request is received by grpc-gateway and are forwarded as-is
// to the grpc server.
// NOTE: key parameter must be in a Canonical format.
func OpenTracingHeadersMatcher(key string) bool {
	return strings.HasPrefix(key, tracingHeaderPrefix)
}

// GetRequestIDFromGRPCIncomingContext extracts from trace-id value from gRPC incoming metadata.
func GetRequestIDFromGRPCIncomingContext(ctx context.Context) string {
	if headers, ok := metadata.FromIncomingContext(ctx); ok {
		if reqIDs := headers.Get(tracingHeaderName); len(reqIDs) != 0 {
			return reqIDs[0]
		}
	}
	return ""
}

// GetRequestIDFromHTTPRequest extracts from trace-id value from incoming HTTP request.
func GetRequestIDFromHTTPRequest(r *http.Request) string {
	if reqID := r.Header.Get(tracingHeaderName); len(reqID) > 0 { //nolint:canonicalheader
		return reqID
	}
	return ""
}

// AddRequestIDToGRPCResponseContext adds trace-id value to gRPC response metadata.
func AddRequestIDToGRPCResponseContext(ctx context.Context, reqID string) {
	_ = grpc.SetHeader(ctx, metadata.Pairs(textproto.CanonicalMIMEHeaderKey(portalRequestID), reqID))
}

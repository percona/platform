// Package grpcutils contains gRPC helpers.
package grpcutils

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Dial creates gRPC connection for a given address.
func Dial(ctx context.Context, address, userAgent string) (grpc.ClientConnInterface, error) {
	opts := []grpc.DialOption{
		grpc.WithBlock(), //nolint:staticcheck
		grpc.WithUserAgent(userAgent),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.DialContext(ctx, address, opts...) //nolint:staticcheck
	if err != nil {
		return nil, errors.Wrapf(err, "failed to dial %s", address)
	}

	return conn, nil
}

// NonblockingDial creates a non-blocking gRPC connection for a given address.
func NonblockingDial(_ context.Context, address, userAgent string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithUserAgent(userAgent),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to dial %s", address)
	}

	return conn, nil
}

package logger

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestHealthcheckWarning(t *testing.T) {
	var buf bytes.Buffer
	zc := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(&buf),
		zap.DebugLevel,
	)
	l := &GRPC{
		SugaredLogger: zap.New(zc, zap.ErrorOutput(zapcore.AddSync(&buf))).Sugar(),
	}

	l.Warningf("grpc: Server.Serve failed to complete security handshake from %q: %v", "10.0.0.115:48912", io.EOF)
	require.NoError(t, l.Sync())
	assert.Empty(t, buf.String())

	l.Warningf("grpc: Server.Serve failed to complete security handshake from lala")
	require.NoError(t, l.Sync())
	assert.Contains(t, buf.String(), "grpc: Server.Serve failed to complete security handshake from lala")
}

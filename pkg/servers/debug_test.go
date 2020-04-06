package servers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDebugServer(t *testing.T) {
	t.Run("InvalidAddr", func(t *testing.T) {
		require.PanicsWithValue(t, "Unexpected server stop: listen tcp: address 99999: invalid port.", func() {
			RunDebugServer(context.Background(), &RunDebugServerOpts{
				Addr: "invalid.port:99999",
			})
		})
	})
}

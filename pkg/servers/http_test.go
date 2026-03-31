package servers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPServer(t *testing.T) {
	t.Parallel()

	t.Run("InvalidAddr", func(t *testing.T) {
		t.Parallel()

		require.PanicsWithValue(t, "Unexpected server stop: listen tcp: address 99999: invalid port.", func() {
			RunHTTPServer(t.Context(), &RunHTTPServerOpts{
				Addr: "invalid.port:99999",
			})
		})
	})
}

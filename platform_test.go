package platform

import (
	"testing"

	"github.com/stretchr/testify/assert"

	authv1 "github.com/percona-platform/platform/gen/auth"
)

func TestTestify(t *testing.T) {
	assert.Equal(t, &authv1.SignUpResponse{}, &authv1.SignUpResponse{})
}

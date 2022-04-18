package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateTier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		tier  Tier
		error string
	}{
		{
			name:  "anonymous",
			tier:  Anonymous,
			error: "",
		}, {
			name:  "registered",
			tier:  Registered,
			error: "",
		}, {
			name:  "paid",
			tier:  Paid,
			error: "",
		}, {
			name:  "unknown",
			tier:  Tier("unknown"),
			error: "unknown check tier: \"unknown\"",
		}, {
			name:  "empty",
			tier:  Tier(""),
			error: "tier is empty",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.tier.Validate()

			if test.error != "" {
				assert.EqualError(t, err, test.error)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestValidateTiers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		tiers []Tier
		error string
	}{
		{
			name:  "normal",
			tiers: []Tier{Anonymous, Registered, Paid},
			error: "",
		}, {
			name:  "invalid",
			tiers: []Tier{Anonymous, Registered, Tier("invalid")},
			error: "unknown check tier: \"invalid\"",
		}, {
			name:  "empty tier",
			tiers: []Tier{Anonymous, Tier("")},
			error: "tier is empty",
		}, {
			name:  "empty tiers array",
			tiers: []Tier{},
			error: "",
		}, {
			name:  "duplicate",
			tiers: []Tier{Anonymous, Anonymous, Registered},
			error: "duplicate tier: \"anonymous\"",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateTiers(test.tiers)

			if test.error != "" {
				assert.EqualError(t, err, test.error)
				return
			}

			assert.NoError(t, err)
		})
	}
}

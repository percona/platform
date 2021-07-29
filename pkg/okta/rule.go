package okta

import (
	"time"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

// This file contains some modified Okta Sign On Policy structures from Okta SDK. Some changes are made:
// removed 'omitempty' annotations from fields which causes Okta API validation errors with empty values,
// some fields names changed to match golang style. Structures from okta package reused when possible.

// signOnPolicyRule is similar to okta.OktaSignOnPolicyRule, but with custom Actions filed type.
type signOnPolicyRule struct {
	Created     *time.Time                           `json:"created,omitempty"`
	ID          string                               `json:"id,omitempty"`
	LastUpdated *time.Time                           `json:"lastUpdated,omitempty"`
	Priority    int64                                `json:"priority,omitempty"`
	Status      string                               `json:"status,omitempty"`
	System      *bool                                `json:"system,omitempty"`
	Type        string                               `json:"type,omitempty"`
	Actions     *signOnPolicyRuleActions             `json:"actions,omitempty"`
	Conditions  *okta.OktaSignOnPolicyRuleConditions `json:"conditions,omitempty"`
	Name        string                               `json:"name,omitempty"`
}

// signOnPolicyRuleActions is similar to okta.OktaSignOnPolicyRuleActions, but with custom SignOn field type.
type signOnPolicyRuleActions struct {
	SignOn *signOnPolicyRuleSignOnActions `json:"signon,omitempty"`
}

// signOnPolicyRuleSignOnActions is similar to okta.OktaSignOnPolicyRuleSignonActions, but with custom Session filed type.
type signOnPolicyRuleSignOnActions struct {
	Access                  string                                `json:"access,omitempty"`
	FactorLifetime          int64                                 `json:"factorLifetime,omitempty"`
	FactorPromptMode        string                                `json:"factorPromptMode,omitempty"`
	RememberDeviceByDefault *bool                                 `json:"rememberDeviceByDefault,omitempty"`
	RequireFactor           *bool                                 `json:"requireFactor,omitempty"`
	Session                 *signOnPolicyRuleSignOnSessionActions `json:"session,omitempty"`
}

// signOnPolicyRuleSignOnSessionActions is similar okta.OktaSignOnPolicyRuleSignonSessionActions,
// but with modifications of json annotations.
type signOnPolicyRuleSignOnSessionActions struct {
	MaxSessionIdleMinutes     int64 `json:"maxSessionIdleMinutes,omitempty"`
	MaxSessionLifetimeMinutes int64 `json:"maxSessionLifetimeMinutes"` // Removed omitempty, filed should be present, but 0 is valid value
	UsePersistentCookie       *bool `json:"usePersistentCookie"`       // Removed omitempty
}

// check compatibility for types which has similar fields types (where such check is possible).
var _ = signOnPolicyRuleSignOnSessionActions(okta.OktaSignOnPolicyRuleSignonSessionActions{})

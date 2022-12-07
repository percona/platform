// Package okta implements methods for interacting with Okta API.
package okta

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/pkg/errors"
)

// Stat is helper struct for parsing Okta stats result.
type Stat struct {
	Count int    `json:"count"`
	Key   string `json:"key"`
}

// getUsersCount returns number of users based on filter.
func (c *Client) getUsersCount(ctx context.Context, since, until time.Time, filter string) (int, error) {
	params := url.Values{}
	params.Add("since", since.Format(time.RFC3339))
	params.Add("until", until.Format(time.RFC3339))
	params.Add("filter", filter)

	path := fmt.Sprintf("/sage/api/v1/logs/count?%s", params.Encode())
	var result []Stat

	err := c.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return 0, convertOktaError(oErr)
		}

		return 0, errors.Wrap(err, "failed to get log events")
	}
	return result[0].Count, nil
}

// GetCreatedUsersCount returns number of created users.
func (c *Client) GetCreatedUsersCount(ctx context.Context, since, until time.Time) (int, error) {
	return c.getUsersCount(ctx, since, until, "eventType eq \"user.lifecycle.create\"")
}

// GetActivatedUsersCount returns number of activated users.
func (c *Client) GetActivatedUsersCount(ctx context.Context, since, until time.Time) (int, error) {
	return c.getUsersCount(ctx, since, until, "eventType eq \"user.lifecycle.activate\"")
}

// GetDeactivatedUsersCount returns number of deactivated users.
func (c *Client) GetDeactivatedUsersCount(ctx context.Context, since, until time.Time) (int, error) {
	return c.getUsersCount(ctx, since, until, "eventType eq \"user.lifecycle.deactivate\"")
}

// getTopLoginAttempts returns top 'limit' user login attempts.
func (c *Client) getTopLoginAttempts(ctx context.Context, since, until time.Time, limit int, filter string, filterEmailDomains []string) ([]Stat, error) {
	var filterEmailDomainsString string
	if len(filterEmailDomains) != 0 {
		for i, s := range filterEmailDomains {
			if i != 0 {
				filterEmailDomainsString += " or "
			}
			filterEmailDomainsString += fmt.Sprintf("actor.alternateId ew \"%s\"", s)
		}
	}
	params := url.Values{}
	params.Add("since", since.Format(time.RFC3339))
	params.Add("until", until.Format(time.RFC3339))
	filterParam := fmt.Sprintf("eventType eq \"user.session.start\" and outcome.result eq \"%s\"", filter)
	if len(filterEmailDomainsString) > 0 {
		filterParam += fmt.Sprintf(" and not (%s)", filterEmailDomainsString)
	}
	params.Add("filter", filterParam)
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("field", "actor.alternateId")

	path := fmt.Sprintf("/sage/api/v1/logs/count?%s", params.Encode())
	var result []Stat

	err := c.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return result, convertOktaError(oErr)
		}

		return result, errors.Wrap(err, "failed to get log events")
	}
	return result, nil
}

// GetTopLoginSuccessfulAttempts returns top 'limit' user login successful attempts.
func (c *Client) GetTopLoginSuccessfulAttempts(ctx context.Context, since, until time.Time, limit int, filterEmailDomains []string) ([]Stat, error) {
	return c.getTopLoginAttempts(ctx, since, until, limit, "SUCCESS", filterEmailDomains)
}

// GetTopLoginFailedAttempts returns top 'limit' user login failed attempts.
func (c *Client) GetTopLoginFailedAttempts(ctx context.Context, since, until time.Time, limit int, filterEmailDomains []string) ([]Stat, error) {
	return c.getTopLoginAttempts(ctx, since, until, limit, "FAILURE", filterEmailDomains)
}

// getLoginTotalAttemptsCount returns total number of login attempts.
func (c *Client) getLoginTotalAttemptsCount(ctx context.Context, since, until time.Time, filter string, filterEmailDomains []string) (int, error) {
	var filterEmailDomainsString string
	if len(filterEmailDomains) != 0 {
		for i, s := range filterEmailDomains {
			if i != 0 {
				filterEmailDomainsString += " or "
			}
			filterEmailDomainsString += fmt.Sprintf("actor.alternateId ew \"%s\"", s)
		}
	}
	params := url.Values{}
	params.Add("since", since.Format(time.RFC3339))
	params.Add("until", until.Format(time.RFC3339))
	filterParam := fmt.Sprintf("eventType eq \"user.session.start\" and outcome.result eq \"%s\"", filter)
	if len(filterEmailDomainsString) > 0 {
		filterParam += fmt.Sprintf(" and not (%s)", filterEmailDomainsString)
	}
	params.Add("filter", filterParam)
	path := fmt.Sprintf("/sage/api/v1/logs/count?%s", params.Encode())
	var result []Stat

	err := c.DoRequest(ctx, "GET", path, nil, &result)
	if err != nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return 0, convertOktaError(oErr)
		}

		return 0, errors.Wrap(err, "failed to get log events")
	}
	return result[0].Count, nil
}

// GetLoginSuccessfulTotalAttemptsCount returns total number of login successful attempts.
func (c *Client) GetLoginSuccessfulTotalAttemptsCount(ctx context.Context, since, until time.Time, filterEmailDomains []string) (int, error) {
	return c.getLoginTotalAttemptsCount(ctx, since, until, "SUCCESS", filterEmailDomains)
}

// GetLoginFailedTotalAttemptsCount returns total number of login failed attempts.
func (c *Client) GetLoginFailedTotalAttemptsCount(ctx context.Context, since, until time.Time, filterEmailDomains []string) (int, error) {
	return c.getLoginTotalAttemptsCount(ctx, since, until, "FAILURE", filterEmailDomains)
}

// GetTotalUsersCount returns total number of users.
func (c *Client) GetTotalUsersCount(ctx context.Context) (int, error) {
	if len(c.oktaEveryoneGroupID) == 0 {
		qp := query.NewQueryParams(
			query.WithQ("Everyone"),
			query.WithFilter("type eq \"BUILT_IN\""),
		)
		groups, _, err := c.c.Group.ListGroups(ctx, qp)
		if err != nil {
			var oErr *okta.Error
			if errors.As(err, &oErr) {
				return 0, convertOktaError(oErr)
			}

			return 0, errors.Wrap(err, "failed to find everyone users group")
		}

		if len(groups) != 1 {
			return 0, fmt.Errorf("expect only one 'Everyone' group search result, got %d", len(groups))
		}
		c.oktaEveryoneGroupID = groups[0].Id
	}

	var groupStats struct {
		Count int `json:"usersCount"`
	}

	path := fmt.Sprintf("/api/v1/groups/%s/stats", c.oktaEveryoneGroupID)
	err := c.DoRequest(ctx, "GET", path, nil, &groupStats)
	if err != nil {
		var oErr *okta.Error
		if errors.As(err, &oErr) {
			return 0, convertOktaError(oErr)
		}

		return 0, errors.Wrap(err, "failed to get group stats")
	}

	return groupStats.Count, nil
}

/*
 * // Copyright (C) 2024 Percona LLC
 * //
 * // This program is free software: you can redistribute it and/or modify
 * // it under the terms of the GNU Affero General Public License as published by
 * // the Free Software Foundation, either version 3 of the License, or
 * // (at your option) any later version.
 * //
 * // This program is distributed in the hope that it will be useful,
 * // but WITHOUT ANY WARRANTY; without even the implied warranty of
 * // MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * // GNU Affero General Public License for more details.
 * //
 * // You should have received a copy of the GNU Affero General Public License
 * // along with this program. If not, see <https://www.gnu.org/licenses/>.
 */

package client

import (
	"bytes"
	"context"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"

	genericv1 "github.com/percona-platform/platform/gen/telemetry/generic"
)

// SendTelemetry sends telemetry data to Percona Platform.
func (c *Client) SendTelemetry(ctx context.Context, accessToken string, report *genericv1.ReportRequest) error {
	const path = "/v1/telemetry/GenericReport"

	body, err := protojson.Marshal(report)
	if err != nil {
		return err
	}

	err = c.sendPostRequest(ctx, path, accessToken, bytes.NewReader(body), nil)
	if err != nil {
		return errors.Wrap(err, "failed to send telemetry data")
	}

	return nil
}

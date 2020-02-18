package telemetry

import (
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	telemetryv1beta1 "github.com/Percona-Platform/platform/gen/telemetry"
	pmmv1beta1 "github.com/Percona-Platform/platform/gen/telemetry/events/pmm"
)

func TestValidators(t *testing.T) {
	id := uuid.New()
	event := &pmmv1beta1.ServerUptimeEvent{
		Id:                 id[:],
		Version:            "1.2.3",
		UpDuration:         ptypes.DurationProto(42 * time.Second),
		DistributionMethod: pmmv1beta1.DistributionMethod_DOCKER,
	}
	err := event.Validate()
	require.NoError(t, err)

	eventB, err := proto.Marshal(event)
	require.NoError(t, err)

	id = uuid.New()
	req := &telemetryv1beta1.ReportRequest{
		Events: []*telemetryv1beta1.Event{{
			Id:   id[:],
			Time: ptypes.TimestampNow(),
			Event: &telemetryv1beta1.AnyEvent{
				TypeUrl: proto.MessageName(event),
				Binary:  eventB,
			},
		}},
	}
	err = req.Validate()
	require.NoError(t, err)
}

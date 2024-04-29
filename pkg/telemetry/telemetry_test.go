package telemetry

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pmmv1 "github.com/percona/platform/gen/telemetry/events/pmm"
	reporterv1 "github.com/percona/platform/gen/telemetry/reporter"
)

func TestValidators(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	event := &pmmv1.ServerUptimeEvent{
		Id:                 id[:],
		Version:            "1.2.3",
		UpDuration:         durationpb.New(42 * time.Second),
		DistributionMethod: pmmv1.DistributionMethod_DOCKER,
	}
	err := event.Validate()
	require.NoError(t, err)
	assert.NotContains(t, event.ProtoReflect().Descriptor().FullName(), "beta")

	eventB, err := proto.Marshal(event)
	require.NoError(t, err)

	id = uuid.New()
	req := &reporterv1.ReportRequest{
		Events: []*reporterv1.Event{{
			Id:   id[:],
			Time: timestamppb.Now(),
			Event: &reporterv1.AnyEvent{
				TypeUrl: string(event.ProtoReflect().Descriptor().FullName()),
				Binary:  eventB,
			},
		}},
	}
	err = req.Validate()
	require.NoError(t, err)
	assert.NotContains(t, req.ProtoReflect().Descriptor().FullName(), "beta")
}

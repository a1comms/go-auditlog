package audit

import (
	"context"
	"fmt"

	metadata "cloud.google.com/go/compute/metadata"
	cl "cloud.google.com/go/logging"
	viap "github.com/a1comms/go-middleware-validate-iap"
)

var logger *cl.Logger

func init() {
	projectID, err := metadata.ProjectID()
	if err != nil {
		panic(err)
	}

	logClient, err := cl.NewClient(context.Background(), projectID)
	if err != nil {
		panic(err)
	}

	logger = logClient.Logger(
		"audit.a1comms.com/audit_log",
	)
}

type AuditEvent struct {
	EventType      string      `json:"event_type"`
	ActingIdentity string      `json:"acting_identity"`
	DeviceID       string      `json:"device_id,omitempty"`
	Context        interface{} `json:"context"`
}

func LogEvent(ctx context.Context, eventType, actingIdentity string, context interface{}) error {
	claim, _ := viap.GetGoogleClaimFromContext(ctx)

	err := logger.LogSync(ctx, cl.Entry{
		Severity: cl.Info,
		Payload: &AuditEvent{
			EventType:      eventType,
			ActingIdentity: actingIdentity,
			DeviceID:       claim.DeviceID,
			Context:        context,
		},
	})
	if err != nil {
		return fmt.Errorf("Failed to write audit log: %s", err)
	}

	return nil
}

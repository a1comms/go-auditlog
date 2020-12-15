package audit

import (
	"context"

	metadata "cloud.google.com/go/compute/metadata"
	cl "cloud.google.com/go/logging"
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
	Context        interface{} `json:"context"`
}

func LogEvent(ctx context.Context, eventType, actingIdentity string, context interface{}) error {
	return logger.LogSync(ctx, cl.Entry{
		Severity: cl.Info,
		Payload: &AuditEvent{
			EventType:      eventType,
			ActingIdentity: actingIdentity,
			Context:        context,
		},
	})
}

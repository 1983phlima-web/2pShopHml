package contracts

import (
	"time"

	"github.com/google/uuid"
)

type EventEnvelope struct {
	ID            string         `json:"id"`
	Type          string         `json:"type"`
	Version       int            `json:"version"`
	OccurredAt    time.Time      `json:"occurred_at"`
	TenantID      string         `json:"tenant_id"`
	AggregateType string         `json:"aggregate_type"`
	AggregateID   string         `json:"aggregate_id"`
	TraceID       string         `json:"trace_id"`
	SpanID        string         `json:"span_id"`
	Payload       map[string]any `json:"payload"`
}

func NewEvent(eventType, aggregateType, aggregateID, tenantID, traceID, spanID string, payload map[string]any) *EventEnvelope {
	return &EventEnvelope{
		ID:            uuid.Must(uuid.NewV7()).String(),
		Type:          eventType,
		Version:       1,
		OccurredAt:    time.Now().UTC(),
		TenantID:      tenantID,
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		TraceID:       traceID,
		SpanID:        spanID,
		Payload:       payload,
	}
}

type Pagination struct {
	Cursor    string `json:"cursor,omitempty"`
	Limit     int    `json:"limit"`
	Direction string `json:"direction,omitempty"`
}

type PaginatedResponse struct {
	Data       []any  `json:"data"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
	Total      int64  `json:"total,omitempty"`
}

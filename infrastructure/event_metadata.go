package infrastructure

import "github.com/google/uuid"

type EventMetadata struct {
	CorrelationId CorrelationId
	CausationId   CausationId
	Position      int
}

func NewEventMetadata(correlationId, causationId uuid.UUID, position int) EventMetadata {
	return EventMetadata{
		CorrelationId: CorrelationId{Value: correlationId},
		CausationId: CausationId{Value: causationId},
		Position: position,
	}
}

func NewEventMetadataFrom(m CommandMetadata) EventMetadata {
	return EventMetadata{
		CorrelationId: m.CorrelationId,
		CausationId:   m.CausationId,
	}
}

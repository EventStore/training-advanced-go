package infrastructure

import "github.com/google/uuid"

type CorrelationId struct {
	Value uuid.UUID
}

type CausationId struct {
	Value uuid.UUID
}

type CommandMetadata struct {
	CorrelationId CorrelationId
	CausationId   CausationId
}

func NewCommandMetadata(correlationId, causationId uuid.UUID) CommandMetadata {
	return CommandMetadata{
		CorrelationId: CorrelationId{Value: correlationId},
		CausationId: CausationId{Value: causationId},
	}
}

func NewCommandMetadataFrom(m EventMetadata) CommandMetadata {
	return CommandMetadata{
		CorrelationId: m.CorrelationId,
		CausationId: m.CausationId,
	}
}

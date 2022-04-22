package projections

import (
	"github.com/EventStore/training-introduction-go/infrastructure"
)

type Subscription interface {
	Project(event interface{}, m infrastructure.EventMetadata)
}

type Projector struct {
	Subscription

	projection infrastructure.EventHandler
}

func NewProjector(p infrastructure.EventHandler) *Projector {
	return &Projector{
		projection: p,
	}
}

func (p Projector) Project(event interface{}, m infrastructure.EventMetadata) {
	t := infrastructure.GetValueType(event)
	if !p.projection.CanHandle(t) {
		return
	}

	p.projection.Handle(t, event, m)
}

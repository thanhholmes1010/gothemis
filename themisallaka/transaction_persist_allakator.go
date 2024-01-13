package themisallaka

import (
	"github.com/thaianhsoft/gothemis/themisgraphorm/schema"
	"github.com/thaianhsoft/gothemis/themisgraphorm/schema/edge"
	"github.com/thaianhsoft/gothemis/themisgraphorm/schema/field"
)

type PersistEvent struct {
	*schema.Schema
	EventId  uint32
	EntityId uint32
	Version  uint32
	Data     any
}

func (p *PersistEvent) DefineFields() []field.IField {
	return []field.IField{
		field.Integer(32).Name("EventId").Unsigned().Null(false),
		field.Integer(32).Name("EntityId").Unsigned().Null(false),
		field.Integer(32).Name("Version").Unsigned().Null(false),
		field.JSONField(map[string]any{}).Name("Data").Null(false),
	}
}

func (p *PersistEvent) DefineEdges() []edge.IEdge {
	return []edge.IEdge{}
}

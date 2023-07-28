package schema

import (
	"github.com/thaianhsoft/gothemis/themisgraphorm/schema/edge"
	"github.com/thaianhsoft/gothemis/themisgraphorm/schema/field"
)

type Migrator interface {
	DefineFields() []field.IField
	DefineEdges() []edge.IEdge
	migrate()
}

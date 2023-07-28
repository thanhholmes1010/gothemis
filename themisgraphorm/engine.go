package themisgraphorm

import (
	"github.com/thaianhsoft/gothemis/themisgraphorm/schema"
	"strings"
)

type ThemisGraphEngine struct {
	createStmtBuilder *strings.Builder
}

func (tge *ThemisGraphEngine) Migrate(schemaClasses ...schema.Migrator) string {
	for _, schema := range schemaClasses {
		schema.DefineFields()
		schema.DefineEdges()
	}
	return tge.createStmtBuilder.String()
}

package schema

import "strings"

type Schema struct {
	*strings.Builder
}

// one function implement signature for method Migrator interface
func (s *Schema) migrate() {}

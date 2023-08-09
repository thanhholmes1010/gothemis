package container

import (
	"testing"
)

type FuncCommand func()

func TestRecommendSearch(t *testing.T) {
	p := NewPrefixTrie(".", ":")
	p.Insert("gen.schema.:name.fields.:field_names:options", FuncCommand(func() {
	}))
	p.Insert("gen.migrate.:name_1.option.all", FuncCommand(func() {
	}))
	p.Insert("gen.migrate.:name_2.option", FuncCommand(func() {

	}))
	data := p.SearchAllPossibles("gen", p.DrawLevelWithTabSpace)
	switch fn := data.(type) {
	case FuncCommand:
		fn()
	}
}

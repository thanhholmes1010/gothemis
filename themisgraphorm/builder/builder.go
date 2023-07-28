package builder

import (
	"fmt"
	"strings"
)

type QuerierInterface interface {
	Query() (string, []any)
}
type Builder struct {
	*strings.Builder
}

func (b *Builder) writeString(v string) *Builder {
	if b.Builder == nil {
		b.Builder = &strings.Builder{}
	}
	b.Builder.WriteString(v)
	return b
}

func (b *Builder) quoteIdent(v string) *Builder {
	b.writeString("`")
	b.writeString(v)
	b.writeString("`")
	return b
}

func (b *Builder) space() *Builder {
	b.writeString(" ")
	return b
}

func (b *Builder) comma() *Builder {
	b.writeString(",")
	return b
}

func (b *Builder) dot() *Builder {
	b.writeString(".")
	return b
}

type TableSelector struct {
	name string
}

func Table(name string) *TableSelector {
	return &TableSelector{
		name: name,
	}
}

func (tb *TableSelector) C(col string) string {
	return fmt.Sprintf("`%v`.`%v`", tb.name, col)
}

func (tb *TableSelector) v() {} // implement view query

type Predicate struct {
	preds   []*Predicate
	mergeOp prefixOp
	*Builder
	val         any
	args        []any
	predicateOp OpPredicate
	col         string
}

type OpPredicate uint16

const (
	EQ OpPredicate = iota + 1
	NEQ
	IN
	INN
	LIKE
	GTE
	GT
	LTE
	LT
)

var mapPredicateOp = [...]string{"=", "!=", "IS NULL", "IS NOT NULL", "LIKE", ">=", ">", "<=", "<"}

func (p *Predicate) mergePredicate() {

}
func (p *Predicate) Query() (string, []any) {
	p.mergePredicate() // merge nested predicate into
	switch view := p.val.(type) {
	case *Selector:
		view_query, view_args := view.Query()
		p.writeString("(")
		p.writeString(view_query)
		p.args = append(p.args, view_args)
		p.writeString(")")
	case string:
		// only value normal
		opString := mapPredicateOp[p.predicateOp]
		p.space().writeString(opString)
		if p.predicateOp != IN && p.predicateOp != INN {
			p.writeString("?")
		}
	}
	return p.String(), p.args
}

type prefixOp uint8

const (
	orPrefix prefixOp = iota + 1
	andPrefix
)

func newP(preds ...*Predicate) *Predicate {
	return &Predicate{
		Builder: &Builder{},
		preds:   preds,
	}
}

func P(col string, op OpPredicate, val any) *Predicate {
	p := newP()
	p.col = col
	p.predicateOp = op
	p.val = val
	return p
}
func Or(preds ...*Predicate) *Predicate {
	p := newP(preds...)
	p.mergeOp = orPrefix
	return p
}

func And(preds ...*Predicate) *Predicate {
	p := newP(preds...)
	p.mergeOp = andPrefix
	return p
}

type Selector struct {
	cols     []string
	from     queryViewer
	where    *Predicate
	*Builder       // use to write query into this
	args     []any // this to call in database pass argument dynamic on sql.Query() from standard lib
}

func (s *Selector) Query() (string, []any) {
	s.writeString("SELECT")
	s.space() // add one space
	for i, col := range s.cols {
		if i > 0 {
			s.comma() // add one character "," into query builder
		}
		s.quoteIdent(col) // add col name into query builder with ident `` in it
	}
	s.space()
	s.writeString("FROM")
	switch view := s.from.(type) {
	case *Selector:
		view_query, view_args := view.Query() // query from it
		s.writeString("(")
		s.writeString(view_query)
		s.args = append(s.args, view_args...)
		s.writeString(")")
	case *TableSelector:
		s.quoteIdent(view.name) // ident with name table
	}
	s.space()
	s.writeString("WHERE")
	where_query, where_args := s.where.Query()
	s.writeString(where_query)
	s.args = append(s.args, where_args...)
	return s.String(), s.args
}
func (s *Selector) v() {} // implement view query

// why use this
// selector can select from sub-query or from table
type queryViewer interface {
	v()
}

func Select() *Selector {
	return &Selector{
		Builder: &Builder{},
	}
}
func (s *Selector) Select(cols ...string) *Selector {
	s.cols = make([]string, len(cols))
	for i, _ := range cols {
		s.cols[i] = cols[i]
	}
	return s
}

func (s *Selector) From(qv queryViewer) *Selector {
	s.from = qv
	return s
}

func (s *Selector) Where(p *Predicate) *Selector {
	s.where = p
	return s
}

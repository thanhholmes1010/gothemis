package mysql

import (
	"fmt"
	"reflect"
	"strings"
)

type QuerierInterface interface {
	Query() (string, []any)
}

func newBuilder() *Builder {
	return &Builder{
		builder: &strings.Builder{},
	}
}

type Builder struct {
	builder *strings.Builder
}

func (b *Builder) writeString(v string) *Builder {
	b.builder.WriteString(v)
	return b
}

func (b *Builder) quoteIdent(v string) *Builder {
	if string(v[0]) == "`" {
		b.writeString(v)
		return b
	}
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
	as   string
}

func Table(name string) *TableSelector {
	return &TableSelector{
		name: name,
	}
}

func (tb *TableSelector) C(col string) string {
	v := tb.name
	if tb.as != "" {
		v = tb.as
	}
	return fmt.Sprintf("`%v`.`%v`", v, col)
}

func (tb *TableSelector) v() {}

func (tb *TableSelector) As(s string) *TableSelector {
	tb.as = s
	return tb
} // implement view query

type Predicate struct {
	preds   []*Predicate
	mergeOp prefixOp
	*Builder
	val         any
	args        []any
	predicateOp OpPredicate
	col         string
	depth       int
}

type OpPredicate uint16

const (
	EQ OpPredicate = iota + 1
	NEQ
	NULL
	NOTNULL
	LIKE
	GTE
	GT
	LTE
	LT
	IN
	EXIST
)

var mapPredicateOp = [...]string{
	"",
	"=", "!=", "IS NULL",
	"IS NOT NULL", "LIKE", ">=",
	">", "<=", "<",
	"IN", "EXIST",
}

func (p *Predicate) mergePredicate() {
	for i, childPredicate := range p.preds {
		n := len(childPredicate.preds)
		if i > 0 {
			p.space()
			p.writeString(p.mergeOp.ToOpString())
		}
		query, args := childPredicate.Query()
		if query != "" {
			p.space()
			p.mergeWithParenthese(&query, n)
		}
		p.args = append(p.args, args...)
	}
}

func (p *Predicate) Query() (string, []any) {
	p.mergePredicate() // merge nested predicate int // o
	if p.mergeOp > 0 {
		return p.builder.String(), p.args
	}
	opString := mapPredicateOp[p.predicateOp]
	p.quoteIdent(p.col).space()
	p.writeString(opString).space()
	switch view := p.val.(type) {
	case *Selector:
		view_query, view_args := view.Query()
		p.writeString("(")
		p.writeString(view_query)
		p.args = append(p.args, view_args)
		p.writeString(")")
	default:
		// only value normal
		if p.mergeOp == 0 {
			//opString := mapPredicateOp[p.predicateOp]
			//p.quoteIdent(p.col).space()
			//p.writeString(opString).space()
			fmt.Println("flush here: ", p.builder.String())
			if p.predicateOp != NULL && p.predicateOp != NOTNULL {
				p.writeString("?")
			}
			if p.predicateOp == IN || p.predicateOp == EXIST {
				rv := reflect.Indirect(reflect.ValueOf(p.val))
				nestValArg := ""
				switch rv.Kind() {
				case reflect.Slice:
					for i := 0; i < rv.Len(); i++ {
						if i > 0 {
							nestValArg += ", "
						}
						valStr := fmt.Sprintf("%v", rv.Index(i).Interface())
						nestValArg += valStr
					}
				}
				p.val = nestValArg
			}
			p.args = append(p.args, p.val)
		}
	}
	return p.builder.String(), p.args
}

func (p *Predicate) mergeWithParenthese(s *string, lenChild int) {
	if lenChild > 1 {
		p.writeString("(")
	}
	p.writeString(*s)
	if lenChild > 1 {
		p.writeString(")")
	}
}

type prefixOp uint8

func (p prefixOp) ToOpString() string {
	if p == orPrefix {
		return "OR"
	}
	return "AND"
}

const (
	orPrefix prefixOp = iota + 1
	andPrefix
)

func newP(preds ...*Predicate) *Predicate {
	p := &Predicate{
		Builder: newBuilder(),
		preds:   preds,
	}
	return p
}

func P(col string, op OpPredicate, val any) *Predicate {
	p := newP()
	p.col = col
	p.predicateOp = op
	p.val = val
	return p
}
func Or(preds ...*Predicate) *Predicate {
	fmt.Println("call or")
	p := newP(preds...)
	p.mergeOp = orPrefix
	return p
}

func And(preds ...*Predicate) *Predicate {
	fmt.Println("call and")
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
	joins    []*joiner
}

func (s *Selector) Query() (string, []any) {
	s.writeString("SELECT")
	s.space() // add one space
	for i, col := range s.cols {
		if i > 0 {
			s.comma() // add one character "," into query builder
			s.space()
		}
		s.quoteIdent(col) // add col name into query builder with ident `` in it
	}
	if s.from != nil {
		s.space()
		s.writeString("FROM").space()
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
	}

	// build join if have
	if len(s.joins) > 0 {
		for _, joiner := range s.joins {
			s.space()
			s.writeString(joiner.op.toOpJoin())
			s.space()
			switch jv := joiner.joinTable.(type) {
			//case *Selector:
			//	jvquery, jvargs := jv.Query()
			//	s.writeString("(")
			//	s.writeString(jvquery)
			//	s.writeString(")")
			//	s.args = append(jvargs)
			case *TableSelector:
				s.writeString(jv.name)
				if jv.as != "" {
					s.writeString(" AS ")
					s.quoteIdent(jv.as)
				}
			}
			s.writeString(" ON ")
			s.quoteIdent(joiner.c1)
			s.writeString(" = ")
			s.quoteIdent(joiner.c2)
		}
	}
	if s.where != nil {
		s.space()
		s.writeString("WHERE")
		where_query, where_args := s.where.Query()
		s.writeString(where_query)
		s.args = append(s.args, where_args...)
	}
	return s.builder.String(), s.args
}
func (s *Selector) v() {} // implement view query

// why use this
// selector can select from sub-query or from table
type queryViewer interface {
	v()
}

func newSelector() *Selector {
	return &Selector{
		Builder: newBuilder(),
	}
}
func Select(cols ...string) *Selector {
	s := newSelector()
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
	fmt.Println("p here: ", len(p.preds))
	return s
}

func (s *Selector) joinSelect(js queryViewer, op joinOp) *Selector {
	s.joins = append(s.joins, &joiner{
		joinTable: js,
		op:        op,
	})
	return s
}

func (s *Selector) InnerJoin(js queryViewer) *Selector {
	return s.joinSelect(js, innerJoin)
}

func (s *Selector) LeftJoin(js queryViewer) *Selector {
	return s.joinSelect(js, leftJoin)
}

func (s *Selector) RightJoin(js queryViewer) *Selector {
	return s.joinSelect(js, rightJoin)
}

func (s *Selector) FullJoin(js queryViewer) *Selector {
	return s.joinSelect(js, fullJoin)
}

func (s *Selector) On(c1, c2 string) *Selector {
	lastJoiner := s.joins[len(s.joins)-1]
	lastJoiner.c1 = c1
	lastJoiner.c2 = c2
	return s
}

type joinOp uint8

func (o joinOp) toOpJoin() string {
	if o == innerJoin {
		return "INNER JOIN"
	}
	if o == leftJoin {
		return "LEFT JOIN"
	}
	if o == rightJoin {
		return "RIGHT JOIN"
	}
	return "FULL JOIN"
}

const (
	innerJoin joinOp = iota
	leftJoin
	rightJoin
	fullJoin
)

type joiner struct {
	joinTable queryViewer
	op        joinOp
	c1        string
	c2        string
}

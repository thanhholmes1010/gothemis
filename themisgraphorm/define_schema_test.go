package themisgraphorm

import (
	"fmt"
	"github.com/thaianhsoft/gothemis/themisgraphorm/schema"
	"github.com/thaianhsoft/gothemis/themisgraphorm/schema/edge"
	"github.com/thaianhsoft/gothemis/themisgraphorm/schema/field"
	"testing"
)

var (
	StudentType edge.Type = "student"
	BookType    edge.Type = "book"
)

type Student struct {
	*schema.Schema
}

func (u *Student) DefineFields() []field.IField {
	return []field.IField{
		field.Integer(32).Name("Id").Unsigned().AI(),
		field.Varchar(40).Name("Email").Unique().Default("thaianhsoft@gmail.com"),
	}
}

func (u *Student) DefineEdges() []edge.IEdge {
	return []edge.IEdge{
		edge.PointTo("HasBooks", BookType),
	}
}

type Book struct {
	*schema.Schema
}

func (b *Book) DefineFields() []field.IField {
	return []field.IField{
		field.Integer(32).Name("Id").Unsigned().AI(),
		field.Varchar(30).Name("Title").Unique().Null(false),
	}
}

func (b *Book) DefineEdges() []edge.IEdge {
	return []edge.IEdge{
		// define cause table add one column hold relation for book and student
		// column name OwnerRelId must is unique
		// but for Book one row is id  unique, don't create indexing on OwnerRelId
		edge.PointBack("Owner", StudentType).RefOn("HasBooks", StudentType).Unique(),
	}
}

func TestDefineSchema(t *testing.T) {
	engine := &ThemisGraphEngine{}
	createTableStmt := engine.Migrate(&Student{}, &Book{})
	fmt.Println(createTableStmt)
}

package schema

import (
	"fmt"
	"github.com/thaianhsoft/gothemis/themisgraphorm/schema/edge"
	"github.com/thaianhsoft/gothemis/themisgraphorm/schema/field"
	"testing"
)

type User struct {
	*Schema
	Id        uint32
	Age       uint8
	Address   string
	LastName  string
	FirstName string
	Email     string
}

func (u *User) DefineFields() []field.IField {
	return []field.IField{
		field.Integer(32).Name("Id").Unsigned().AI(),
		field.Integer(8).Name("Age").Unsigned().Null(false),
		field.Varchar(10).Name("FirstName").Null(false),
		field.Varchar(10).Name("LastName").Null(false),
		field.Varchar(45).Name("Address").Null(false),
		field.Varchar(45).Name("Email").Null(false),
	}
}

func (u *User) DefineEdges() []edge.IEdge {
	return []edge.IEdge{}
}

type UserMessage struct {
	UserAge       uint8
	UserLastName  string
	UserFirstName string
	UserAddress   string
	UserEmail     string
}

func TestChangeset(t *testing.T) {
	u := &User{}
	msg := &UserMessage{
		UserAge:       32,
		UserLastName:  "Le",
		UserFirstName: "Thai Anh",
		UserAddress:   "50/b, nguyen van luong, go vap",
		UserEmail:     "thaianhsoft@gmail.com",
	}
	cs := Cast(msg, u)
	err := cs.
		ValidateRequired().
		ValidPattern("Email", []byte(`.*@gmail.com`)).
		IsValidAll()
	if err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Printf("user entity after cast: %v\n", u)
	}
}

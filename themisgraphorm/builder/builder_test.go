package builder

import (
	"testing"
)

func TestBuilder(t *testing.T) {
	t1 := Table("Student")
	t2 := Table("Book")
	s1 := Select()
	s2 := Select()
	s1.
		Select("Age", "Name", "Id").
		From(
			s2.
				Select("Id", "BookTitle").
				From(t2).
				Where(P(t2.C("Id"), EQ, 1)),
		).
		Where(
			And(
				Or(
					P(t1.C("Age"), GTE, 15),
					P(t1.C("Age"), LTE, 32),
				),
				P(t1.C("Name"), LIKE, "%thaianhsoft%"),
			),
		)
}

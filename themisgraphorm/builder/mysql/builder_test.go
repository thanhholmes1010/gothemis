package mysql

import (
	"fmt"
	"testing"
)

func TestPredicate(t *testing.T) {
	t1 := Table("User")
	//t2 := Table("Friend")
	t3 := Table("Book")
	s1 := Select(
		t1.C("Id"), t1.C("Name"), t1.C("Age"),
		t3.C("Id"), t3.C("Title"),
	).
		From(t1).
		InnerJoin(t3).
		On(t1.C("Id"), t3.C("OwnerRelId")).
		Where(
			And(
				Or(
					P(t1.C("Age"), GTE, 15),
					P(t1.C("Age"), LTE, 35),
				),
				Or(
					P(t3.C("Title"), LIKE, "%VanHocNgheThuat%"),
					P(t3.C("Title"), LIKE, "%VanHocDanGian%"),
				),
			),
		)
	query, args := s1.Query()
	fmt.Println(query, args)
}

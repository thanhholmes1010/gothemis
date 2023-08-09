package container

type Pair struct {
	First any
	Right any
}

func NewPair(first, right any) *Pair {
	return &Pair{
		first,
		right,
	}
}

package container

type stackNode struct {
	down *stackNode
	data any
}

type stackIterator struct {
	top  *stackNode
	size int
}

func (s *stackIterator) Empty() bool {
	return s.top == nil
}

func (s *stackIterator) Size() int {
	return s.size
}
func (s *stackIterator) Pop() any {
	old := s.top
	s.top = s.top.down
	s.size--
	return old.data
}

func (s *stackIterator) Peak() any {
	return s.top.data
}

func (s *stackIterator) Push(isPushFront bool, data ...any) {
	for _, v := range data {
		n := &stackNode{data: v}
		if s.top == nil {
			s.top = n
		} else {
			s.top, n.down = n, s.top
		}
		s.size++
	}
}

func NewStack() IteratorContainer {
	return &stackIterator{}
}

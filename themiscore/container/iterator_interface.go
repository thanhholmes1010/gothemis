package container

type IteratorContainer interface {
	Empty() bool // check if has next element
	Pop() any    // get and pop
	Peak() any   // lookahead but don't pop before
	Push(isPushFront bool, data ...any)
	Size() int
}

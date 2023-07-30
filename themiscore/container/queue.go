package container

type queueNode struct {
	next *queueNode
	data any
}

type queueIterator struct {
	f *queueNode
	b *queueNode
}

func (q *queueIterator) Empty() bool {
	return q.f == nil
}

func (q *queueIterator) Pop() any {
	old := q.f
	q.f = old.next
	return old.data
}

func (q *queueIterator) Peak() any {
	return q.f.data
}

func (q *queueIterator) Push(isPushFront bool, data ...any) {
	for _, oneData := range data {
		v := &queueNode{data: oneData}
		if isPushFront {
			q.pushFront(v)
		} else {
			q.pushBack(v)
		}
	}
}

func (q *queueIterator) pushBack(v *queueNode) {
	if q.f == nil {
		q.f, q.b = v, v
	} else {
		q.b, q.b.next = v, v
	}
}

func (q *queueIterator) pushFront(v *queueNode) {
	if q.f == nil {
		q.f, q.b = v, v
	} else {
		q.f, v.next = v, q.f
	}
}

func NewQueue() IteratorContainer {
	return &queueIterator{}
}

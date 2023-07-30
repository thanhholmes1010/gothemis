package container

import (
	"github.com/thaianhsoft/gothemis/themisutils"
	"math"
	"sync/atomic"
)

type SplayNode struct {
	f    *SplayNode
	c    [2]*SplayNode
	val  uint64
	fre  int64
	data any
}

func newNode(key uint64, data any) *SplayNode {
	return &SplayNode{
		c:    [2]*SplayNode{},
		val:  key,
		fre:  0,
		data: data,
	}
}

func (n *SplayNode) GetChilds() [2]*SplayNode {
	return n.c
}

func (n *SplayNode) GetData() any {
	return n.data
}

func (n *SplayNode) UpdateData(update_func func(new_data any) any) {
	n.data = update_func(n.data)
}

type SplayTree struct {
	root        *SplayNode
	t0          int64
	replaceFunc func(f1, f2 int) bool
}

func (T *SplayTree) GetRoot() *SplayNode {
	return T.root
}

func (T *SplayTree) rotate(n *SplayNode) {
	v := themisutils.AssertAssignBit(n.f.c[0] == n)
	//fmt.Println(v)
	p, m := n.f, n.c[v]
	//fmt.Println("p: ", p, " m: ", m)
	if p.f != nil {
		pv := themisutils.AssertAssignBit(p.f.c[1] == p)
		p.f.c[pv] = n
	}
	n.f, n.c[v] = p.f, p
	p.f, p.c[v^1] = n, m
	if m != nil {
		m.f = p
	}
}

func (T *SplayTree) splay(n *SplayNode, s *SplayNode) {
	for n.f != s {
		m, l := n.f, n.f.f
		//fmt.Println("splay")
		if l == s {
			T.rotate(n)
		} else {
			if (m.c[0] == n) == (l.c[0] == m) {
				T.rotate(m)
				T.rotate(n)
			} else {
				T.rotate(n)
				T.rotate(n)
			}
		}
	}
	if s == nil {
		T.root = n
	}
}

func (T *SplayTree) Find(v uint64, sp ...bool) *SplayNode {
	n := T.root
	for n != nil {
		if n.val == v {
			atomic.AddInt64(&T.t0, 1)
			atomic.AddInt64(&n.fre, 1)
			//n.fre++
			break
		}
		if n.val > v {
			if n.c[0] == nil {
				break
			}
			n = n.c[0]
		} else {
			if n.c[1] == nil {
				break
			}
			n = n.c[1]
		}
	}
	is_splay := true
	if len(sp) == 1 && !sp[0] {
		is_splay = false
	}
	if is_splay {
		T.splay(n, nil)
	}
	return n
}

func (T *SplayTree) Insert(v uint64, data any) {
	if T.root == nil {
		T.root = newNode(v, data)
		return
	}
	n := T.Find(v, false)

	new_node := newNode(v, data)
	new_node.f = n
	if n.val < v {
		n.c[1] = new_node
	} else {
		n.c[0] = new_node
	}
	T.splay(new_node, nil)
}

func (T *SplayTree) Erase(v uint64) {
	n := T.Find(v, false)
	if n.val == v {
		T.splay(n, nil) // splay node find to root
		// then change root is right child, deallocate parent of it
		L, R := n.c[0], n.c[1]
		if L != nil {
			L.f = nil
		}
		if R != nil {
			R.f = nil
		}
		n = L
		if L == nil || R == nil {
			if R != nil {
				n = R
			}
		} else {
			for n.c[1] != nil {
				n = n.c[1]
			}
			T.splay(n, nil)
			n.c[1] = R
			R.f = n
		}
		T.root = n
	}
}

func (t *SplayTree) ExtractOrder() *[]uint64 {
	order := []uint64{}
	t.walk_down(t.root, &order)
	return &order
}

func (t *SplayTree) walk_down(n *SplayNode, order *[]uint64) {
	if n == nil {
		return
	}
	*order = append(*order, n.val)
	t.walk_down(n.c[0], order)
	t.walk_down(n.c[1], order)
}

func (t *SplayTree) ReplaceLeastFreq() {
	n := t.root
	var lf int64 = math.MaxInt64
	for n != nil {
		v := lf
		if n.c[0] != nil {
			if n.c[0].fre < v {
				n = n.c[0]
				v = n.c[0].fre
			}
		}
		if n.c[1] != nil {
			if n.c[1].fre < v {
				n = n.c[1]
				v = n.c[1].fre
			}
		}
	}
}

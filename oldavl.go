// Package avl implements an AVL balanced binary tree.
package oldavl

// Tree holds elements of the AVL tree.
type Tree struct {
	root *Node
	size int
	cmp  func(interface{}, interface{}) int
}

// A Node holds an Ordered element of the AVL tree in
// the Val field.
type Node struct {
	Val interface{}
	c   [2]*Node
	p   *Node
	b   int8
}

func NewTree(cmp func(interface{}, interface{}) int) *Tree {
	return &Tree{cmp: cmp}
}

// Size returns the number of elements stored in the tree.
func (t *Tree) Size() int {
	return t.size
}

// Lookup looks up val and returns the matching element if
// it is found.
//
// Val's Less implementation must be able to handle
// comparisons to elements stored in this tree.
func (t *Tree) Lookup(val interface{}) (match interface{}, ok bool) {
	if t == nil {
		return
	}
	n := t.root
	for n != nil {
		cmp := t.cmp(val, n.Val)
		switch {
		case cmp < 0:
			n = n.c[0]
		default:
			return n.Val, true
		case cmp > 0:
			n = n.c[1]
		}
	}
	return
}

// Insert looks up val and inserts it into the tree.
// If a matching element is found in the tree then the
// found element delval is removed from the tree and returned.
//
// Val's Less implementation must be able to handle
// comparisons to elements stored in this tree.
func (t *Tree) Insert(val interface{}) {
	t.insert(val, nil, &t.root)
}

func (t *Tree) insert(val interface{}, p *Node, qp **Node) bool {
	q := *qp
	if q == nil {
		t.size++
		*qp = &Node{Val: val, p: p}
		return true
	}

	c := t.cmp(val, q.Val)
	if c == 0 {
		q.Val = val
		return false
	}

	switch {
	case c < 0:
		c = -1
	case c > 0:
		c = 1
	}
	a := (c + 1) / 2
	fix := t.insert(val, q, &q.c[a])
	if fix {
		return insertFix(int8(c), qp)
	}
	return false

}

// Delete looks up val and dels the matching element
// from the tree. The found element oldval is returned.
//
// Val's Less implementation must be able to handle
// comparisons to elements stored in this tree.
func (t *Tree) Delete(val interface{}) {
	if t == nil {
		return
	}
	t.del(val, &t.root)
}

func (t *Tree) del(val interface{}, qp **Node) bool {
	q := *qp
	if q == nil {
		return false
	}

	c := t.cmp(val, q.Val)
	if c == 0 {
		t.size--
		if q.c[1] == nil {
			if q.c[0] != nil {
				q.c[0].p = q.p
			}
			*qp = q.c[0]
			return true
		}
		fix := delmin(&q.c[1], &q.Val)
		if fix {
			return delFix(-1, qp)
		}
		return false
	}
	switch {
	case c < 0:
		c = -1
	case c > 0:
		c = 1
	}
	a := (c + 1) / 2
	fix := t.del(val, &q.c[a])
	if fix {
		return delFix(int8(-c), qp)
	}
	return false
}

func delmin(qp **Node, min *interface{}) bool {
	q := *qp
	if q.c[0] == nil {
		*min = q.Val
		if q.c[1] != nil {
			q.c[1].p = q.p
		}
		*qp = q.c[1]
		return true
	}
	fix := delmin(&q.c[0], min)
	if fix {
		return delFix(1, qp)
	}
	return false
}

func insertFix(c int8, t **Node) bool {
	s := *t
	if s.b == 0 {
		s.b = c
		return true
	}

	if s.b == -c {
		s.b = 0
		return false
	}

	if s.c[(c+1)/2].b == c {
		s = singlerot(c, s)
	} else {
		s = doublerot(c, s)
	}
	*t = s
	return false
}

func delFix(c int8, t **Node) bool {
	s := *t
	if s.b == 0 {
		s.b = c
		return false
	}

	if s.b == -c {
		s.b = 0
		return true
	}

	a := (c + 1) / 2
	if s.c[a].b == 0 {
		s = rotate(c, s)
		s.b = -c
		*t = s
		return false
	}

	if s.c[a].b == c {
		s = singlerot(c, s)
	} else {
		s = doublerot(c, s)
	}
	*t = s
	return true
}

func singlerot(c int8, s *Node) *Node {
	s.b = 0
	s = rotate(c, s)
	s.b = 0
	return s
}

func doublerot(c int8, s *Node) *Node {
	a := (c + 1) / 2
	r := s.c[a]
	s.c[a] = rotate(-c, s.c[a])
	p := rotate(c, s)
	if r.p != p || s.p != p {
		panic("doublerot: bad parents")
	}

	switch {
	default:
		s.b = 0
		r.b = 0
	case p.b == c:
		s.b = -c
		r.b = 0
	case p.b == -c:
		s.b = 0
		r.b = c
	}

	p.b = 0
	return p
}

func rotate(c int8, s *Node) *Node {
	a := (c + 1) / 2
	r := s.c[a]
	s.c[a] = r.c[a^1]
	if s.c[a] != nil {
		s.c[a].p = s
	}
	r.c[a^1] = s
	r.p = s.p
	s.p = r
	return r
}

// Min returns the minimum element of the AVL tree
// or nil if the tree is empty.
func (t *Tree) Min() *Node {
	return t.bottom(0)
}

// Max returns the maximum element of the AVL tree
// or nil if the tree is empty.
func (t *Tree) Max() *Node {
	return t.bottom(1)
}

func (t *Tree) bottom(d int) *Node {
	n := t.root
	if n == nil {
		return nil
	}

	for c := n.c[d]; c != nil; c = n.c[d] {
		n = c
	}
	return n
}

// Prev returns the previous element in an inorder
// walk of the AVL tree.
func (n *Node) Prev() *Node {
	return n.walk1(0)
}

// Next returns the next element in an inorder
// walk of the AVL tree.
func (n *Node) Next() *Node {
	return n.walk1(1)
}

func (n *Node) walk1(a int) *Node {
	if n == nil {
		return nil
	}

	if n.c[a] != nil {
		n = n.c[a]
		for n.c[a^1] != nil {
			n = n.c[a^1]
		}
		return n
	}

	p := n.p
	for p != nil && p.c[a] == n {
		n = p
		p = p.p
	}
	return p
}

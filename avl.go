// Package avl implements an AVL balanced binary tree.
package avl

import (
	"io/ioutil"
	"log"
)

var dbgLog = log.New(ioutil.Discard, "avl: ", log.LstdFlags)

// Tree holds elements of the AVL tree.
type Tree struct {
	root *Node
	size int
}

// Ordered defines the comparison used to store
// elements in the AVL tree.
type Ordered interface {
	Less(Ordered) bool
}

// A Node holds an Ordered element of the AVL tree in
// the Val field.
type Node struct {
	Val Ordered
	c   [2]*Node
	p   *Node
	b   int8
}

// Size returns the number of elements stored in the tree.
func (t *Tree) Size() int {
	return t.size
}

// Insert inserts the element val into the tree. If
// the value being inserted is already in the tree
// then the original value and true are returned.
// Otherwise the boolean return value is false.
//
// Val's Less implementation must be able to handle
// comparisons to elements stored in this tree.
func (t *Tree) Insert(val Ordered) (Ordered, bool) {
	var insert func(p, q *Node) (*Node, bool)

	var old Ordered
	dup := false
	insert = func(p, q *Node) (*Node, bool) {
		if q == nil {
			new := &Node{Val: val, p: p}
			t.size++
			return new, true
		}

		c := cmp(val, q.Val)
		if c == 0 {
			old = q.Val
			dup = true
			q.Val = val
			return q, false
		}

		a := (c + 1) / 2
		var fix bool
		q.c[a], fix = insert(q, q.c[a])
		if fix {
			return insertfix(c, q)
		}
		return q, false
	}

	t.root, _ = insert(nil, t.root)
	return old, dup
}

// Delete deletes the element Val from the tree and returns
// whether the item was found and deletion was successful.
// Val's Less implementation must be able to handle
// comparisons to elements stored in this tree.
func (t *Tree) Delete(val Ordered) bool {
	del := false
	t.root, _ = delete(t.root, val, &del)
	return del
}

func delete(q *Node, val Ordered, del *bool) (*Node, bool) {
	if q == nil {
		return nil, false
	}

	c := cmp(val, q.Val)
	if c == 0 {
		*del = true
		if q.c[1] == nil {
			if q.c[0] != nil {
				q.c[0].p = q.p
			}
			return q.c[0], true
		}
		var min Ordered
		var fix bool
		q.c[1], fix = deletemin(q.c[1], &min)
		q.Val = min
		if fix {
			return deletefix(-1, q)
		}
		return q, false
	}
	a := (c + 1) / 2
	var fix bool
	q.c[a], fix = delete(q.c[a], val, del)
	if fix {
		return deletefix(-c, q)
	}
	return q, false
}

func deletemin(q *Node, min *Ordered) (*Node, bool) {
	if q.c[0] == nil {
		*min = q.Val
		if q.c[1] != nil {
			q.c[1].p = q.p
		}
		return q.c[1], true
	}
	var fix bool
	q.c[0], fix = deletemin(q.c[0], min)
	if fix {
		return deletefix(1, q)
	}
	return q, false
}

func cmp(a, b Ordered) int8 {
	if a.Less(b) {
		return -1
	}
	if b.Less(a) {
		return 1
	}
	return 0
}

func insertfix(c int8, s *Node) (*Node, bool) {
	if s.b == 0 {
		s.b = c
		return s, true
	}

	if s.b == -c {
		s.b = 0
	} else if s.c[(c+1)/2].b == c {
		s = singlerot(c, s)
	} else {
		s = doublerot(c, s)
	}
	return s, false
}

func deletefix(c int8, s *Node) (*Node, bool) {
	if s.b == 0 {
		s.b = c
		return s, false
	}
	if s.b == -c {
		s.b = 0
		return s, true
	}
	a := (c + 1) / 2
	if s.c[a].b == 0 {
		s = rotate(c, s)
		s.b = -c
		return s, false
	}
	if s.c[a].b == c {
		s = singlerot(c, s)
	} else {
		s = doublerot(c, s)
	}
	return s, true
}

func singlerot(c int8, s *Node) *Node {
	dbgLog.Printf("singlerot: enter %p:%v %d\n", s, s, c)
	s.b = 0
	s = rotate(c, s)
	s.b = 0
	dbgLog.Printf("singlerot: exit %p:%v\n", s, s)
	return s
}

func doublerot(c int8, s *Node) *Node {
	dbgLog.Printf("doublerot: enter %p:%v %d\n", s, s, c)
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
	dbgLog.Printf("doublerot: exit %p:%v\n", s, s)
	return p
}

func rotate(c int8, s *Node) *Node {
	dbgLog.Printf("rotate: enter %p:%v %d\n", s, s, c)
	a := (c + 1) / 2
	r := s.c[a]
	s.c[a] = r.c[a^1]
	if s.c[a] != nil {
		s.c[a].p = s
	}
	r.c[a^1] = s
	r.p = s.p
	s.p = r
	dbgLog.Printf("rotate: exit %p:%v\n", r, r)
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

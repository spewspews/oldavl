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
}

// Ordered defines the comparison used to store
// elements in the AVL tree.
type Ordered interface {
	Less(interface{}) bool
}

// A Node holds an Ordered element of the AVL tree in
// the Val field.
type Node struct {
	Val Ordered
	c   [2]*Node
	p   *Node
	b   int8
}

// Insert inserts the element Val into the tree. Val's Less
// implementation must be able to handle comparisons to
// elements stored in this tree.
func (t *Tree) Insert(Val Ordered) {
	new := &Node{Val: Val}
	t.root, _ = insert(nil, t.root, new)
}

func insert(p, q, new *Node) (*Node, bool) {
	if q == nil {
		new.p = p
		dbgLog.Printf("insert: Inserting %p:%v\n", new, new)
		return new, true
	}

	c := cmp(new.Val, q.Val)
	if c == 0 {
		dbgLog.Printf("insert: collision: %p:%v %p:%v\n", q, q, new, new)
		q.Val = new.Val
		return q, false
	}

	a := (c + 1) / 2
	ch, fix := insert(q, q.c[a], new)
	q.c[a] = ch
	if fix {
		return insertfix(c, q)
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

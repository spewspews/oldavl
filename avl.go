// Package avl implements an AVL balanced binary tree.
package avl

// Tree holds elements of the AVL tree.
type Tree struct {
	root *node
}

// Ordered defines the comparison used to store
// elements in the AVL tree.
type Ordered interface {
	Less(interface{}) bool
}

type node struct {
	val     Ordered
	child   [2]*node
	parent  *node
	balance int8
}

// Insert inserts the element val into the tree. val's Less
// implementation must be able to handle comparisons to
// elements stored in this tree.
func (tree *Tree) Insert(val Ordered) {
	new := &node{val: val}
	tree.root, _ = insert(nil, tree.root, new)
}

func insert(p, q, new *node) (*node, bool) {
	if q == nil {
		new.parent = p
		return new, true
	}

	c := cmp(new.val, q.val)
	if c == 0 {
		return new, true
	}

	a := (c + 1) / 2
	child, fix := insert(q, q.child[a], new)
	q.child[a] = child

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

func insertfix(c int8, s *node) (*node, bool) {
	if s.balance == 0 {
		s.balance = c
		return s, true
	}
	if s.balance == -c {
		s.balance = 0
		return s, false
	}
	if s.child[(c+1)/2].balance == c {
		s = singlerot(c, s)
	} else {
		s = doublerot(c, s)
	}
	return s, false
}

func singlerot(c int8, s *node) *node {
	s.balance = 0
	s = rotate(c, s)
	s.balance = 0
	return s
}

func doublerot(c int8, s *node) *node {
	a := (c + 1) / 2
	r := s.child[a]
	s.child[a] = rotate(-c, s.child[a])
	p := rotate(c, s)
	if r.parent != p || s.parent != p {
		panic("doublerot: bad parents")
	}

	switch {
	default:
		s.balance = 0
		r.balance = 0
	case p.balance == c:
		s.balance = -c
		r.balance = 0
	case p.balance == -c:
		s.balance = 0
		r.balance = c
	}

	p.balance = 0
	return p

}

func rotate(c int8, s *node) *node {
	a := (c + 1) / 2
	r := s.child[a]
	s.child[a] = r.child[a^1]
	if s.child[a] != nil {
		s.child[a].parent = s
	}
	r.child[a^1] = s
	r.parent = s.parent
	s.parent = r
	return r
}

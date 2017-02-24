package avl

import (
	"math/rand"
	"testing"
	"time"
)

const (
	randMax = 100000
	nodes   = 1000
	dels    = 300
)

type Int int

func (i Int) Less(j Ordered) bool {
	return i < j.(Int)
}

func TestInsertOrdered(t *testing.T) {
	tree := newIntTree(nodes)
	tree.checkOrdered(t)
}

func TestInsertBalanced(t *testing.T) {
	tree := newIntTree(nodes)
	tree.checkBalance(t)
}

func TestDeleteOrdered(t *testing.T) {
	tree := newIntTree(nodes)
	tree.deleteSome(dels)
	tree.checkOrdered(t)
}

func TestBalanceOrdered(t *testing.T) {
	tree := newIntTree(nodes)
	tree.deleteSome(dels)
	tree.checkBalance(t)
}

func (tree *Tree) checkOrdered(t *testing.T) {
	n := tree.Min()
	for next := n.Next(); next != nil; next = n.Next() {
		if next.Val.(Int) <= n.Val.(Int) {
			t.Errorf("Tree not ordered: %d ≮ %d", n.Val.(Int), next.Val.(Int))
		}
		n = next
	}
}

func (tree *Tree) checkBalance(t *testing.T) {
	for n := tree.Min(); n != nil; n = n.Next() {
		if !checkBalance(n) {
			t.Errorf("Tree not balanced")
		}
	}
}

func newIntTree(n int) *Tree {
	tree := new(Tree)
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < n; i++ {
		r := Int(rng.Intn(randMax))
		tree.Insert(r)
	}
	return tree
}

func (tree *Tree) deleteSome(n int) {
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < n; i++ {
		r := Int(rng.Intn(randMax))
		tree.Delete(r)
	}
}

func checkBalance(n *Node) bool {
	left := depth(n.c[0])
	right := depth(n.c[1])
	b := right - left
	if int8(b) != n.b {
		return false
	}
	return true
}

func depth(n *Node) int {
	if n == nil {
		return 0
	}

	ld := depth(n.c[0])
	rd := depth(n.c[1])
	if ld >= rd {
		return ld + 1
	}
	return rd + 1
}

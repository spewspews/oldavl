package avl

import (
	"math/rand"
	"testing"
	"time"
)

const (
	randMax = 1000000
	nodes   = 1000
)

type Int int

func (i Int) Less(j interface{}) bool {
	return i < j.(Int)
}

func TestInsertOrdered(t *testing.T) {
	tree := newIntTree(nodes)
	n := tree.Min()
	for next := n.Next(); next != nil; next = n.Next() {
		if next.Val.(Int) <= n.Val.(Int) {
			t.Errorf("Tree not ordered: %d ≮ %d", n.Val.(Int), next.Val.(Int))
		}
		n = next
	}
}

func TestInsertBalanced(t *testing.T) {
	tree := newIntTree(nodes)
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

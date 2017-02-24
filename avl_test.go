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

type IntString struct {
	key int
	val string
}

func (is *IntString) Less(j Ordered) bool {
	return is.key < j.(*IntString).key
}

func TestInsertOrdered(t *testing.T) {
	tree := newIntTree(nodes)
	tree.checkOrdered(t)
}

func TestInsertBalanced(t *testing.T) {
	tree := newIntTree(nodes)
	tree.checkBalance(t)
}

func TestInsertDuplicates(t *testing.T) {
	var old *Ordered
	tree := new(Tree)

	old = tree.Insert(Int(5))
	if old != nil {
		t.Errorf("got bad duplicate %d\n", (*old).(Int))
	}

	old = tree.Insert(Int(6))
	if old != nil {
		t.Errorf("got bad duplicate %d\n", (*old).(Int))
	}

	old = tree.Insert(Int(5))
	if old == nil {
		t.Error("Should have gotten duplicate")
	}
	t.Logf("Duplicate value is %d\n", (*old).(Int))
}

func TestInsertKeyValDuplicates(t *testing.T) {
	var old *Ordered
	tree := new(Tree)

	old = tree.Insert(&IntString{3, "three"})
	if old != nil {
		t.Errorf("got bad duplicate %d\n", (*old).(*IntString).key)
	}

	old = tree.Insert(&IntString{4, "four"})
	if old != nil {
		t.Errorf("got bad duplicate %d\n", (*old).(*IntString).key)
	}

	old = tree.Insert(&IntString{3, "newthree"})
	if old == nil {
		t.Error("Should have gotten duplicate")
	}
	s := (*old).(*IntString).val
	if s != "three" {
		t.Error("Got the wrong value")
	}
	t.Logf("Duplicate value is %s\n", s)
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

package avl

import (
	"math/rand"
	"testing"
	"time"
)

const (
	randMax = 10000
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
	tree, _ := newIntTree(nodes)
	tree.checkOrdered(t)
}

func TestInsertBalanced(t *testing.T) {
	tree, _ := newIntTree(nodes)
	tree.checkBalance(t)
}

func TestInsertSize(t *testing.T) {
	tree, size := newIntTree(nodes)
	if size != tree.Size() {
		t.Errorf("Size does not match: size %d, tree.Size() %d\n", size, tree.Size())
	}
}

func TestInsertDuplicates(t *testing.T) {
	tree := new(Tree)

	old, found := tree.Insert(Int(5))
	if found {
		t.Errorf("got bad duplicate %d\n", old.(Int))
	}

	old, found = tree.Insert(Int(6))
	if found {
		t.Errorf("got bad duplicate %d\n", old.(Int))
	}

	old, found = tree.Insert(Int(5))
	if !found {
		t.Error("Should have gotten duplicate")
	}
	v := old.(Int)
	if v != 5 {
		t.Error("Got the wrong value")
	}
}

func TestInsertKeyValDuplicates(t *testing.T) {
	tree := new(Tree)

	old, found := tree.Insert(&IntString{3, "three"})
	if found {
		v := old.(*IntString)
		t.Errorf("got bad duplicate %d\n", v.key)
	}

	old, found = tree.Insert(&IntString{4, "four"})
	if found {
		v := old.(*IntString)
		t.Errorf("got bad duplicate %d\n", v.key)
	}

	old, found = tree.Insert(&IntString{3, "three"})
	if !found {
		t.Error("Should have gotten duplicate")
	}
	v := old.(*IntString)
	if v.key != 3 || v.val != "three" {
		t.Errorf("Got the wrong values: %d %s", v.key, v.val)
	}
}

func TestDeleteOrdered(t *testing.T) {
	tree, _ := newIntTree(nodes)
	tree.deleteSome(dels)
	tree.checkOrdered(t)
}

func TestBalanceOrdered(t *testing.T) {
	tree, _ := newIntTree(nodes)
	tree.deleteSome(dels)
	tree.checkBalance(t)
}

func TestDeleteSize(t *testing.T) {
	tree, _ := newIntTree(nodes)
	oldsize := tree.Size()
	dels := tree.deleteSome(dels)
	if tree.Size() != oldsize-dels {
		t.Errorf("Size does not match: oldsize-dels %d, tree.Size() %d", oldsize-dels, tree.Size())
	}
}

func TestLookups(t *testing.T) {
	tree := new(Tree)
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	vals := make(map[Int]bool)
	for i := 0; i < nodes; i++ {
		r := Int(rng.Intn(randMax))
		tree.Insert(r)
		vals[r] = true
	}
	for i := 0; i < randMax; i++ {
		inMap := vals[Int(i)]
		_, inTree := tree.Lookup(Int(i))
		msg := ""
		switch inMap {
		case true:
			if !inTree {
				msg = "value found in map but not in tree"
			}
		case false:
			if inTree {
				msg = "value found in tree but not in map"
			}
		}
		if msg != "" {
			t.Errorf("Mismatch between map and tree: %s", msg)
		}
	}
}

func TestLookupsAfterDeletions(t *testing.T) {
	tree := new(Tree)
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	vals := make(map[Int]bool)
	for i := 0; i < nodes; i++ {
		r := Int(rng.Intn(randMax))
		tree.Insert(r)
		vals[r] = true
	}
	if len(vals) != tree.Size() {
		t.Errorf("Size mismatch between map and tree: %d %d\n", len(vals), tree.Size())
	}
	t.Logf("Inserted %d elements\n", tree.Size())
	oldSize := tree.Size()
	deleted := 0
	for i := 0; i < dels; i++ {
		r := Int(rng.Intn(randMax))
		if _, found := tree.Delete(r); found {
			deleted++
		}
		delete(vals, r)
	}
	newSize := oldSize - deleted
	if len(vals) != newSize {
		t.Errorf("There should be %d values in the map\n", newSize)
	}
	if tree.Size() != newSize {
		t.Errorf("there should be %d values in the tree\n", newSize)
	}
	t.Logf("Succesfully deleted %d values\n", newSize)
	for i := 0; i < randMax; i++ {
		inMap := vals[Int(i)]
		_, inTree := tree.Lookup(Int(i))
		msg := ""
		switch inMap {
		case true:
			if !inTree {
				msg = "value found in map but not in tree"
			}
		case false:
			if inTree {
				msg = "value found in tree but not in map"
			}
		}
		if msg != "" {
			t.Errorf("Mismatch between map and tree: %s", msg)
		}
	}
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

func newIntTree(n int) (*Tree, int) {
	tree := new(Tree)
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	ins := 0
	for i := 0; i < n; i++ {
		r := Int(rng.Intn(randMax))
		if _, found := tree.Insert(r); !found {
			ins++
		}
	}
	return tree, ins
}

func (tree *Tree) deleteSome(n int) (dels int) {
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < n; i++ {
		r := Int(rng.Intn(randMax))
		if _, found := tree.Delete(r); found {
			dels++
		}
	}
	return
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

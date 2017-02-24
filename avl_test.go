package avl

import (
	"math/rand"
	"testing"
	"time"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

const (
	randMax = 2000
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
	tree, _ := newIntTree(nodes, randMax)
	tree.checkOrdered(t)
}

func TestInsertBalanced(t *testing.T) {
	tree, _ := newIntTree(nodes, randMax)
	tree.checkBalance(t)
}

func TestInsertSize(t *testing.T) {
	tree, size := newIntTree(nodes, randMax)
	if size != tree.Size() {
		t.Errorf("Size does not match: size %d, tree.Size() %d\n", size, tree.Size())
	}
}

func TestWalk(t *testing.T) {
	tree, _ := newIntTree(nodes, randMax)
	i := 0
	for n := tree.Min(); n != nil; n = n.Next() {
		i++
	}
	if i != tree.Size() {
		t.Errorf("Walk up not the same size as tree: %d %d\n", i, tree.Size())
	}
	i = 0
	for n := tree.Max(); n != nil; n = n.Prev() {
		i++
	}
	if i != tree.Size() {
		t.Errorf("Walk down not the same size as tree: %d %d\n", i, tree.Size())
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
	tree, _ := newIntTree(nodes, randMax)
	tree.deleteSome(dels)
	tree.checkOrdered(t)
}

func TestDeleteBalanced(t *testing.T) {
	tree, _ := newIntTree(nodes, randMax)
	t.Logf("Tree has %d elements\n", tree.Size())
	d := tree.deleteSome(dels)
	t.Logf("Deleted %d elements\n", d)
	tree.checkBalance(t)
}

func TestDeleteSize(t *testing.T) {
	tree, _ := newIntTree(nodes, randMax)
	oldsize := tree.Size()
	dels := tree.deleteSome(dels)
	if tree.Size() != oldsize-dels {
		t.Errorf("Size does not match: oldsize-dels %d, tree.Size() %d", oldsize-dels, tree.Size())
	}
}

func TestDeleteWalk(t *testing.T) {
	tree, _ := newIntTree(nodes, randMax)
	t.Logf("Tree has %d elements\n", tree.Size())
	d := tree.deleteSome(dels)
	t.Logf("Deleted %d elements\n", d)
	i := 0
	for n := tree.Min(); n != nil; n = n.Next() {
		i++
	}
	if i != tree.Size() {
		t.Errorf("Walk up not the same size as tree: %d %d\n", i, tree.Size())
	}
	i = 0
	for n := tree.Max(); n != nil; n = n.Prev() {
		i++
	}
	if i != tree.Size() {
		t.Errorf("Walk down not the same size as tree: %d %d\n", i, tree.Size())
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
		if !n.checkBalance(t) {
			t.Errorf("Tree not balanced")
		}
	}
}

func (n *Node) checkBalance(t *testing.T) bool {
	left := depth(n.c[0])
	right := depth(n.c[1])
//	t.Logf("Balance is %d %d\n", left, right)
	b := right - left
	if int8(b) != n.b {
		return false
	}
	return true
}

func newIntTree(n, randMax int) (*Tree, int) {
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

func BenchmarkLookup100(b *testing.B) {
	benchmarkLookup(b, 100)
}

func BenchmarkLookup1000(b *testing.B) {
	benchmarkLookup(b, 1000)
}

func BenchmarkLookup10000(b *testing.B) {
	benchmarkLookup(b, 10000)
}

func BenchmarkLookup100000(b *testing.B) {
	benchmarkLookup(b, 100000)
}

func benchmarkLookup(b *testing.B, size int) {
	b.StopTimer()
	tree := new(Tree)
	for n := 0; n < size; n++ {
		tree.Insert(Int(n))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Lookup(Int(n))
		}
	}
}

func BenchmarkLookupRandom100(b *testing.B) {
	benchmarkLookupRandom(b, 100)
}

func BenchmarkLookupRandom1000(b *testing.B) {
	benchmarkLookupRandom(b, 1000)
}

func BenchmarkLookupRandom10000(b *testing.B) {
	benchmarkLookupRandom(b, 10000)
}

func BenchmarkLookupRandom100000(b *testing.B) {
	benchmarkLookupRandom(b, 100000)
}

func benchmarkLookupRandom(b *testing.B, size int) {
	b.StopTimer()
	tree, _ := newIntTree(size, size)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Lookup(Int(n))
		}
	}
}

func BenchmarkRedBlackGet100(b *testing.B) {
	benchmarkRedBlackGet(b, 100)
}

func BenchmarkRedBlackGet1000(b *testing.B) {
	benchmarkRedBlackGet(b, 1000)
}

func BenchmarkRedBlackGet10000(b *testing.B) {
	benchmarkRedBlackGet(b, 10000)
}

func BenchmarkRedBlackGet100000(b *testing.B) {
	benchmarkRedBlackGet(b, 100000)
}

func benchmarkRedBlackGet(b *testing.B, size int) {
	b.StopTimer()
	tree := rbt.NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Get(n)
		}
	}
}

func BenchmarkRedBlackGetRandom100(b *testing.B) {
	benchmarkRedBlackGetRandom(b, 100)
}

func BenchmarkRedBlackGetRandom1000(b *testing.B) {
	benchmarkRedBlackGetRandom(b, 1000)
}

func BenchmarkRedBlackGetRandom10000(b *testing.B) {
	benchmarkRedBlackGetRandom(b, 10000)
}

func BenchmarkRedBlackGetRandom100000(b *testing.B) {
	benchmarkRedBlackGetRandom(b, 100000)
}

func benchmarkRedBlackGetRandom(b *testing.B, size int) {
	b.StopTimer()
	tree := rbt.NewWithIntComparator()
	seed := time.Now().UTC().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	for n := 0; n < size; n++ {
		tree.Put(rng.Intn(size), struct{}{})
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Get(n)
		}
	}
}

package avl

import (
	"math/rand"
	"testing"
	"time"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

const (
	randMax = 2000
	nNodes   = 1000
	nDels    = 300
)

var rng *rand.Rand

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

func TestMain(m *testing.M) {
	seed := time.Now().UTC().UnixNano()
	rng = rand.New(rand.NewSource(seed))
	m.Run()
}

func TestInsertOrdered(t *testing.T) {
	tree := newIntTree(nNodes, randMax)
	tree.checkOrdered(t)
}

func TestInsertBalanced(t *testing.T) {
	tree := newIntTree(nNodes, randMax)
	tree.checkBalance(t)
}

func TestInsertSize(t *testing.T) {
	tree, vals := newIntTreeAndMap(nNodes, randMax)
	if len(vals) != tree.Size() {
		t.Errorf("Size does not match: size %d, tree.Size() %d\n", len(vals), tree.Size())
	}
}

func TestInsertReturn(t *testing.T) {
	tree := new(Tree)
	for i := 0; i < 10; i += 2 {
		_, found := tree.Insert(Int(i))
		if found {
			t.Errorf("Should not have found duplicate on first loop: %d\n", i)
		}
	}
	for i := 0; i < 10; i += 2 {
		j, found := tree.Insert(Int(i))
		if !found {
			t.Errorf("Did not find duplicate on second loop: %d\n", i)
		}
		if j.(Int) != Int(i) {
			t.Errorf("Got the wrong value %d %d\n", i, j)
		}
	}
	for i := 1; i < 10; i += 2 {
		_, found := tree.Insert(Int(i))
		if found {
			t.Errorf("Should not have found duplicate on second loop: %d\n", i)
		}
	}
}

func TestWalk(t *testing.T) {
	tree := newIntTree(nNodes, randMax)
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
	tree := newIntTree(nNodes, randMax)
	tree.deleteSome(nDels)
	tree.checkOrdered(t)
}

func TestDeleteBalanced(t *testing.T) {
	tree := newIntTree(nNodes, randMax)
	t.Logf("Tree has %d elements\n", tree.Size())
	d := tree.deleteSome(nDels)
	t.Logf("Deleted %d elements\n", d)
	tree.checkBalance(t)
}

func TestDeleteSize(t *testing.T) {
	tree := newIntTree(nNodes, randMax)
	oldsize := tree.Size()
	nDels := tree.deleteSome(nDels)
	if tree.Size() != oldsize-nDels {
		t.Errorf("Size does not match: oldsize-nDels %d, tree.Size() %d", oldsize-nDels, tree.Size())
	}
}

func TestDeleteWalk(t *testing.T) {
	tree := newIntTree(nNodes, randMax)
	t.Logf("Tree has %d elements\n", tree.Size())
	d := tree.deleteSome(nDels)
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
	tree, vals := newIntTreeAndMap(nNodes, randMax)
	tree.checkLookups(t, vals, randMax)
}

func TestLookupsAfterDeletions(t *testing.T) {
	tree, vals := newIntTreeAndMap(nNodes, randMax)
	tree.deleteSomeAndMap(nDels, vals)
	tree.checkLookups(t, vals, randMax)
}

func (tree *Tree) checkLookups(t *testing.T, vals map[Int]bool, max int) {
	for i := 0; i < max; i++ {
		inMap := vals[Int(i)]
		_, inTree := tree.Lookup(Int(i))
		if inMap && !inTree {
			t.Errorf("Mismatch: value found in map but not in tree")
		}
		if inTree && !inMap {
			t.Errorf("Mismatch: value found in tree but not in map")
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

func newIntTree(n, randMax int) *Tree {
	tree := new(Tree)
	for i := 0; i < n; i++ {
		tree.Insert(Int(rng.Intn(randMax)))
	}
	return tree
}

func newIntTreeAndMap(n, randMax int) (tree *Tree, vals map[Int]bool) {
	tree = new(Tree)
	vals = make(map[Int]bool)
	for i := 0; i < nNodes; i++ {
		r := Int(rng.Intn(randMax))
		tree.Insert(r)
		vals[r] = true
	}
	return
}

func (tree *Tree) deleteSome(n int) (nDels int) {
	for i := 0; i < n; i++ {
		r := Int(rng.Intn(randMax))
		if _, found := tree.Delete(r); found {
			nDels++
		}
	}
	return
}

func (tree *Tree) deleteSomeAndMap(n int, vals map[Int]bool) (nDels int) {
	for i := 0; i < n; i++ {
		r := Int(rng.Intn(randMax))
		if _, found := tree.Delete(r); found {
			nDels++
		}
		delete(vals, r)
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
	tree := newIntTree(size, size)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Lookup(Int(n))
		}
	}
}

func BenchmarkInsert100(b *testing.B) {
	benchmarkInsert(b, 100)
}

func BenchmarkInsert1000(b *testing.B) {
	benchmarkInsert(b, 1000)
}

func BenchmarkInsert10000(b *testing.B) {
	benchmarkInsert(b, 10000)
}

func BenchmarkInsert100000(b *testing.B) {
	benchmarkInsert(b, 100000)
}

func benchmarkInsert(b *testing.B, size int) {
	tree := new(Tree)
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Insert(Int(n))
		}
		tree = new(Tree)
	}
}

func BenchmarkInsertRandom100(b *testing.B) {
	benchmarkInsertRandom(b, 100)
}

func BenchmarkInsertRandom1000(b *testing.B) {
	benchmarkInsertRandom(b, 1000)
}

func BenchmarkInsertRandom10000(b *testing.B) {
	benchmarkInsertRandom(b, 10000)
}

func BenchmarkInsertRandom100000(b *testing.B) {
	benchmarkInsertRandom(b, 100000)
}

func benchmarkInsertRandom(b *testing.B, size int) {
	b.StopTimer()
	vals := make([]Int, size)
	for i := range vals {
		vals[i] = Int(rng.Intn(size))
	}
	tree := new(Tree)
	b.StartTimer()
	for t := 0; t < b.N; t++ {
		for _, n := range vals {
			tree.Insert(n)
		}
		tree = new(Tree)
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

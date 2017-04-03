package avl

import (
	"math/rand"
	"testing"
	"time"

	"github.com/emirpasic/gods/trees/avltree"
)

const (
	randMax = 2000
	nNodes  = 1000
	nDels   = 300
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
	tree := newRandIntTree(nNodes, randMax)
	tree.checkOrdered(t)
}

func TestInsertBalanced(t *testing.T) {
	tree := newRandIntTree(nNodes, randMax)
	tree.checkBalance(t)
}

func TestInsertSize(t *testing.T) {
	tree, vals := newRandIntTreeAndMap(nNodes, randMax)
	if len(vals) != tree.Size() {
		t.Errorf("Size does not match: size %d, tree.Size() %d\n", len(vals), tree.Size())
	}
}

func TestWalk(t *testing.T) {
	tree := newRandIntTree(nNodes, randMax)
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

func TestDeleteOrdered(t *testing.T) {
	tree := newRandIntTree(nNodes, randMax)
	tree.deleteSome(nDels)
	tree.checkOrdered(t)
}

func TestDeleteBalanced(t *testing.T) {
	tree := newRandIntTree(nNodes, randMax)
	t.Logf("Tree has %d elements\n", tree.Size())
	tree.deleteSome(nDels)
	tree.checkBalance(t)
}

func TestDeleteWalk(t *testing.T) {
	tree := newRandIntTree(nNodes, randMax)
	t.Logf("Tree has %d elements\n", tree.Size())
	tree.deleteSome(nDels)
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
	tree, vals := newRandIntTreeAndMap(nNodes, randMax)
	tree.checkLookups(t, vals, randMax)
}

func TestLookupsAfterDeletions(t *testing.T) {
	tree, vals := newRandIntTreeAndMap(nNodes, randMax)
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

func newRandIntTree(n, randMax int) *Tree {
	tree := new(Tree)
	for i := 0; i < n; i++ {
		tree.Insert(Int(rng.Intn(randMax)))
	}
	return tree
}

func newRandIntTreeAndMap(n, randMax int) (tree *Tree, vals map[Int]bool) {
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
		tree.Delete(Int(rng.Intn(randMax)))
	}
	return
}

func (tree *Tree) deleteSomeAndMap(n int, vals map[Int]bool) (nDels int) {
	for i := 0; i < n; i++ {
		r := Int(rng.Intn(randMax))
		tree.Delete(r)
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

func BenchmarkGoDSGet100(b *testing.B) {
	benchmarkGoDSGet(b, 100)
}

func BenchmarkGoDSGet1000(b *testing.B) {
	benchmarkGoDSGet(b, 1000)
}

func BenchmarkGoDSGet10000(b *testing.B) {
	benchmarkGoDSGet(b, 10000)
}

func BenchmarkGoDSGet100000(b *testing.B) {
	benchmarkGoDSGet(b, 100000)
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

func BenchmarkGoDSGetRandom100(b *testing.B) {
	benchmarkGoDSGetRandom(b, 100)
}

func BenchmarkGoDSGetRandom1000(b *testing.B) {
	benchmarkGoDSGetRandom(b, 1000)
}

func BenchmarkGoDSGetRandom10000(b *testing.B) {
	benchmarkGoDSGetRandom(b, 10000)
}

func BenchmarkGoDSGetRandom100000(b *testing.B) {
	benchmarkGoDSGetRandom(b, 100000)
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

func BenchmarkGoDSPut100(b *testing.B) {
	benchmarkGoDSPut(b, 100)
}

func BenchmarkGoDSPut1000(b *testing.B) {
	benchmarkGoDSPut(b, 1000)
}

func BenchmarkGoDSPut10000(b *testing.B) {
	benchmarkGoDSPut(b, 10000)
}

func BenchmarkGoDSPut100000(b *testing.B) {
	benchmarkGoDSPut(b, 100000)
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

func BenchmarkGoDSPutRandom100(b *testing.B) {
	benchmarkGoDSPutRandom(b, 100)
}

func BenchmarkGoDSPutRandom1000(b *testing.B) {
	benchmarkGoDSPutRandom(b, 1000)
}

func BenchmarkGoDSPutRandom10000(b *testing.B) {
	benchmarkGoDSPutRandom(b, 10000)
}

func BenchmarkGoDSPutRandom100000(b *testing.B) {
	benchmarkGoDSPutRandom(b, 100000)
}

func benchmarkInsertRandom(b *testing.B, size int) {
	b.StopTimer()
	vals := make([]Int, size);
	for i := range vals {
		vals[i] = Int(rng.Intn(size))
	}
	b.StartTimer()
	for t := 0; t < b.N; t++ {
		tree := new(Tree)
		for _, n := range vals {
			tree.Insert(n)
		}
	}
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

func benchmarkInsert(b *testing.B, size int) {
	for i := 0; i < b.N; i++ {
		tree := new(Tree)
		for n := 0; n < size; n++ {
			tree.Insert(Int(n))
		}
	}
}

func benchmarkGoDSPut(b *testing.B, size int) {
	for i := 0; i < b.N; i++ {
		tree := avltree.NewWithIntComparator()
		for n := 0; n < size; n++ {
			tree.Put(n, nil)
		}
	}
}

func benchmarkGoDSPutRandom(b *testing.B, size int) {
	b.StopTimer()
	vals := make([]int, size);
	for i := range vals {
		vals[i] = rng.Intn(size)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree := avltree.NewWithIntComparator()
		for _, n := range vals {
			tree.Put(n, nil)
		}
	}
}

func benchmarkGoDSGet(b *testing.B, size int) {
	b.StopTimer()
	tree := avltree.NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.Put(n, nil)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Get(n)
		}
	}
}

func benchmarkLookupRandom(b *testing.B, size int) {
	b.StopTimer()
	tree := newRandIntTree(size, size)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Lookup(Int(n))
		}
	}
}

func benchmarkGoDSGetRandom(b *testing.B, size int) {
	b.StopTimer()
	tree := avltree.NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.Put(rng.Intn(size), nil)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Get(n)
		}
	}
}

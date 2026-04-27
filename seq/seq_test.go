package seq

import (
	"iter"
	"slices"
	"strconv"
	"strings"
	"testing"

	testcmp "github.com/google/go-cmp/cmp"
	"github.com/kagadar/go-pipeline/predicates"
)

func TestAny(t *testing.T) {
	if !Any(slices.Values([]int{1, 2, 3}), predicates.IsEven) {
		t.Error("Any({1,2,3}, IsEven) got false, want true")
	}
	if Any(slices.Values([]int{1, 3, 5}), predicates.IsEven) {
		t.Error("Any({1,3,5}, IsEven) got true, want false")
	}
}

func TestChunk(t *testing.T) {
	if diff := testcmp.Diff(
		slices.Collect(Chunk(slices.Values([]int{1, 2, 3}), 2)),
		[][]int{{1, 2}, {3}},
	); diff != "" {
		t.Errorf("Chunk({1,2,3}, 2) unexpected diff (-got +want):\n%s", diff)
	}
}

func TestChunk_Panic(t *testing.T) {
	defer func() {
		if got, want := recover(), "size cannot be less than 1"; got != want {
			t.Errorf("Chunk(size 0) did not panic with expected error, got %v, want %q", got, want)
		}
	}()
	Chunk(slices.Values([]int{1}), 0)
}

func TestConcat(t *testing.T) {
	if diff := testcmp.Diff(
		slices.Collect(Concat(slices.Values([]int{1, 2}), slices.Values([]int{3}))),
		[]int{1, 2, 3},
	); diff != "" {
		t.Errorf("Concat({1,2}, {3}) unexpected diff (-got +want):\n%s", diff)
	}
}

func TestContains(t *testing.T) {
	if !Contains(slices.Values([]int{1, 2, 3}), 2) {
		t.Error("Contains({1,2,3}, 2) got false, want true")
	}
	if Contains(slices.Values([]int{1, 2, 3}), 4) {
		t.Error("Contains({1,2,3}, 4) got true, want false")
	}
}

func TestCount(t *testing.T) {
	if got := Count(slices.Values([]int{1, 2, 3})); got != 3 {
		t.Errorf("Count({1,2,3}) got %d, want 3", got)
	}
}

func TestFind(t *testing.T) {
	if e, ok := Find(slices.Values([]int{1, 2, 3}), predicates.IsEven); !ok || e != 2 {
		t.Errorf("Find({1,2,3}, IsEven) got %d %t, want 2 true", e, ok)
	}
	if e, ok := Find(slices.Values([]int{1, 3, 5}), predicates.IsEven); ok || e != 0 {
		t.Errorf("Find({1,3,5}, IsEven) got %d %t, want 0 false", e, ok)
	}
}

func TestGroupBy(t *testing.T) {
	if diff := testcmp.Diff(
		GroupBy(slices.Values([]int{1, 2, 3, 4}), predicates.FindParity),
		map[predicates.Parity][]int{
			predicates.Odd:  {1, 3},
			predicates.Even: {2, 4},
		},
	); diff != "" {
		t.Errorf("GroupBy(Parity) unexpected diff (-got +want):\n%s", diff)
	}
}

func TestLast(t *testing.T) {
	if got := Last(slices.Values([]int{1, 2, 3})); got != 3 {
		t.Errorf("Last({1,2,3}) got %d, want 3", got)
	}
}

func TestPartition(t *testing.T) {
	evens, odds := Partition(slices.Values([]int{1, 2, 3, 4}), predicates.IsEven)
	if ediff, odiff := testcmp.Diff(evens, []int{2, 4}), testcmp.Diff(odds, []int{1, 3}); ediff != "" || odiff != "" {
		t.Errorf("Partition({1,2,3,4}, IsEven) unexpected diff (-got +want):\n%s", strings.Join([]string{ediff, odiff}, "\n"))
	}
}

func TestSkip(t *testing.T) {
	if diff := testcmp.Diff(
		slices.Collect(Skip(slices.Values([]int{1, 2, 3}), 1)),
		[]int{2, 3},
	); diff != "" {
		t.Errorf("Skip({1,2,3}, 1) unexpected diff (-got +want):\n%s", diff)
	}
}

func TestTake(t *testing.T) {
	if diff := testcmp.Diff(
		slices.Collect(Take(slices.Values([]int{1, 2, 3}), 2)),
		[]int{1, 2},
	); diff != "" {
		t.Errorf("Take({1,2,3}, 2) unexpected diff (-got +want):\n%s", diff)
	}
}

func coverBreak[E any](i iter.Seq[E]) {
	for range i {
		break
	}
}

func TestYieldCoverage(t *testing.T) {
	v := slices.Values([]int{2, 1, 6})
	var err error
	coverBreak(Chunk(v, 1))
	coverBreak(Concat(v))
	coverBreak(Dedupe(v))
	coverBreak(Filter(v, predicates.IsEven))
	coverBreak(Flatten(slices.Values([][]int{{1}})))
	coverBreak(FlattenSeq(slices.Values([]iter.Seq[int]{slices.Values([]int{1})})))
	coverBreak(Skip(v, 1))
	coverBreak(SkipWhile(v, predicates.IsEven))
	coverBreak(Take(v, 1))
	coverBreak(TakeWhile(v, predicates.IsEven))
	coverBreak2(ToSeq2(v, func(int) (int, int) { return 0, 0 }))
	coverBreak(Transform(v, strconv.Itoa))
	coverBreak(TransformErr(v, func(int) (int, error) { return 0, nil }, &err))
}

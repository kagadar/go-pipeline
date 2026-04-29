package seq

import (
	"errors"
	"iter"
	"maps"
	"slices"
	"testing"

	testcmp "github.com/google/go-cmp/cmp"
	"github.com/kagadar/go-pipeline/predicates"
)

func coverBreak2[K, V any](i iter.Seq2[K, V]) {
	for range i {
		break
	}
}

func TestAppendErr(t *testing.T) {
	got, err := AppendErr([]int{}, slices.All([]error{nil, nil}))
	if err != nil {
		t.Fatalf("AppendErr() unexpected error: %v", err)
	}
	if diff := testcmp.Diff(got, []int{0, 1}); diff != "" {
		t.Errorf("AppendErr() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestAppendErr_Error(t *testing.T) {
	want := errors.New("2")
	got, err := AppendErr([]int{}, slices.All([]error{nil, want}))
	if got != nil || err != want {
		t.Errorf("AppendErr(fail) got %v %v, want [] %v", got, err, want)
	}
}

func TestCollectErr(t *testing.T) {
	got, err := CollectErr(slices.All([]error{nil, nil}))
	if err != nil {
		t.Fatalf("CollectErr() unexpected error: %v", err)
	}
	if diff := testcmp.Diff(got, []int{0, 1}); diff != "" {
		t.Errorf("CollectErr() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestCollectErr_Error(t *testing.T) {
	want := errors.New("2")
	got, err := CollectErr(slices.All([]error{nil, want}))
	if got != nil || err != want {
		t.Errorf("CollectErr(fail) got %v %v, want [] %v", got, err, want)
	}
}

func TestConcat2(t *testing.T) {
	if diff := testcmp.Diff(
		maps.Collect(Concat2(maps.All(map[int]struct{}{1: {}, 2: {}}), maps.All(map[int]struct{}{3: {}}))),
		map[int]struct{}{1: {}, 2: {}, 3: {}},
	); diff != "" {
		t.Errorf("Concat2({1,2}, {3}) unexpected diff (-got +want):\n%s", diff)
	}
}

func TestCount2(t *testing.T) {
	if got := Count2(slices.All([]int{1, 2, 3})); got != 3 {
		t.Errorf("Count2({1,2,3}) got %d, want 3", got)
	}
}

func TestFind2(t *testing.T) {
	if k, v, ok := Find2(maps.All(map[int]string{1: "1", 2: "2", 3: "3"}), predicates.Keys[int, string](predicates.IsEven)); !ok || k != 2 || v != "2" {
		t.Errorf("Find2({1,2,3}, IsEven) got %d %q %t, want 2 \"2\" true", k, v, ok)
	}
	if k, v, ok := Find2(maps.All(map[int]string{1: "1", 3: "3", 5: "5"}), predicates.Keys[int, string](predicates.IsEven)); ok || k != 0 || v != "" {
		t.Errorf("Find2({1,3,5}, IsEven) got %d %q %t, want 0 \"\" false", k, v, ok)
	}
}

func TestInsertMap(t *testing.T) {
	in := map[int]struct{}{1: {}, 2: {}, 3: {}}
	got := InsertMap(in, maps.All(map[int]struct{}{3: {}, 4: {}, 5: {}}))
	want := map[int]struct{}{1: {}, 2: {}, 3: {}, 4: {}, 5: {}}
	if diff := testcmp.Diff(in, want); diff != "" {
		t.Errorf("input map after InsertMap() unexpected diff (-got +want):\n%s", diff)
	}
	if diff := testcmp.Diff(got, want); diff != "" {
		t.Errorf("InsertMap() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestYieldCoverage2(t *testing.T) {
	v := slices.All([]int{2, 1, 6})
	var err error
	coverBreak(CatchErr(maps.All(map[int]error{1: nil}), &err))
	coverBreak2(Concat2(v))
	coverBreak2(Filter2(v, predicates.Keys[int, int](predicates.IsEven)))
	coverBreak2(Invert(v))
	coverBreak(ToSeq(v, func(int, int) int { return 0 }))
	coverBreak2(Transform2(v, func(int, int) (int, int) { return 0, 0 }))
	coverBreak2(TransformErr2(v, func(int, int) (int, int, error) { return 0, 0, nil }, &err))
}

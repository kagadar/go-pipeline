package maps

import (
	"cmp"
	"errors"
	"strconv"
	"testing"

	testcmp "github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestKeys(t *testing.T) {
	if diff := testcmp.Diff(
		Keys(map[int]struct{}{1: {}, 2: {}, 3: {}}),
		[]int{1, 2, 3},
		cmpopts.SortSlices(cmp.Less[int]),
	); diff != "" {
		t.Errorf("Keys() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestValues(t *testing.T) {
	if diff := testcmp.Diff(
		Values(map[int]int{1: 9, 2: 8, 3: 7}),
		[]int{9, 8, 7},
		cmpopts.SortSlices(cmp.Less[int]),
	); diff != "" {
		t.Errorf("Values() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestFilter(t *testing.T) {
	if diff := testcmp.Diff(Filter(map[int]int{1: 1, 2: 3, 4: 4}, func(k, v int) bool { return k == v }),
		map[int]int{1: 1, 4: 4}); diff != "" {
		t.Errorf("Filter() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestToSlice(t *testing.T) {
	if diff := testcmp.Diff(
		ToSlice(map[int]struct{}{1: {}, 2: {}, 3: {}}, func(k int, _ struct{}) int { return k }),
		[]int{1, 2, 3},
		cmpopts.SortSlices(cmp.Less[int]),
	); diff != "" {
		t.Errorf("ToSlice() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestMapMapInsert(t *testing.T) {
	got := map[int]map[string]bool{}
	MapMapInsert(got, 1, "2", true)
	MapMapInsert(got, 1, "3", false)
	if diff := testcmp.Diff(
		got,
		map[int]map[string]bool{1: {"2": true, "3": false}},
	); diff != "" {
		t.Errorf("MapMapInsert() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestTransform(t *testing.T) {
	if diff := testcmp.Diff(
		Transform(
			map[int]struct{}{1: {}, 2: {}, 3: {}},
			func(k int, v struct{}) (string, struct{}) { return strconv.Itoa(k), v },
		),
		map[string]struct{}{"1": {}, "2": {}, "3": {}},
	); diff != "" {
		t.Errorf("Transform() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestRange(t *testing.T) {
	var got []string
	if err := Range(
		map[int]string{1: "a", 2: "b", 3: "c"},
		[]int{1, 2},
		func(i int, s string) error {
			got = append(got, s)
			return nil
		},
	); err != nil {
		t.Fatalf("Range() unexpected error: %v", err)
	}
	if diff := testcmp.Diff(got, []string{"a", "b"}); diff != "" {
		t.Errorf("Range() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestRange_Err(t *testing.T) {
	want := errors.New("test error")
	if err := Range(
		map[int]string{1: "a", 2: "b", 3: "c"},
		[]int{1, 2}, func(i int, s string) error {
			return want
		},
	); err != want {
		t.Errorf("Range() unexpected error: got %v want %v", err, want)
	}
}

func TestSortedRange(t *testing.T) {
	var got []string
	if err := SortedRange(
		map[int]string{1: "a", 2: "b", 3: "c"},
		func(i int, s string) error {
			got = append(got, s)
			return nil
		},
	); err != nil {
		t.Fatalf("SortedRange() unexpected error: %v", err)
	}
	if diff := testcmp.Diff(got, []string{"a", "b", "c"}); diff != "" {
		t.Errorf("SortedRange() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestSortedRangeFunc(t *testing.T) {
	type key struct {
		a string
		b int
	}
	var got []string
	if err := SortedRangeFunc(
		map[key]string{{a: "z", b: 1}: "a", {a: "y", b: 2}: "b", {a: "x", b: 3}: "c"},
		func(x, y key) int {
			return x.b - y.b
		},
		func(k key, s string) error {
			got = append(got, s)
			return nil
		},
	); err != nil {
		t.Fatalf("SortedRangeFunc() unexpected error: %v", err)
	}
	if diff := testcmp.Diff(got, []string{"a", "b", "c"}); diff != "" {
		t.Errorf("SortedRangeFunc() unexpected diff (-got +want):\n%s", diff)
	}
}

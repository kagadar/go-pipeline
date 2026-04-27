package maps

import (
	"cmp"
	"errors"
	"strconv"
	"strings"
	"testing"

	testcmp "github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/kagadar/go-pipeline/predicates"
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

func TestAll(t *testing.T) {
	if !All(map[int]struct{}{2: {}, 4: {}, 6: {}}, predicates.Keys[int, struct{}](predicates.IsEven)) {
		t.Error("All({2,4,6}, IsEven) got false, want true")
	}
	if All(map[int]struct{}{1: {}, 2: {}, 6: {}}, predicates.Keys[int, struct{}](predicates.IsEven)) {
		t.Error("All({1,2,6}, IsEven) got true, want false")
	}
}

func TestAny(t *testing.T) {
	if !Any(map[int]struct{}{1: {}, 2: {}, 3: {}}, predicates.Keys[int, struct{}](predicates.IsEven)) {
		t.Error("Any({1,2,3}, IsEven) got false, want true")
	}
	if Any(map[int]struct{}{1: {}, 3: {}, 5: {}}, predicates.Keys[int, struct{}](predicates.IsEven)) {
		t.Error("Any({1,3,5}, IsEven) got true, want false")
	}
}

func TestFilter(t *testing.T) {
	if diff := testcmp.Diff(
		Filter(map[int]int{1: 1, 2: 3, 4: 4}, func(k, v int) bool { return k == v }),
		map[int]int{1: 1, 4: 4},
	); diff != "" {
		t.Errorf("Filter() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestInvert(t *testing.T) {
	if diff := testcmp.Diff(
		Invert(map[string]int{"1": 9, "2": 8}),
		map[int]string{9: "1", 8: "2"},
	); diff != "" {
		t.Errorf(`Invert({"1": 9, "2": 8}) unexpected diff (-got +want):
%s`, diff)
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

func TestPartition(t *testing.T) {
	evens, odds := Partition(map[int]struct{}{1: {}, 2: {}, 3: {}, 4: {}}, predicates.Keys[int, struct{}](predicates.IsEven))
	if ediff, odiff := testcmp.Diff(evens, map[int]struct{}{2: {}, 4: {}}), testcmp.Diff(odds, map[int]struct{}{1: {}, 3: {}}); ediff != "" || odiff != "" {
		t.Errorf("Partition({1,2,3,4}, IsEven) unexpected diff (-got +want):\n%s", strings.Join([]string{ediff, odiff}, "\n"))
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

func TestReduce(t *testing.T) {
	if got, want := Reduce(map[int]struct{}{1: {}, 2: {}, 3: {}}, func(o, k int, v struct{}) int { return o + k }), 6; got != want {
		t.Errorf("Reduce() got %d want %d", got, want)
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

func TestToSlice(t *testing.T) {
	if diff := testcmp.Diff(
		ToSlice(map[int]struct{}{1: {}, 2: {}, 3: {}}, func(k int, _ struct{}) int { return k }),
		[]int{1, 2, 3},
		cmpopts.SortSlices(cmp.Less[int]),
	); diff != "" {
		t.Errorf("ToSlice() unexpected diff (-got +want):\n%s", diff)
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
		t.Errorf("Transform(Itoa) unexpected diff (-got +want):\n%s", diff)
	}
}

func TestTransformErr(t *testing.T) {
	got, err := TransformErr(
		map[int]struct{}{1: {}, 2: {}, 3: {}},
		func(k int, v struct{}) (string, struct{}, error) { return strconv.Itoa(k), v, nil },
	)
	if err != nil {
		t.Fatalf("TransformErr(Itoa) unexpected error: %v", err)
	}
	if diff := testcmp.Diff(got, map[string]struct{}{"1": {}, "2": {}, "3": {}}); diff != "" {
		t.Errorf("TransformErr(Itoa) unexpected diff (-got +want):\n%s", diff)
	}
}

func TestTransformErr_Error(t *testing.T) {
	want := errors.New("tef")
	errors := []error{nil, want}
	got, err := TransformErr(
		map[int]struct{}{1: {}, 2: {}, 3: {}},
		func(k int, v struct{}) (string, struct{}, error) {
			err := errors[0]
			errors = errors[1:]
			return "", v, err
		},
	)
	if got != nil || err != want {
		t.Fatalf("TransformErr(fail) got %v %v, want map[] %v", got, err, want)
	}
}

func TestValueSortedRange(t *testing.T) {
	var got []int
	if err := ValueSortedRange(
		map[int]string{1: "z", 2: "y", 3: "x"},
		func(i int, s string) error {
			got = append(got, i)
			return nil
		},
	); err != nil {
		t.Fatalf("ValueSortedRange() unexpected error: %v", err)
	}
	if diff := testcmp.Diff(got, []int{3, 2, 1}); diff != "" {
		t.Errorf("ValueSortedRange() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestValueSortedRangeFunc(t *testing.T) {
	type value struct {
		a string
		b int
	}
	var got []int
	if err := ValueSortedRangeFunc(
		map[int]value{1: {"z", 9}, 2: {"y", 8}, 3: {"x", 7}},
		func(x, y value) int {
			return x.b - y.b
		},
		func(i int, v value) error {
			got = append(got, i)
			return nil
		},
	); err != nil {
		t.Fatalf("ValueSortedRangeFunc() unexpected error: %v", err)
	}
	if diff := testcmp.Diff(got, []int{3, 2, 1}); diff != "" {
		t.Errorf("ValueSortedRangeFunc() unexpected diff (-got +want):\n%s", diff)
	}
}

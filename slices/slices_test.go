package slices

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	testcmp "github.com/google/go-cmp/cmp"
	"github.com/kagadar/go-pipeline/predicates"
)

func TestAll(t *testing.T) {
	if !All([]int{2, 4, 6}, predicates.IsEven) {
		t.Error("All({2,4,6}, IsEven) got false, want true")
	}
	if All([]int{1, 2, 6}, predicates.IsEven) {
		t.Error("All({2,4,6}, IsEven) got true, want false")
	}
}

func TestCollateMap(t *testing.T) {
	if diff := testcmp.Diff(
		CollateMap([]int{1, 3, 2, 3}, func(e int) (k int, v struct{}) { return e, v }),
		map[int][]struct{}{1: {{}}, 2: {{}}, 3: {{}, {}}},
	); diff != "" {
		t.Errorf("CollateMap() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestDedupe(t *testing.T) {
	if diff := testcmp.Diff(Dedupe([]int{3, 2, 1, 1, 3}), []int{3, 2, 1}); diff != "" {
		t.Errorf("Dedupe() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestFilter(t *testing.T) {
	if diff := testcmp.Diff(
		Filter([]int{1, 2, 3, 4, 5, 6}, predicates.IsEven),
		[]int{2, 4, 6},
	); diff != "" {
		t.Errorf("Filter({1,2,3,4,5,6}, IsEven) unexpected diff (-got +want):\n%s", diff)
	}
}

func TestFirst(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   []int
		want int
	}{
		{
			name: "populated",
			in:   []int{1, 2, 3},
			want: 1,
		},
		{
			name: "empty",
			in:   []int{},
			want: 0,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := First(tt.in); got != tt.want {
				t.Errorf("Last() got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestFlatten(t *testing.T) {
	if diff := testcmp.Diff(
		Flatten([][]int{{1}, {1, 2}, {1, 2, 3}}),
		[]int{1, 1, 2, 1, 2, 3},
	); diff != "" {
		t.Errorf("Flatten() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestGroupBy(t *testing.T) {
	if diff := testcmp.Diff(
		GroupBy([]int{1, 2, 3, 4}, predicates.FindParity),
		map[predicates.Parity][]int{
			predicates.Odd:  {1, 3},
			predicates.Even: {2, 4},
		},
	); diff != "" {
		t.Errorf("GroupBy(Parity) unexpected diff (-got +want):\n%s", diff)
	}
}

func TestLast(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   []int
		want int
	}{
		{
			name: "populated",
			in:   []int{1, 2, 3},
			want: 3,
		},
		{
			name: "empty",
			in:   []int{},
			want: 0,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := Last(tt.in); got != tt.want {
				t.Errorf("Last() got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestPartition(t *testing.T) {
	evens, odds := Partition([]int{1, 2, 3, 4}, predicates.IsEven)
	if ediff, odiff := testcmp.Diff(evens, []int{2, 4}), testcmp.Diff(odds, []int{3, 1}); ediff != "" || odiff != "" {
		t.Errorf("Partition({1,2,3,4}, IsEven) unexpected diff (-got +want):\n%s", strings.Join([]string{ediff, odiff}, "\n"))
	}
}

func TestReduce(t *testing.T) {
	if got, want := Reduce([]int{1, 2, 3}, func(o string, e int) string { return o + strconv.Itoa(e) }), "123"; got != want {
		t.Errorf("Reduce(Itoa) got %s, want %s", got, want)
	}
}

func TestSkipWhile(t *testing.T) {
	if diff := testcmp.Diff(SkipWhile([]int{0, 2, 4, 1, 2}, predicates.IsEven), []int{1, 2}); diff != "" {
		t.Errorf("SkipWhile() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestTakeWhile(t *testing.T) {
	if diff := testcmp.Diff(TakeWhile([]int{0, 2, 4, 1, 2}, predicates.IsEven), []int{0, 2, 4}); diff != "" {
		t.Errorf("TakeWhile() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestToMap(t *testing.T) {
	if diff := testcmp.Diff(
		ToMap([]int{1, 2, 3}, func(e int) (k int, v struct{}) { return e, v }),
		map[int]struct{}{1: {}, 2: {}, 3: {}},
	); diff != "" {
		t.Errorf("ToMap() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestTransform(t *testing.T) {
	if diff := testcmp.Diff(
		Transform([]int{1, 2, 3}, strconv.Itoa),
		[]string{"1", "2", "3"},
	); diff != "" {
		t.Errorf("Transform(Itoa) unexpected diff (-got +want):\n%s", diff)
	}
}

func TestTransformErr(t *testing.T) {
	got, err := TransformErr([]int{1, 2, 3}, func(e int) (string, error) { return strconv.Itoa(e), nil })
	if err != nil {
		t.Fatalf("TransformErr(Itoa) unexpected error: %v", err)
	}
	if diff := testcmp.Diff(got, []string{"1", "2", "3"}); diff != "" {
		t.Errorf("TransformErr(Itoa) unexpected diff (-got +want):\n%s", diff)
	}
}

func TestTransformErr_Error(t *testing.T) {
	want := errors.New("tef")
	errors := []error{nil, want}
	got, err := TransformErr([]int{1, 2, 3}, func(e int) (string, error) {
		err := errors[0]
		errors = errors[1:]
		return "", err
	})
	if got != nil || err != want {
		t.Errorf("TransformErr(fail) got %v %v, want [] %v", got, err, want)
	}
}

func TestZip(t *testing.T) {
	if l := len(Zip[[]int]()); l != 0 {
		t.Errorf("Zip(nothing) got len %d, want 0", l)
	}
	if diff := testcmp.Diff(
		Zip([]int{1, 2, 3, 4, 5}, []int{5, 6, 7, 8, 9}, []int{9, 8}),
		[][]int{{1, 5, 9}, {2, 6, 8}},
	); diff != "" {
		t.Errorf("Zip() unexpected diff (-got +want):\n%s", diff)
	}
}

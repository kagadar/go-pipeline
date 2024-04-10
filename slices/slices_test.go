package slices

import (
	"fmt"
	"strconv"
	"testing"

	testcmp "github.com/google/go-cmp/cmp"
)

func TestDedupe(t *testing.T) {
	if diff := testcmp.Diff(Dedupe([]int{3, 2, 1, 1, 3}), []int{3, 2, 1}); diff != "" {
		t.Errorf("Dedupe() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestFilter(t *testing.T) {
	if diff := testcmp.Diff(
		Filter([]int{1, 2, 3, 4, 5, 6}, func(e int) bool { return e%3 == 0 }),
		[]int{3, 6},
	); diff != "" {
		t.Errorf("Filter() unexpected diff (-got +want):\n%s", diff)
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
				t.Errorf("Last() got %d want %d", got, tt.want)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	if got, want := Reduce([]int{1, 2, 3}, func(o string, e int) string { return o + strconv.Itoa(e) }), "123"; got != want {
		t.Errorf("Reduce() got %s want %s", got, want)
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
		Transform([]int{1, 2, 3}, func(e int) string { return strconv.Itoa(e) }),
		[]string{"1", "2", "3"},
	); diff != "" {
		t.Errorf("Transform() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestTransformErr(t *testing.T) {
	got, err := TransformErr([]int{1, 2, 3}, func(e int) (string, error) { return strconv.Itoa(e), nil })
	if err != nil {
		t.Fatalf("TransformErr() unexpected error: %v", err)
	}
	if diff := testcmp.Diff(got, []string{"1", "2", "3"}); diff != "" {
		t.Errorf("TransformErr() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestTransformErr_Error(t *testing.T) {
	_, err := TransformErr([]int{1, 2}, func(e int) (string, error) { return "", fmt.Errorf("%d", e) })
	if err == nil {
		t.Fatal("TransformErr() expected error but got nil")
	}
	if err.Error() != "1" {
		t.Fatalf("TransformErr() unexpected error: %v", err)
	}
}

func TestZip(t *testing.T) {
	if diff := testcmp.Diff(Zip([]int{1, 2, 3, 4, 5}, []int{5, 6, 7, 8, 9}, []int{9, 8}), [][]int{{1, 5, 9}, {2, 6, 8}}); diff != "" {
		t.Errorf("Zip() unexpected diff (-got +want):\n%s", diff)
	}
}

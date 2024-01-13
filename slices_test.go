package pipeline

import (
	"strconv"
	"testing"

	testcmp "github.com/google/go-cmp/cmp"
)

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

func TestFlatten(t *testing.T) {
	if diff := testcmp.Diff(
		Flatten([][]int{{1}, {1, 2}, {1, 2, 3}}),
		[]int{1, 1, 2, 1, 2, 3},
	); diff != "" {
		t.Errorf("Flatten() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestSliceToMap(t *testing.T) {
	if diff := testcmp.Diff(
		SliceToMap([]int{1, 2, 3}, func(e int) (k int, v struct{}) { return e, v }),
		map[int]struct{}{1: {}, 2: {}, 3: {}},
	); diff != "" {
		t.Errorf("SliceToMap() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestTransformSlice(t *testing.T) {
	if diff := testcmp.Diff(
		TransformSlice([]int{1, 2, 3}, func(e int) string { return strconv.Itoa(e) }),
		[]string{"1", "2", "3"},
	); diff != "" {
		t.Errorf("TransformSlice() unexpected diff (-got +want):\n%s", diff)
	}
}

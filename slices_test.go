package pipeline

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFlatten(t *testing.T) {
	if diff := cmp.Diff(
		Flatten([][]int{{1}, {1, 2}, {1, 2, 3}}),
		[]int{1, 1, 2, 1, 2, 3},
	); diff != "" {
		t.Errorf("Flatten() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestLess(t *testing.T) {
	if Less(2, 1) {
		t.Errorf("Less(2, 1) returned true, expected false")
	}
}

func TestSliceToMap(t *testing.T) {
	if diff := cmp.Diff(
		SliceToMap([]int{1, 2, 3}, func(e int) (k int, v struct{}) { return e, v }),
		map[int]struct{}{1: {}, 2: {}, 3: {}},
	); diff != "" {
		t.Errorf("SliceToMap() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestTransformSlice(t *testing.T) {
	if diff := cmp.Diff(
		TransformSlice([]int{1, 2, 3}, func(e int) string { return strconv.Itoa(e) }),
		[]string{"1", "2", "3"},
	); diff != "" {
		t.Errorf("TransformSlice() unexpected diff (-got +want):\n%s", diff)
	}
}

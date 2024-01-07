package pipeline

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestMapToSlice(t *testing.T) {
	if diff := cmp.Diff(
		MapToSlice(map[int]struct{}{1: {}, 2: {}, 3: {}}, func(k int, _ struct{}) int { return k }),
		[]int{1, 2, 3},
		cmpopts.SortSlices(Less[int]),
	); diff != "" {
		t.Errorf("MapToSlice() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestTransformMap(t *testing.T) {
	if diff := cmp.Diff(
		TransformMap(map[int]struct{}{1: {}, 2: {}, 3: {}}, func(k int, v struct{}) (string, struct{}) { return strconv.Itoa(k), v }),
		map[string]struct{}{"1": {}, "2": {}, "3": {}},
	); diff != "" {
		t.Errorf("TransformMap() unexpected diff (-got +want):\n%s", diff)
	}
}

package mapchains

import (
	"strconv"
	"testing"

	testcmp "github.com/google/go-cmp/cmp"
	"github.com/kagadar/go-pipeline/chains"
)

func TestTransform(t *testing.T) {
	got, err := Transform(
		Transform(
			chains.NewMapLink(map[int]struct{}{1: {}, 2: {}, 3: {}}),
			func(k int, v struct{}) (string, int, error) { return strconv.Itoa(k), k, nil },
		),
		func(k string, v int) (int, string, error) { return v, k, nil },
	)()
	if err != nil {
		t.Fatalf("Transform() unexpected error: %v", err)
	}
	if diff := testcmp.Diff(
		got,
		map[int]string{1: "1", 2: "2", 3: "3"},
	); diff != "" {
		t.Errorf("Transform() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestTransform_ErrKeyCollision(t *testing.T) {
	_, err := Transform(
			chains.NewMapLink(map[int]string{1: "1", 2: "1", 3: "3"}),
			func(k int, v string) (string, int, error) { return v, k, nil },
		)()
	if !IsKeyCollisionErr(err) {
		t.Errorf("Transform() error = %v, want %v", err, errKeyCollision)
	}
}

func TestFilter(t *testing.T) {
	got, err := Filter(chains.NewMapLink(map[int]int{1: 1, 2: 3, 4: 4}), func(k, v int) bool { return k == v })()
	if err != nil {
		t.Fatalf("Filter() unexpected error: %v", err)
	}
	if diff := testcmp.Diff(
		got,
		map[int]int{1: 1, 4: 4},
	); diff != "" {
		t.Errorf("Filter() unexpected diff (-got +want):\n%s", diff)
	}
}

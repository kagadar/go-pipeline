package api

import (
	"errors"
	"testing"

	testcmp "github.com/google/go-cmp/cmp"
)

func data(last int, pageSize int) []int {
	d := []int{1, 2, 3, 4, 5, 6}
	if last >= len(d) {
		return []int{}
	}
	return d[last:min(last+pageSize, len(d))]
}

func TestPaginate(t *testing.T) {
	got, err := Paginate(func(i int) ([]int, error) {
		d := data(i, 3)
		return d, nil
	})
	if err != nil {
		t.Fatalf("Paginate() unexpected error: %v", err)
	}
	if diff := testcmp.Diff(got, []int{1, 2, 3, 4, 5, 6}); diff != "" {
		t.Errorf("Paginate() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestPaginate_Err(t *testing.T) {
	want := errors.New("test error")
	_, err := Paginate(func(i int) ([]int, error) {
		return nil, want
	})
	if err != want {
		t.Errorf("Paginate() unexpected error: got %v want %v", err, want)
	}
}

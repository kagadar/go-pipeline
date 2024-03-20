package api

import (
	"errors"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const pageSize = 3

var data = []int{1, 2, 3, 4, 5, 6}

func getDataFromIdx(idx int) ([]int, int, error) {
	if idx >= len(data) {
		return nil, 0, nil
	}
	end := min(idx+pageSize, len(data))
	return data[idx:end], end, nil
}

func TestPaginateToken(t *testing.T) {
	got, err := PaginateToken(getDataFromIdx)
	if err != nil {
		t.Fatalf("PaginateToken() unexpected error: %v", err)
	}
	if diff := cmp.Diff(got, data); diff != "" {
		t.Errorf("PaginateToken() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestPaginateToken_Error(t *testing.T) {
	want := errors.New("test error")
	if _, err := PaginateToken(func(i int) ([]int, int, error) {
		return nil, 0, want
	}); err != want {
		t.Errorf("PaginateToken() unexpected error: got %v want %v", err, want)
	}
}

func TestPaginate158(t *testing.T) {
	got, err := Paginate158(func(token string) ([]int, string, error) {
		var idx int
		if token != "" {
			var err error
			idx, err = strconv.Atoi(token)
			if err != nil {
				return nil, "", err
			}
		}
		var next string
		d, ni, err := getDataFromIdx(idx)
		if len(d) != 0 {
			next = strconv.Itoa(ni)
		}
		return d, next, err
	})
	if err != nil {
		t.Fatalf("Paginate158() unexpected error: %v", err)
	}
	if diff := cmp.Diff(got, data); diff != "" {
		t.Errorf("Paginate158() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestPaginate158_Error(t *testing.T) {
	want := errors.New("test error")
	if _, err := Paginate158(func(t string) ([]int, string, error) {
		return nil, "", want
	}); err != want {
		t.Errorf("Paginate158() unexpected error: got %v want %v", err, want)
	}
}

func TestPaginate(t *testing.T) {
	got, err := Paginate(func(val int) ([]int, error) {
		for i, v := range data {
			if v > val {
				return data[i:min(i+pageSize, len(data))], nil
			}
		}
		return nil, nil
	})
	if err != nil {
		t.Fatalf("Paginate() unexpected error: %v", err)
	}
	if diff := cmp.Diff(got, data); diff != "" {
		t.Errorf("Paginate() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestPaginate_Error(t *testing.T) {
	want := errors.New("test error")
	if _, err := Paginate(func(i int) ([]int, error) {
		return nil, want
	}); err != want {
		t.Errorf("Paginate() unexpected error: got %v want %v", err, want)
	}
}

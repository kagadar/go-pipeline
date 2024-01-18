package maps

import (
	"cmp"
	"testing"

	testcmp "github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestDifference(t *testing.T) {
	if diff := testcmp.Diff(
		Difference(
			map[int]struct{}{1: {}, 2: {}, 3: {}},
			map[int]struct{}{1: {}, 2: {}},
		),
		map[int]struct{}{3: {}},
	); diff != "" {
		t.Errorf("Difference() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestDisjoint(t *testing.T) {
	if diff := testcmp.Diff(
		Disjoint(
			map[int]struct{}{1: {}, 2: {}, 3: {}},
			map[int]struct{}{2: {}, 3: {}, 4: {}},
		),
		map[int]struct{}{1: {}, 4: {}},
	); diff != "" {
		t.Errorf("Disjoint() unexpected diff (-got +want):\n%s", diff)
	}
}

func TestIntersect(t *testing.T) {
	for _, tt := range []struct {
		name string
		a, b map[int]struct{}
		want []int
	}{
		{
			name: "less a",
			a:    map[int]struct{}{1: {}, 2: {}},
			b:    map[int]struct{}{2: {}, 3: {}, 4: {}},
			want: []int{2},
		},
		{
			name: "less b",
			a:    map[int]struct{}{1: {}, 2: {}, 3: {}},
			b:    map[int]struct{}{2: {}, 3: {}},
			want: []int{2, 3},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := testcmp.Diff(
				Intersect(
					tt.a,
					tt.b,
				),
				tt.want,
				cmpopts.SortSlices(cmp.Less[int]),
			); diff != "" {
				t.Errorf("Intersect() unexpected diff (-got +want):\n%s", diff)
			}
		})
	}
}

func TestSubset(t *testing.T) {
	for _, tt := range []struct {
		name string
		a, b map[int]struct{}
		want bool
	}{
		{
			name: "true",
			a:    map[int]struct{}{1: {}, 2: {}, 3: {}},
			b:    map[int]struct{}{1: {}, 2: {}, 3: {}, 4: {}},
			want: true,
		},
		{
			name: "false",
			a:    map[int]struct{}{1: {}, 2: {}, 3: {}},
			b:    map[int]struct{}{2: {}, 3: {}, 4: {}},
			want: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := Subset(tt.a, tt.b); got != tt.want {
				t.Errorf("Subset() got %t want %t", got, tt.want)
			}
		})
	}
}

func TestUnion(t *testing.T) {
	if diff := testcmp.Diff(
		Union(
			map[int]string{1: "a", 2: "b", 3: "c"},
			map[int]string{2: "x", 3: "y", 4: "z"},
		),
		map[int]string{1: "a", 2: "x", 3: "y", 4: "z"},
	); diff != "" {
		t.Errorf("Union() unexpected diff (-got +want):\n%s", diff)
	}
}

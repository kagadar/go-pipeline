package maps

import (
	"testing"

	testcmp "github.com/google/go-cmp/cmp"
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
		name  string
		input []map[int]string
		want  map[int][]string
	}{
		{
			name: "empty",
			want: nil,
		},
		{
			name:  "less first",
			input: []map[int]string{{1: "a", 2: "b"}, {2: "c", 3: "d", 4: "e"}},
			want:  map[int][]string{2: {"b", "c"}},
		},
		{
			name:  "less last",
			input: []map[int]string{{1: "a", 2: "b", 3: "c"}, {2: "d", 3: "e", 4: "f", 5: "g"}, {2: "h", 3: "i"}},
			want:  map[int][]string{2: {"b", "d", "h"}, 3: {"c", "e", "i"}},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if diff := testcmp.Diff(Intersect(tt.input...), tt.want); diff != "" {
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
			map[int]string{2: "i", 3: "j", 4: "k"},
			map[int]string{3: "x", 4: "y", 5: "z"},
		),
		map[int][]string{1: {"a"}, 2: {"b", "i"}, 3: {"c", "j", "x"}, 4: {"k", "y"}, 5: {"z"}},
	); diff != "" {
		t.Errorf("Union() unexpected diff (-got +want):\n%s", diff)
	}
}

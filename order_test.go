package pipeline

import "testing"

func TestMax(t *testing.T) {
	for _, tt := range []struct {
		name       string
		x, y, want int
	}{
		{name: "x", x: 2, y: 1, want: 2},
		{name: "y", x: 1, y: 2, want: 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := Max(tt.x, tt.y); got != tt.want {
				t.Errorf("Max() got %d want %d", got, tt.want)
			}
		})
	}
}

func TestMin(t *testing.T) {
	for _, tt := range []struct {
		name       string
		x, y, want int
	}{
		{name: "x", x: 1, y: 2, want: 1},
		{name: "y", x: 2, y: 1, want: 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.x, tt.y); got != tt.want {
				t.Errorf("Min() got %d want %d", got, tt.want)
			}
		})
	}
}

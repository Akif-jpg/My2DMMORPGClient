package shapes

import (
	"testing"
)

func TestRectangleIntersects_Detailed(t *testing.T) {
	tests := []struct {
		name     string
		r1       Rectangle
		r2       Rectangle
		expected bool
	}{
		{
			name:     "Separated on Y axis only",
			r1:       Rectangle{Point{0, 0}, 4, 2},
			r2:       Rectangle{Point{0, 5}, 4, 2},
			expected: false,
		},
		{
			name:     "Separated on X axis only",
			r1:       Rectangle{Point{0, 0}, 2, 4},
			r2:       Rectangle{Point{5, 0}, 2, 4},
			expected: false,
		},
		{
			name:     "Asymmetric overlap Y axis",
			r1:       Rectangle{Point{0, 10}, 2, 2},
			r2:       Rectangle{Point{0, 0}, 2, 20},
			expected: true,
		},
		{
			name:     "Tall rectangle intersects short wide rectangle",
			r1:       Rectangle{Point{0, 0}, 2, 10},
			r2:       Rectangle{Point{0, 4}, 6, 2},
			expected: true,
		},
		{
			name:     "Touching only at corner",
			r1:       Rectangle{Point{0, 0}, 2, 2},
			r2:       Rectangle{Point{2, 2}, 2, 2},
			expected: false,
		},
		{
			name:     "One inside other but shifted on Y",
			r1:       Rectangle{Point{0, 0}, 10, 10},
			r2:       Rectangle{Point{0, 3}, 2, 2},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r1.IntersectsRectangle(tt.r2); got != tt.expected {
				t.Errorf("Intersects() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestRectangleContains_Detailed(t *testing.T) {
	tests := []struct {
		name     string
		outer    Rectangle
		inner    Rectangle
		expected bool
	}{
		{
			name:     "Fully contained centered",
			outer:    Rectangle{Point{0, 0}, 10, 10},
			inner:    Rectangle{Point{0, 0}, 2, 2},
			expected: true,
		},
		{
			name:     "Contained but shifted on Y",
			outer:    Rectangle{Point{0, 0}, 10, 10},
			inner:    Rectangle{Point{0, 3}, 2, 2},
			expected: true,
		},
		{
			name:     "Exceeds on Y axis only",
			outer:    Rectangle{Point{0, 0}, 10, 10},
			inner:    Rectangle{Point{0, 6}, 2, 2},
			expected: false,
		},
		{
			name:     "Exceeds on X axis only",
			outer:    Rectangle{Point{0, 0}, 10, 10},
			inner:    Rectangle{Point{6, 0}, 2, 2},
			expected: false,
		},
		{
			name:     "Tall inner exceeds height but not width",
			outer:    Rectangle{Point{0, 0}, 10, 6},
			inner:    Rectangle{Point{0, 0}, 2, 8},
			expected: false,
		},
		{
			name:     "Wide inner exceeds width but not height",
			outer:    Rectangle{Point{0, 0}, 6, 10},
			inner:    Rectangle{Point{0, 0}, 8, 2},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.outer.ContainsRectangle(tt.inner); got != tt.expected {
				t.Errorf("Contains() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

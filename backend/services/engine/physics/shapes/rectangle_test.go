package shapes

import (
	"testing"
)

func Test_MyStruct_Implements_MyInterface_Runtime(t *testing.T) {
	var s any = NewRectangle(*NewPoint(0, 0), 0, 0)

	if _, ok := s.(Shape); !ok {
		t.Fatalf("Rectangle does not implement Shape")
	}
}

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
		{
			name:     "One inside other but shifted on X",
			r1:       Rectangle{Point{0, 0}, 10, 10},
			r2:       Rectangle{Point{3, 0}, 2, 2},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r1.IntersectsRectangle(&tt.r2); got != tt.expected {
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
			if got := tt.outer.ContainsRectangle(&tt.inner); got != tt.expected {
				t.Errorf("Contains() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestRectangle_IntersectsCircle(t *testing.T) {
	tests := []struct {
		name     string
		rect     Rectangle
		circle   Circle
		expected bool
	}{
		{
			name:     "Circle overlaps rectangle center",
			rect:     Rectangle{Point{0, 0}, 10, 10},
			circle:   Circle{Point{0, 0}, 3},
			expected: true,
		},
		{
			name:     "Circle touches rectangle edge",
			rect:     Rectangle{Point{0, 0}, 10, 10},
			circle:   Circle{Point{8, 0}, 3},
			expected: true,
		},
		{
			name:     "Circle touches rectangle corner",
			rect:     Rectangle{Point{0, 0}, 10, 10},
			circle:   Circle{Point{7, 7}, 3},
			expected: true,
		},
		{
			name:     "Circle completely outside rectangle",
			rect:     Rectangle{Point{0, 0}, 10, 10},
			circle:   Circle{Point{20, 20}, 3},
			expected: false,
		},
		{
			name:     "Circle completely inside rectangle",
			rect:     Rectangle{Point{0, 0}, 20, 20},
			circle:   Circle{Point{0, 0}, 3},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rect.IntersectsCircle(&tt.circle); got != tt.expected {
				t.Errorf("IntersectsCircle() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestRectangle_ContainsCircle(t *testing.T) {
	tests := []struct {
		name     string
		rect     Rectangle
		circle   Circle
		expected bool
	}{
		{
			name:     "Rectangle fully contains small circle",
			rect:     Rectangle{Point{0, 0}, 20, 20},
			circle:   Circle{Point{0, 0}, 3},
			expected: true,
		},
		{
			name:     "Circle touches rectangle edge from inside",
			rect:     Rectangle{Point{0, 0}, 10, 10},
			circle:   Circle{Point{0, 0}, 5},
			expected: true,
		},
		{
			name:     "Circle extends beyond rectangle",
			rect:     Rectangle{Point{0, 0}, 10, 10},
			circle:   Circle{Point{0, 0}, 6},
			expected: false,
		},
		{
			name:     "Circle completely outside rectangle",
			rect:     Rectangle{Point{0, 0}, 10, 10},
			circle:   Circle{Point{20, 20}, 3},
			expected: false,
		},
		{
			name:     "Circle at corner extends beyond",
			rect:     Rectangle{Point{0, 0}, 10, 10},
			circle:   Circle{Point{3, 3}, 4},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rect.ContainsCircle(&tt.circle); got != tt.expected {
				t.Errorf("ContainsCircle() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestRectangle_ContainsPoint(t *testing.T) {
	tests := []struct {
		name     string
		rect     Rectangle
		point    Point
		expected bool
	}{
		{
			name:     "Point at center",
			rect:     Rectangle{Point{0, 0}, 10, 10},
			point:    Point{0, 0},
			expected: true,
		},
		{
			name:     "Point at edge (within bounds)",
			rect:     Rectangle{Point{0, 0}, 10, 10},
			point:    Point{5, 0},
			expected: true,
		},
		{
			name:     "Point at corner (within bounds)",
			rect:     Rectangle{Point{0, 0}, 10, 10},
			point:    Point{5, 5},
			expected: true,
		},
		{
			name:     "Point outside on X axis",
			rect:     Rectangle{Point{0, 0}, 10, 10},
			point:    Point{10, 0},
			expected: false,
		},
		{
			name:     "Point outside on Y axis",
			rect:     Rectangle{Point{0, 0}, 10, 10},
			point:    Point{0, 10},
			expected: false,
		},
		{
			name:     "Point completely outside",
			rect:     Rectangle{Point{0, 0}, 10, 10},
			point:    Point{20, 20},
			expected: false,
		},
		{
			name:     "Point inside shifted rectangle",
			rect:     Rectangle{Point{5, 5}, 10, 10},
			point:    Point{7, 7},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rect.ContainsPoint(&tt.point); got != tt.expected {
				t.Errorf("ContainsPoint() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

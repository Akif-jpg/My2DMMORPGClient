package geometry

import (
	"testing"
)

func TestCircle_Implements_Shape_Interface(t *testing.T) {
	var s any = &Circle{Center: Point{0, 0}, Radius: 5}

	if _, ok := s.(Shape); !ok {
		t.Fatalf("Circle does not implement Shape interface")
	}
}

func TestCircle_IntersectsLine(t *testing.T) {
	tests := []struct {
		name     string
		circle   Circle
		line     Line
		expected bool
	}{
		{
			name:     "Line passes through circle center",
			circle:   Circle{Center: Point{0, 0}, Radius: 5},
			line:     Line{Start: Point{-10, 0}, End: Point{10, 0}},
			expected: true,
		},
		{
			name:     "Line touches circle edge",
			circle:   Circle{Center: Point{0, 0}, Radius: 5},
			line:     Line{Start: Point{-10, 5}, End: Point{10, 5}},
			expected: true,
		},
		{
			name:     "Line outside circle",
			circle:   Circle{Center: Point{0, 0}, Radius: 5},
			line:     Line{Start: Point{-10, 10}, End: Point{10, 10}},
			expected: false,
		},
		{
			name:     "Line segment starts inside circle",
			circle:   Circle{Center: Point{0, 0}, Radius: 5},
			line:     Line{Start: Point{0, 0}, End: Point{10, 10}},
			expected: true,
		},
		{
			name:     "Short line segment inside circle",
			circle:   Circle{Center: Point{0, 0}, Radius: 10},
			line:     Line{Start: Point{1, 1}, End: Point{2, 2}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.circle.IntersectsLine(&tt.line); got != tt.expected {
				t.Errorf("IntersectsLine() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestCircle_IntersectsPoint(t *testing.T) {
	tests := []struct {
		name     string
		circle   Circle
		point    Point
		expected bool
	}{
		{
			name:     "Point at center",
			circle:   Circle{Center: Point{0, 0}, Radius: 5},
			point:    Point{0, 0},
			expected: true,
		},
		{
			name:     "Point on edge",
			circle:   Circle{Center: Point{0, 0}, Radius: 5},
			point:    Point{5, 0},
			expected: true,
		},
		{
			name:     "Point outside circle",
			circle:   Circle{Center: Point{0, 0}, Radius: 5},
			point:    Point{10, 10},
			expected: false,
		},
		{
			name:     "Point inside circle",
			circle:   Circle{Center: Point{5, 5}, Radius: 10},
			point:    Point{7, 7},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.circle.IntersectsPoint(&tt.point); got != tt.expected {
				t.Errorf("IntersectsPoint() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestCircle_ContainsLine(t *testing.T) {
	tests := []struct {
		name     string
		circle   Circle
		line     Line
		expected bool
	}{
		{
			name:     "Line fully inside circle",
			circle:   Circle{Center: Point{0, 0}, Radius: 10},
			line:     Line{Start: Point{1, 1}, End: Point{2, 2}},
			expected: true,
		},
		{
			name:     "Line partially outside circle",
			circle:   Circle{Center: Point{0, 0}, Radius: 5},
			line:     Line{Start: Point{0, 0}, End: Point{10, 0}},
			expected: false,
		},
		{
			name:     "Line completely outside circle",
			circle:   Circle{Center: Point{0, 0}, Radius: 5},
			line:     Line{Start: Point{10, 10}, End: Point{20, 20}},
			expected: false,
		},
		{
			name:     "Line with both endpoints on circle edge",
			circle:   Circle{Center: Point{0, 0}, Radius: 5},
			line:     Line{Start: Point{5, 0}, End: Point{0, 5}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.circle.ContainsLine(&tt.line); got != tt.expected {
				t.Errorf("ContainsLine() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestCircle_IntersectsCircle(t *testing.T) {
	tests := []struct {
		name     string
		c1       Circle
		c2       Circle
		expected bool
	}{
		{
			name:     "Overlapping circles",
			c1:       Circle{Center: Point{0, 0}, Radius: 5},
			c2:       Circle{Center: Point{3, 0}, Radius: 5},
			expected: true,
		},
		{
			name:     "Touching circles",
			c1:       Circle{Center: Point{0, 0}, Radius: 5},
			c2:       Circle{Center: Point{10, 0}, Radius: 5},
			expected: true,
		},
		{
			name:     "Separated circles",
			c1:       Circle{Center: Point{0, 0}, Radius: 5},
			c2:       Circle{Center: Point{20, 0}, Radius: 5},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c1.IntersectsCircle(&tt.c2); got != tt.expected {
				t.Errorf("IntersectsCircle() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestCircle_IntersectsRectangle(t *testing.T) {
	tests := []struct {
		name     string
		circle   Circle
		rect     Rectangle
		expected bool
	}{
		{
			name:     "Circle overlaps rectangle",
			circle:   Circle{Center: Point{0, 0}, Radius: 5},
			rect:     Rectangle{Center: Point{3, 0}, Width: 4, Height: 4},
			expected: true,
		},
		{
			name:     "Circle and rectangle separated",
			circle:   Circle{Center: Point{0, 0}, Radius: 2},
			rect:     Rectangle{Center: Point{10, 10}, Width: 4, Height: 4},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.circle.IntersectsRectangle(&tt.rect); got != tt.expected {
				t.Errorf("IntersectsRectangle() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

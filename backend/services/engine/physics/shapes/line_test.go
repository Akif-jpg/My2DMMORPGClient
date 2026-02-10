package shapes

import (
	"math"
	"testing"
)

func TestLine_Implements_Shape_Interface(t *testing.T) {
	var s any = &Line{Start: Point{0, 0}, End: Point{10, 10}}

	if _, ok := s.(Shape); !ok {
		t.Fatalf("Line does not implement Shape interface")
	}
}

func TestLine_GetType(t *testing.T) {
	line := Line{Start: Point{0, 0}, End: Point{10, 10}}
	if got := line.GetType(); got != "Line" {
		t.Errorf("GetType() = %v, expected 'Line'", got)
	}
}

func TestLine_GetCenter(t *testing.T) {
	tests := []struct {
		name     string
		line     Line
		expected Point
	}{
		{
			name:     "Horizontal line",
			line:     Line{Start: Point{0, 0}, End: Point{10, 0}},
			expected: Point{5, 0},
		},
		{
			name:     "Vertical line",
			line:     Line{Start: Point{0, 0}, End: Point{0, 10}},
			expected: Point{0, 5},
		},
		{
			name:     "Diagonal line",
			line:     Line{Start: Point{0, 0}, End: Point{10, 10}},
			expected: Point{5, 5},
		},
		{
			name:     "Negative coordinates",
			line:     Line{Start: Point{-10, -10}, End: Point{10, 10}},
			expected: Point{0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.line.GetCenter()
			if got.X != tt.expected.X || got.Y != tt.expected.Y {
				t.Errorf("GetCenter() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestLine_IntersectsRectangle(t *testing.T) {
	tests := []struct {
		name     string
		line     Line
		rect     Rectangle
		expected bool
	}{
		{
			name:     "Line passes through rectangle",
			line:     Line{Start: Point{-10, 0}, End: Point{10, 0}},
			rect:     Rectangle{Center: Point{0, 0}, Width: 4, Height: 4},
			expected: true,
		},
		{
			name:     "Line outside rectangle",
			line:     Line{Start: Point{10, 10}, End: Point{20, 20}},
			rect:     Rectangle{Center: Point{0, 0}, Width: 4, Height: 4},
			expected: false,
		},
		{
			name:     "Line starts inside rectangle",
			line:     Line{Start: Point{0, 0}, End: Point{10, 10}},
			rect:     Rectangle{Center: Point{0, 0}, Width: 4, Height: 4},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.line.IntersectsRectangle(&tt.rect); got != tt.expected {
				t.Errorf("IntersectsRectangle() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestLine_IntersectsCircle(t *testing.T) {
	tests := []struct {
		name     string
		line     Line
		circle   Circle
		expected bool
	}{
		{
			name:     "Line passes through circle",
			line:     Line{Start: Point{-10, 0}, End: Point{10, 0}},
			circle:   Circle{Center: Point{0, 0}, Radius: 5},
			expected: true,
		},
		{
			name:     "Line outside circle",
			line:     Line{Start: Point{10, 10}, End: Point{20, 20}},
			circle:   Circle{Center: Point{0, 0}, Radius: 5},
			expected: false,
		},
		{
			name:     "Line touches circle edge",
			line:     Line{Start: Point{-10, 5}, End: Point{10, 5}},
			circle:   Circle{Center: Point{0, 0}, Radius: 5},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.line.IntersectsCircle(&tt.circle); got != tt.expected {
				t.Errorf("IntersectsCircle() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestLine_IntersectsPoint(t *testing.T) {
	tests := []struct {
		name     string
		line     Line
		point    Point
		expected bool
	}{
		{
			name:     "Point on line",
			line:     Line{Start: Point{0, 0}, End: Point{10, 10}},
			point:    Point{5, 5},
			expected: true,
		},
		{
			name:     "Point at start",
			line:     Line{Start: Point{0, 0}, End: Point{10, 10}},
			point:    Point{0, 0},
			expected: true,
		},
		{
			name:     "Point at end",
			line:     Line{Start: Point{0, 0}, End: Point{10, 10}},
			point:    Point{10, 10},
			expected: true,
		},
		{
			name:     "Point off line",
			line:     Line{Start: Point{0, 0}, End: Point{10, 10}},
			point:    Point{5, 0},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.line.IntersectsPoint(&tt.point); got != tt.expected {
				t.Errorf("IntersectsPoint() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestLine_ContainsRectangle(t *testing.T) {
	line := Line{Start: Point{0, 0}, End: Point{10, 10}}
	rect := Rectangle{Center: Point{5, 5}, Width: 2, Height: 2}

	// A line can never contain a 2D shape
	if got := line.ContainsRectangle(&rect); got != false {
		t.Errorf("ContainsRectangle() = %v, expected false", got)
	}
}

func TestLine_ContainsCircle(t *testing.T) {
	line := Line{Start: Point{0, 0}, End: Point{10, 10}}
	circle := Circle{Center: Point{5, 5}, Radius: 2}

	// A line can never contain a 2D shape
	if got := line.ContainsCircle(&circle); got != false {
		t.Errorf("ContainsCircle() = %v, expected false", got)
	}
}

func TestLine_ContainsLine(t *testing.T) {
	tests := []struct {
		name     string
		line1    Line
		line2    Line
		expected bool
	}{
		{
			name:     "Line contains smaller collinear line",
			line1:    Line{Start: Point{0, 0}, End: Point{10, 10}},
			line2:    Line{Start: Point{2, 2}, End: Point{5, 5}},
			expected: true,
		},
		{
			name:     "Lines are identical",
			line1:    Line{Start: Point{0, 0}, End: Point{10, 10}},
			line2:    Line{Start: Point{0, 0}, End: Point{10, 10}},
			expected: true,
		},
		{
			name:     "Lines are not collinear",
			line1:    Line{Start: Point{0, 0}, End: Point{10, 10}},
			line2:    Line{Start: Point{0, 0}, End: Point{10, 0}},
			expected: false,
		},
		{
			name:     "Line extends beyond container",
			line1:    Line{Start: Point{0, 0}, End: Point{5, 5}},
			line2:    Line{Start: Point{0, 0}, End: Point{10, 10}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.line1.ContainsLine(&tt.line2); got != tt.expected {
				t.Errorf("ContainsLine() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestLine_ContainsPoint(t *testing.T) {
	tests := []struct {
		name     string
		line     Line
		point    Point
		expected bool
	}{
		{
			name:     "Point on horizontal line",
			line:     Line{Start: Point{0, 5}, End: Point{10, 5}},
			point:    Point{5, 5},
			expected: true,
		},
		{
			name:     "Point on vertical line",
			line:     Line{Start: Point{5, 0}, End: Point{5, 10}},
			point:    Point{5, 5},
			expected: true,
		},
		{
			name:     "Point on diagonal line",
			line:     Line{Start: Point{0, 0}, End: Point{10, 10}},
			point:    Point{5, 5},
			expected: true,
		},
		{
			name:     "Point not on line (parallel)",
			line:     Line{Start: Point{0, 0}, End: Point{10, 10}},
			point:    Point{5, 6},
			expected: false,
		},
		{
			name:     "Point beyond line segment",
			line:     Line{Start: Point{0, 0}, End: Point{5, 5}},
			point:    Point{10, 10},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.line.ContainsPoint(&tt.point)
			if got != tt.expected {
				t.Errorf("ContainsPoint() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestLine_Length(t *testing.T) {
	tests := []struct {
		name     string
		line     Line
		expected float64
	}{
		{
			name:     "Horizontal line",
			line:     Line{Start: Point{0, 0}, End: Point{10, 0}},
			expected: 10.0,
		},
		{
			name:     "Vertical line",
			line:     Line{Start: Point{0, 0}, End: Point{0, 10}},
			expected: 10.0,
		},
		{
			name:     "Diagonal line (3-4-5 triangle)",
			line:     Line{Start: Point{0, 0}, End: Point{3, 4}},
			expected: 5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.line.Length()
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("Length() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

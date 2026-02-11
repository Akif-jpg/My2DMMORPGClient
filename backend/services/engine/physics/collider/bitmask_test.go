package collider

import "testing"

func TestBitmask_SetBit(t *testing.T) {
	bm := NewBitmask()
	bm.SetBit(0)
	bm.SetBit(5)

	if !bm.IsSet(0) {
		t.Error("Bit 0 set edilmedi")
	}
	if !bm.IsSet(5) {
		t.Error("Bit 5 set edilmedi")
	}
	if bm.IsSet(3) {
		t.Error("Bit 3 set edilmemeli")
	}
}

func TestBitmask_ClearBit(t *testing.T) {
	bm := NewBitmask()
	bm.SetBit(2)
	bm.ClearBit(2)

	if bm.IsSet(2) {
		t.Error("Bit 2 temizlenmedi")
	}
}

func TestBitmask_CanMatch(t *testing.T) {
	tests := []struct {
		name     string
		mask1    []uint32
		mask2    []uint32
		expected bool
	}{
		{
			name:     "Aynı layer - match olmalı",
			mask1:    []uint32{0, 1},
			mask2:    []uint32{1, 2},
			expected: true,
		},
		{
			name:     "Farklı layer - match olmamalı",
			mask1:    []uint32{0},
			mask2:    []uint32{1},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm1 := NewBitmask()
			bm2 := NewBitmask()

			for _, bit := range tt.mask1 {
				bm1.SetBit(bit)
			}
			for _, bit := range tt.mask2 {
				bm2.SetBit(bit)
			}

			if got := bm1.CanMatch(bm2); got != tt.expected {
				t.Errorf("CanMatch() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBitmask_EdgeCases(t *testing.T) {
	bm := NewBitmask()

	// 31. bit (maksimum)
	bm.SetBit(31)
	if !bm.IsSet(31) {
		t.Error("31. bit set edilemedi")
	}

	// Overflow testi (32 ve üzeri)
	// Bu durumda ne olacağına karar verin
}

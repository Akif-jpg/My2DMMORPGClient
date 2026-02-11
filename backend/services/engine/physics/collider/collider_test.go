package collider

import (
	"testing"

	"github.com/Akif-jpg/MyHobieMMORPGGame/services/engine/physics/geometry"
)

func TestCollider_LayerMatching(t *testing.T) {
	// Player collider: Layer 0, Enemy ve Wall ile çarpışabilir
	player := &Collider{
		ShapeList: []geometry.Shape{&geometry.Circle{Radius: 10}},
	}
	player.LayerMask.SetBit(0) // Player layer
	player.MatchMask.SetBit(1) // Enemy ile çarpışabilir
	player.MatchMask.SetBit(3) // Wall ile çarpışabilir

	// Enemy collider: Layer 1, Player ve Projectile ile çarpışabilir
	enemy := &Collider{
		ShapeList: []geometry.Shape{&geometry.Circle{Radius: 10}},
	}
	enemy.LayerMask.SetBit(1) // Enemy layer
	enemy.MatchMask.SetBit(0) // Player ile çarpışabilir
	enemy.MatchMask.SetBit(2) // Projectile ile çarpışabilir ← BURAYI EKLEDİK

	// Projectile: Layer 2, sadece Enemy ile çarpışabilir
	projectile := &Collider{
		ShapeList: []geometry.Shape{&geometry.Circle{Radius: 5}},
	}
	projectile.LayerMask.SetBit(2)
	projectile.MatchMask.SetBit(1) // Enemy ile çarpışabilir

	tests := []struct {
		name        string
		collider1   *Collider
		collider2   *Collider
		shouldMatch bool
	}{
		{"Player-Enemy match", player, enemy, true},
		{"Player-Projectile no match", player, projectile, false},
		{"Enemy-Projectile match", enemy, projectile, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// İki yönlü kontrol: En az biri diğeriyle match edebilmeli
			match1 := tt.collider1.MatchMask.CanMatch(tt.collider2.LayerMask)
			match2 := tt.collider2.MatchMask.CanMatch(tt.collider1.LayerMask)
			match := match1 || match2

			if match != tt.shouldMatch {
				t.Errorf("Expected %v, got %v (c1->c2: %v, c2->c1: %v)",
					tt.shouldMatch, match, match1, match2)
			}
		})
	}
}

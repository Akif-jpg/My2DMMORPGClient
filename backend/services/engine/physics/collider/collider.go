package collider

import "github.com/Akif-jpg/MyHobieMMORPGGame/services/engine/physics/shapes"

type Collider struct {
	Shape     shapes.Shape
	LayerMask bitmask // Bu collider hangi layer'da
	MatchMask bitmask // Hangi layer'larla çarpışabilir

	// Yeni alanlar
	IsTrigger bool   // Fiziksel çarpışma mı, sadece tetikleme mi?
	Enabled   bool   // Collider aktif mi?
	Tag       string // Opsiyonel: hızlı tanımlama için

	EntityID string      // veya uint64, her ne kullanıyorsan
	UserData interface{} // Opsiyonel: ek data taşımak için
}

func (c *Collider) CanCollideWith(other *Collider) bool {
	if !c.Enabled || !other.Enabled {
		return false
	}
	return c.MatchMask.CanMatch(other.LayerMask)
}

package collision

import (
	"fmt"

	"github.com/Akif-jpg/MyHobieMMORPGGame/services/engine/physics/collider"
	"github.com/Akif-jpg/MyHobieMMORPGGame/services/engine/physics/shapes"
)

// Global hata tanımları
var (
	ErrNilBody      = fmt.Errorf("collision body is nil")
	ErrNilCollider  = fmt.Errorf("collider is nil")
	ErrUnknownShape = fmt.Errorf("unknown shape type")
)

// CollisionData çarpışma hakkında detaylı bilgi taşır
type CollisionData struct {
	EntityA   string             // Birinci entity'nin ID'si
	EntityB   string             // İkinci entity'nin ID'si
	ColliderA *collider.Collider // Birinci collider referansı
	ColliderB *collider.Collider // İkinci collider referansı
	IsTrigger bool               // Eğer true ise fiziksel tepki verilmez, sadece event tetiklenir

	// Narrow Phase (Detaylı çarpışma) sonucu doldurulacak alanlar:
	Normal       *shapes.Point // Çarpışma normali (A'dan B'ye itme yönü)
	Penetration  float64       // İç içe geçme miktarı (derinlik)
	ContactPoint *shapes.Point // Temas noktası
}

// CollisionBody represents a collidable entity in the world
type CollisionBody struct {
	EntityID  string
	Transform shapes.Point // World position
	Collider  *collider.Collider
}

// NewCollisionBody creates a new collision body with error checking
func NewCollisionBody(entityID string, position shapes.Point, coll *collider.Collider) (*CollisionBody, error) {
	if coll == nil {
		return nil, ErrNilCollider
	}
	// Collider'a entity ID'yi ata (Referans bütünlüğü için)
	coll.EntityID = entityID

	return &CollisionBody{
		EntityID:  entityID,
		Transform: position,
		Collider:  coll,
	}, nil
}

// GetWorldShape returns the shape transformed to world space
func (cb *CollisionBody) GetWorldShape() (shapes.Shape, error) {
	if cb.Collider == nil {
		return nil, ErrNilCollider
	}
	return cb.translateShape(cb.Collider.Shape, cb.Transform)
}

// translateShape helper - shape'i world space'e taşır
func (cb *CollisionBody) translateShape(shape shapes.Shape, offset shapes.Point) (shapes.Shape, error) {
	switch s := shape.(type) {
	case *shapes.Circle:
		worldCircle := *s
		center := s.GetCenter()
		worldCircle.Center = shapes.Point{
			X: center.X + offset.X,
			Y: center.Y + offset.Y,
		}
		return &worldCircle, nil

	case *shapes.Rectangle:
		worldRect := *s
		center := s.GetCenter()
		worldRect.Center = shapes.Point{
			X: center.X + offset.X,
			Y: center.Y + offset.Y,
		}
		return &worldRect, nil

	// Line ve diğerleri eklenebilir...
	default:
		return nil, fmt.Errorf("%w: %T", ErrUnknownShape, shape)
	}
}

// CheckCollision iki body arasındaki çarpışmayı test eder (Narrow Phase)
func CheckCollision(bodyA, bodyB *CollisionBody) (*CollisionData, error) {
	if bodyA == nil || bodyB == nil {
		return nil, ErrNilBody
	}

	// 1. Layer/Mask kontrolü (Çarpışmalı mı?)
	if !bodyA.Collider.CanCollideWith(bodyB.Collider) {
		return nil, nil // Çarpışma yok (kurallar gereği)
	}

	// 2. Şekilleri dünya koordinatlarına taşı
	shapeA, err := bodyA.GetWorldShape()
	if err != nil {
		return nil, err
	}
	shapeB, err := bodyB.GetWorldShape()
	if err != nil {
		return nil, err
	}

	// 3. Geometrik kesişim kontrolü
	// Not: Burada Double Dispatch veya Type Switch kullanmalısın.
	// Örnek olarak Circle-Circle ve Circle-Rect implemente ediyoruz.

	isColliding := false

	switch sA := shapeA.(type) {
	case *shapes.Circle:
		switch sB := shapeB.(type) {
		case *shapes.Circle:
			isColliding = sA.IntersectsCircle(sB)
		case *shapes.Rectangle:
			isColliding = sA.IntersectsRectangle(sB)
		}
	case *shapes.Rectangle:
		switch sB := shapeB.(type) {
		case *shapes.Circle:
			// Rectangle vs Circle (Circle'ın metodunu tersten çağır)
			isColliding = sB.IntersectsRectangle(sA)
		case *shapes.Rectangle:
			// isColliding = sA.IntersectsRectangle(sB) // (Shape paketinde olmalı)
		}
	}

	if isColliding {
		return &CollisionData{
			EntityA:   bodyA.EntityID,
			EntityB:   bodyB.EntityID,
			ColliderA: bodyA.Collider,
			ColliderB: bodyB.Collider,
			IsTrigger: bodyA.Collider.IsTrigger || bodyB.Collider.IsTrigger,
			// Normal, Penetration ve ContactPoint hesaplamaları daha karmaşık matematik gerektirir.
			// Şimdilik sadece bool collision dönüyoruz.
		}, nil
	}

	return nil, nil
}

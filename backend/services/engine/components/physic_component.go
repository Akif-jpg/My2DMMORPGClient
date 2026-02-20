package components

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Akif-jpg/MyHobieMMORPGGame/services/engine/physics/collider"
	"github.com/Akif-jpg/MyHobieMMORPGGame/services/engine/physics/geometry"
	"github.com/google/uuid"
)

type physicType uint16

const (
	staticBody    physicType = 0
	kinematicBody physicType = 1
	rigidBody     physicType = 2
)

// PhysicComponent handles physics simulation for an entity.
// It requires a TransformComponent to update position and rotation based on physics calculations.
type PhysicComponent struct {
	componentID string
	name        string
	isActive    bool

	// Required reference to TransformComponent
	transform *TransformComponent

	// Physics properties
	collider         *collider.Collider
	PhysicType       *physicType
	Velocity         *geometry.Vector2
	RotationVelocity float64

	// Additional physics properties (for future use)
	Mass            float64
	Friction        float64
	Restitution     float64 // Bounciness
	LinearDamping   float64 // Velocity damping over time
	AngularDamping  float64 // Rotation damping over time
}

// NewPhysicComponent creates a new PhysicComponent with the given physics type and collider.
// transform parameter is required and cannot be nil.
func NewPhysicComponent(pt physicType, c *collider.Collider, transform *TransformComponent) (*PhysicComponent, error) {
	if transform == nil {
		return nil, errors.New("TransformComponent is required for PhysicComponent")
	}

	return &PhysicComponent{
		componentID:      "physic-" + uuid.New().String(),
		name:             "Physic",
		isActive:         true,
		transform:        transform,
		collider:         c,
		PhysicType:       &pt,
		Velocity:         geometry.NewVector2(0, 0),
		RotationVelocity: 0,
		Mass:             1.0,
		Friction:         0.5,
		Restitution:      0.0,
		LinearDamping:    0.0,
		AngularDamping:   0.0,
	}, nil
}

// SetTransform sets the TransformComponent reference. This is required for physics to work.
func (c *PhysicComponent) SetTransform(transform *TransformComponent) error {
	if transform == nil {
		return errors.New("TransformComponent cannot be nil")
	}
	c.transform = transform
	return nil
}

// GetTransform returns the associated TransformComponent.
func (c *PhysicComponent) GetTransform() *TransformComponent {
	return c.transform
}

// ---------------------------------------------------------------------------
// Component interface implementation
// ---------------------------------------------------------------------------

func (c *PhysicComponent) ComponentID() string {
	return c.componentID
}

func (c *PhysicComponent) Name() string {
	return c.name
}

func (c *PhysicComponent) IsActive() bool {
	return c.isActive
}

func (c *PhysicComponent) SetActive(active bool) bool {
	prev := c.isActive
	c.isActive = active
	return prev
}

func (c *PhysicComponent) Reset() error {
	c.Velocity = geometry.NewVector2(0, 0)
	c.RotationVelocity = 0
	if c.transform != nil {
		c.transform.Reset()
	}
	return nil
}

func (c *PhysicComponent) Start() error {
	if c.transform == nil {
		return errors.New("TransformComponent is required but not set")
	}
	c.isActive = true
	return nil
}

func (c *PhysicComponent) Update(deltaTime float64) {
	if !c.isActive || c.transform == nil {
		return
	}

	// Apply physics based on physic type
	switch *c.PhysicType {
	case staticBody:
		// Static bodies don't move
		c.Velocity = geometry.NewVector2(0, 0)
		c.RotationVelocity = 0

	case kinematicBody, rigidBody:
		// Apply linear damping
		if c.LinearDamping > 0 {
			dampingFactor := 1.0 - (c.LinearDamping * deltaTime)
			if dampingFactor < 0 {
				dampingFactor = 0
			}
			c.Velocity.X *= dampingFactor
			c.Velocity.Y *= dampingFactor
		}

		// Apply angular damping
		if c.AngularDamping > 0 {
			dampingFactor := 1.0 - (c.AngularDamping * deltaTime)
			if dampingFactor < 0 {
				dampingFactor = 0
			}
			c.RotationVelocity *= dampingFactor
		}

		// Update position based on velocity
		displacementX := c.Velocity.X * deltaTime
		displacementY := c.Velocity.Y * deltaTime
		c.transform.Translate(displacementX, displacementY)

		// Update rotation based on rotation velocity
		if c.RotationVelocity != 0 {
			rotationDelta := c.RotationVelocity * deltaTime
			c.transform.Rotate(rotationDelta)
		}

		// Update collider position if it exists
		if c.collider != nil {
			// Collider position should be synced with transform
			// This is a prototype - actual implementation depends on collider API
		}
	}
}

func (c *PhysicComponent) OnCreate() {
	// Validate that transform is set
	if c.transform == nil {
		// Log warning or handle error
		fmt.Printf("Warning: PhysicComponent created without TransformComponent\n")
	}
}

func (c *PhysicComponent) OnDestroy() {
	c.Velocity = geometry.NewVector2(0, 0)
	c.RotationVelocity = 0
	c.transform = nil
	c.collider = nil
}

func (c *PhysicComponent) Serialize() []byte {
	type serializable struct {
		ComponentID     string             `json:"component_id"`
		Name            string             `json:"name"`
		IsActive        bool               `json:"is_active"`
		PhysicType      uint16             `json:"physic_type"`
		Velocity        *geometry.Vector2  `json:"velocity"`
		RotationVelocity float64           `json:"rotation_velocity"`
		Mass            float64            `json:"mass"`
		Friction        float64            `json:"friction"`
		Restitution     float64            `json:"restitution"`
		LinearDamping   float64            `json:"linear_damping"`
		AngularDamping  float64            `json:"angular_damping"`
		// Note: transform and collider are not serialized as they are references
	}

	s := serializable{
		ComponentID:     c.componentID,
		Name:            c.name,
		IsActive:        c.isActive,
		PhysicType:      uint16(*c.PhysicType),
		Velocity:        c.Velocity,
		RotationVelocity: c.RotationVelocity,
		Mass:            c.Mass,
		Friction:        c.Friction,
		Restitution:     c.Restitution,
		LinearDamping:   c.LinearDamping,
		AngularDamping:  c.AngularDamping,
	}

	data, err := json.Marshal(s)
	if err != nil {
		return nil
	}
	return data
}

func (c *PhysicComponent) Deserialize(data []byte) error {
	type serializable struct {
		ComponentID     string            `json:"component_id"`
		Name            string            `json:"name"`
		IsActive        bool              `json:"is_active"`
		PhysicType      uint16            `json:"physic_type"`
		Velocity        *geometry.Vector2 `json:"velocity"`
		RotationVelocity float64          `json:"rotation_velocity"`
		Mass            float64           `json:"mass"`
		Friction        float64           `json:"friction"`
		Restitution     float64           `json:"restitution"`
		LinearDamping   float64           `json:"linear_damping"`
		AngularDamping  float64           `json:"angular_damping"`
	}

	var s serializable
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	c.componentID = s.ComponentID
	c.name = s.Name
	c.isActive = s.IsActive
	pt := physicType(s.PhysicType)
	c.PhysicType = &pt
	c.Velocity = s.Velocity
	c.RotationVelocity = s.RotationVelocity
	c.Mass = s.Mass
	c.Friction = s.Friction
	c.Restitution = s.Restitution
	c.LinearDamping = s.LinearDamping
	c.AngularDamping = s.AngularDamping

	return nil
}

func (c *PhysicComponent) Clone() Component {
	clone := &PhysicComponent{
		componentID:      "physic-" + uuid.New().String(),
		name:            c.name,
		isActive:        c.isActive,
		transform:       nil, // Clone should not copy transform reference
		collider:        nil, // Clone should not copy collider reference
		PhysicType:      new(physicType),
		Velocity:        geometry.NewVector2(c.Velocity.X, c.Velocity.Y),
		RotationVelocity: c.RotationVelocity,
		Mass:            c.Mass,
		Friction:        c.Friction,
		Restitution:     c.Restitution,
		LinearDamping:   c.LinearDamping,
		AngularDamping:  c.AngularDamping,
	}
	*clone.PhysicType = *c.PhysicType
	return clone
}

// ---------------------------------------------------------------------------
// Physics helper methods
// ---------------------------------------------------------------------------

// AddForce applies a force to the physics body (affects velocity).
// For now, this is a simple prototype - in a full physics engine,
// this would integrate with mass and acceleration.
func (c *PhysicComponent) AddForce(force *geometry.Vector2) {
	if *c.PhysicType == staticBody {
		return // Static bodies don't respond to forces
	}

	// Simple prototype: directly add to velocity
	// In a full implementation, this would be: acceleration = force / mass
	c.Velocity.X += force.X / c.Mass
	c.Velocity.Y += force.Y / c.Mass
}

// AddTorque applies rotational force (affects rotation velocity).
func (c *PhysicComponent) AddTorque(torque float64) {
	if *c.PhysicType == staticBody {
		return // Static bodies don't rotate
	}

	// Simple prototype: directly add to rotation velocity
	c.RotationVelocity += torque / c.Mass
}

// SetVelocity sets the linear velocity directly.
func (c *PhysicComponent) SetVelocity(vx, vy float64) {
	if *c.PhysicType == staticBody {
		return
	}
	c.Velocity.X = vx
	c.Velocity.Y = vy
}

// SetRotationVelocity sets the angular velocity directly.
func (c *PhysicComponent) SetRotationVelocity(angularVel float64) {
	if *c.PhysicType == staticBody {
		return
	}
	c.RotationVelocity = angularVel
}

// GetVelocity returns the current linear velocity.
func (c *PhysicComponent) GetVelocity() *geometry.Vector2 {
	return c.Velocity
}

// GetRotationVelocity returns the current angular velocity.
func (c *PhysicComponent) GetRotationVelocity() float64 {
	return c.RotationVelocity
}

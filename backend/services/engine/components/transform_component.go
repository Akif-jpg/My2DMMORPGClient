package components

import (
	"encoding/json"
	"math"

	"github.com/Akif-jpg/MyHobieMMORPGGame/services/engine/physics/geometry"
	"github.com/google/uuid"
)

// TransformComponent represents the spatial state of an entity in a 2D headless
// physics engine: position, rotation (radians), and uniform scale. It also
// tracks the previous-frame position so that other systems (e.g. interpolation,
// broad-phase) can reason about displacement.
type TransformComponent struct {
	componentID string
	name        string
	isActive    bool

	// Current spatial state
	Position *geometry.Point `json:"position"`
	Rotation float64         `json:"rotation"` // radians, counter-clockwise
	Scale    float64         `json:"scale"`    // uniform scale factor

	// Previous frame state – useful for interpolation & sweep tests
	PreviousPosition *geometry.Point `json:"previous_position"`
	PreviousRotation float64         `json:"previous_rotation"`

	// Hierarchical transform (optional parent)
	parent   *TransformComponent
	children []*TransformComponent
}

// ---------------------------------------------------------------------------
// Constructor
// ---------------------------------------------------------------------------

func NewTransformComponent(position *geometry.Point, rotation float64, scale float64) *TransformComponent {
	return &TransformComponent{
		componentID:      "transform-" + uuid.New().String(),
		name:             "Transform",
		isActive:         true,
		Position:         position,
		Rotation:         rotation,
		Scale:            scale,
		PreviousPosition: geometry.NewPoint(position.X, position.Y),
		PreviousRotation: rotation,
		children:         make([]*TransformComponent, 0),
	}
}

// ---------------------------------------------------------------------------
// Component interface implementation
// ---------------------------------------------------------------------------

func (t *TransformComponent) ComponentID() string {
	return t.componentID
}

func (t *TransformComponent) Name() string {
	return t.name
}

func (t *TransformComponent) IsActive() bool {
	return t.isActive
}

func (t *TransformComponent) SetActive(active bool) bool {
	prev := t.isActive
	t.isActive = active
	return prev
}

func (t *TransformComponent) Reset() error {
	t.Position = geometry.NewPoint(0, 0)
	t.Rotation = 0
	t.Scale = 1
	t.PreviousPosition = geometry.NewPoint(0, 0)
	t.PreviousRotation = 0
	t.parent = nil
	t.children = make([]*TransformComponent, 0)
	return nil
}

func (t *TransformComponent) Start() error {
	t.isActive = true
	return nil
}

func (t *TransformComponent) Update(dt float64) {
	// Snapshot current state as "previous" before physics step applies forces.
	t.PreviousPosition = geometry.NewPoint(t.Position.X, t.Position.Y)
	t.PreviousRotation = t.Rotation
}

func (t *TransformComponent) OnCreate() {
	// No-op: spatial state is set via constructor.
}

func (t *TransformComponent) OnDestroy() {
	// Detach from hierarchy on destruction.
	t.SetParent(nil)
	for _, child := range t.children {
		child.parent = nil
	}
	t.children = nil
}

func (t *TransformComponent) Serialize() []byte {
	data, err := json.Marshal(t)
	if err != nil {
		return nil
	}
	return data
}

func (t *TransformComponent) Deserialize(data []byte) error {
	return json.Unmarshal(data, t)
}

func (t *TransformComponent) Clone() Component {
	clone := &TransformComponent{
		componentID:      "transform-" + uuid.New().String(),
		name:             t.name,
		isActive:         t.isActive,
		Position:         geometry.NewPoint(t.Position.X, t.Position.Y),
		Rotation:         t.Rotation,
		Scale:            t.Scale,
		PreviousPosition: geometry.NewPoint(t.PreviousPosition.X, t.PreviousPosition.Y),
		PreviousRotation: t.PreviousRotation,
		children:         make([]*TransformComponent, 0),
		// parent is intentionally nil – clones start as root nodes
	}
	return clone
}

// ---------------------------------------------------------------------------
// Position helpers
// ---------------------------------------------------------------------------

// SetPosition sets the world-space position directly.
func (t *TransformComponent) SetPosition(x, y float64) {
	t.Position.X = x
	t.Position.Y = y
}

// Translate moves the entity by the given delta in world-space.
func (t *TransformComponent) Translate(dx, dy float64) {
	t.Position.X += dx
	t.Position.Y += dy
}

// TranslateVector moves the entity by the given vector.
func (t *TransformComponent) TranslateVector(v *geometry.Vector2) {
	t.Position.X += v.X
	t.Position.Y += v.Y
}

// GetDisplacement returns the displacement vector from the previous frame.
func (t *TransformComponent) GetDisplacement() *geometry.Vector2 {
	return geometry.NewVector2(
		t.Position.X-t.PreviousPosition.X,
		t.Position.Y-t.PreviousPosition.Y,
	)
}

// DistanceTo returns the Euclidean distance to another transform.
func (t *TransformComponent) DistanceTo(other *TransformComponent) float64 {
	return t.Position.DistanceTo(other.Position)
}

// ---------------------------------------------------------------------------
// Rotation helpers
// ---------------------------------------------------------------------------

// SetRotation sets the rotation in radians.
func (t *TransformComponent) SetRotation(radians float64) {
	t.Rotation = normalizeAngle(radians)
}

// Rotate adds the given angle (radians) to the current rotation.
func (t *TransformComponent) Rotate(radians float64) {
	t.Rotation = normalizeAngle(t.Rotation + radians)
}

// SetRotationDegrees sets the rotation in degrees.
func (t *TransformComponent) SetRotationDegrees(degrees float64) {
	t.Rotation = normalizeAngle(degrees * math.Pi / 180.0)
}

// RotationDegrees returns the current rotation in degrees.
func (t *TransformComponent) RotationDegrees() float64 {
	return t.Rotation * 180.0 / math.Pi
}

// LookAt rotates the transform to face the target point.
func (t *TransformComponent) LookAt(target *geometry.Point) {
	dx := target.X - t.Position.X
	dy := target.Y - t.Position.Y
	t.Rotation = normalizeAngle(math.Atan2(dy, dx))
}

// ---------------------------------------------------------------------------
// Direction vectors
// ---------------------------------------------------------------------------

// Forward returns the unit vector pointing in the entity's facing direction.
func (t *TransformComponent) Forward() *geometry.Vector2 {
	return geometry.NewVector2(math.Cos(t.Rotation), math.Sin(t.Rotation))
}

// Right returns the unit vector pointing to the entity's right (90° CW).
func (t *TransformComponent) Right() *geometry.Vector2 {
	return geometry.NewVector2(math.Sin(t.Rotation), -math.Cos(t.Rotation))
}

// Up returns the unit vector orthogonal-left to the forward direction.
func (t *TransformComponent) Up() *geometry.Vector2 {
	return geometry.NewVector2(-math.Sin(t.Rotation), math.Cos(t.Rotation))
}

// ---------------------------------------------------------------------------
// Scale helpers
// ---------------------------------------------------------------------------

// SetScale sets the uniform scale factor.
func (t *TransformComponent) SetScale(s float64) {
	t.Scale = s
}

// ScaleBy multiplies the current scale factor.
func (t *TransformComponent) ScaleBy(factor float64) {
	t.Scale *= factor
}

// ---------------------------------------------------------------------------
// Coordinate conversion (local ↔ world)
// ---------------------------------------------------------------------------

// LocalToWorld converts a point from this transform's local space to world
// space, applying scale → rotation → translation.
func (t *TransformComponent) LocalToWorld(local *geometry.Point) *geometry.Point {
	// Scale
	sx := local.X * t.Scale
	sy := local.Y * t.Scale
	// Rotate
	cos := math.Cos(t.Rotation)
	sin := math.Sin(t.Rotation)
	rx := sx*cos - sy*sin
	ry := sx*sin + sy*cos
	// Translate
	return geometry.NewPoint(rx+t.Position.X, ry+t.Position.Y)
}

// WorldToLocal converts a world-space point into this transform's local space
// (inverse of LocalToWorld).
func (t *TransformComponent) WorldToLocal(world *geometry.Point) *geometry.Point {
	// Inverse translate
	dx := world.X - t.Position.X
	dy := world.Y - t.Position.Y
	// Inverse rotate
	cos := math.Cos(-t.Rotation)
	sin := math.Sin(-t.Rotation)
	rx := dx*cos - dy*sin
	ry := dx*sin + dy*cos
	// Inverse scale
	if t.Scale != 0 {
		rx /= t.Scale
		ry /= t.Scale
	}
	return geometry.NewPoint(rx, ry)
}

// ---------------------------------------------------------------------------
// Hierarchy helpers (parent-child relationships)
// ---------------------------------------------------------------------------

// SetParent attaches this transform to a parent. Pass nil to detach.
func (t *TransformComponent) SetParent(parent *TransformComponent) {
	// Remove from old parent's children list
	if t.parent != nil {
		t.parent.removeChild(t)
	}
	t.parent = parent
	if parent != nil {
		parent.children = append(parent.children, t)
	}
}

// Parent returns the current parent transform, or nil.
func (t *TransformComponent) Parent() *TransformComponent {
	return t.parent
}

// Children returns a copy of the children slice.
func (t *TransformComponent) Children() []*TransformComponent {
	out := make([]*TransformComponent, len(t.children))
	copy(out, t.children)
	return out
}

// WorldPosition returns the absolute world position taking the parent
// hierarchy into account. If there is no parent this is identical to Position.
func (t *TransformComponent) WorldPosition() *geometry.Point {
	if t.parent == nil {
		return geometry.NewPoint(t.Position.X, t.Position.Y)
	}
	return t.parent.LocalToWorld(t.Position)
}

// WorldRotation returns the accumulated rotation through the hierarchy.
func (t *TransformComponent) WorldRotation() float64 {
	if t.parent == nil {
		return t.Rotation
	}
	return normalizeAngle(t.parent.WorldRotation() + t.Rotation)
}

// WorldScale returns the accumulated scale through the hierarchy.
func (t *TransformComponent) WorldScale() float64 {
	if t.parent == nil {
		return t.Scale
	}
	return t.parent.WorldScale() * t.Scale
}

// ---------------------------------------------------------------------------
// Interpolation (for client-side rendering or network smoothing)
// ---------------------------------------------------------------------------

// Lerp returns an interpolated position between the previous and current
// frame positions. alpha should be in [0, 1].
func (t *TransformComponent) Lerp(alpha float64) *geometry.Point {
	return geometry.NewPoint(
		t.PreviousPosition.X+(t.Position.X-t.PreviousPosition.X)*alpha,
		t.PreviousPosition.Y+(t.Position.Y-t.PreviousPosition.Y)*alpha,
	)
}

// LerpRotation returns an interpolated rotation (shortest path).
func (t *TransformComponent) LerpRotation(alpha float64) float64 {
	diff := t.Rotation - t.PreviousRotation
	// Shortest path around the circle
	if diff > math.Pi {
		diff -= 2 * math.Pi
	} else if diff < -math.Pi {
		diff += 2 * math.Pi
	}
	return normalizeAngle(t.PreviousRotation + diff*alpha)
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

func (t *TransformComponent) removeChild(child *TransformComponent) {
	for i, c := range t.children {
		if c == child {
			t.children = append(t.children[:i], t.children[i+1:]...)
			return
		}
	}
}

// normalizeAngle wraps an angle to the range [0, 2π).
func normalizeAngle(a float64) float64 {
	a = math.Mod(a, 2*math.Pi)
	if a < 0 {
		a += 2 * math.Pi
	}
	return a
}

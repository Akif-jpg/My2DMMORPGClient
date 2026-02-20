package components

type Component interface {
	// Define the methods that all components must implement
	ComponentID() string
	Name() string
	IsActive() bool
	SetActive(active bool) bool
	Serialize() []byte
	Deserialize(data []byte) error
	Clone() Component
	Reset() error
	Start() error
	Update(deltaTime float64)
	OnCreate()
	OnDestroy()
}

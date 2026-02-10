package entities

type Entity struct {
	ID        int
	Name      string
	Width     int
	Height    int
	Collision []int
}

type IEntity interface {
	GetID() int
	GetName() string
	GetWidth() int
	GetHeight() int
	GetCollision() []int
}

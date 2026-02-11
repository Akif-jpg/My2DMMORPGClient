package collision

import (
	"fmt"
)

type CollisionSystemError struct {
	Message string
}

func (e *CollisionSystemError) Error() string {
	return fmt.Sprintf("Collision System Error: %s", e.Message)
}

func NewCollisionSystemError(message string) *CollisionSystemError {
	return &CollisionSystemError{Message: message}
}

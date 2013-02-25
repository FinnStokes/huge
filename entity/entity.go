// Package entity implements representation of in-game objects as entities.
package entity

// Entity is a type to represent an object in the game world. It uses components to describe
// its form and function.
type Entity struct {
	id         int
	Components map[string]interface{}
}

// Id returns the integral id of the entity
func (e *Entity) Id() int {
	return e.id
}

// Position is a type used as a component for storing a position in 2d coordinates
type Position struct {
	X, Y float32
}

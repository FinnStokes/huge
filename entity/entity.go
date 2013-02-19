package entity

type Entity struct {
	id         int
	Components map[string]interface{}
}

func (e *Entity) Id() int {
	return e.id
}

type Position struct {
	X, Y int
}

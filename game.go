package huge

import (
	"./entity"
)

type Game struct {
	Entities *entity.Manager
}

func NewGame() *Game {
	g := new(Game)
	g.Entities = entity.NewManager()
	return g
}

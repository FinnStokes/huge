package huge

import (
	"testing"

	"./entity"
	"./sprite"
	"./system"
)

func TestInit(t *testing.T) {
	g := NewGame()
	g.Quit()
	g.Run()
}

func TestRender(t *testing.T) {
	g := NewGame()
	g.Systems.AddSystem(system.Normal, sprite.NewManager())
	e := g.Entities.New()
	e.Components["pos"] = &entity.Position{100, 100}
	g.Run()
}

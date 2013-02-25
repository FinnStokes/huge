package huge

import (
	"testing"

	"github.com/FinnStokes/huge/entity"
	"github.com/FinnStokes/huge/sprite"
	"github.com/FinnStokes/huge/system"
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
	var err error
	e.Components["sprite"], err = g.Resources.GetSprite("sprite")
	if err != nil {
		t.Fatal(err)
	}
	g.Run()
}

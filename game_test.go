package huge

import (
	"image"
	_ "image/png"
	"os"
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
	file, err := os.Open("sprite.png")
	if err != nil {
		t.Fatalf("error opening image 'sprite.png': %v", err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		t.Fatalf("error decoding image 'sprite.png': %v", err)
	}
	anim := &sprite.Animation{
		[]int{0, 1, 2, 3},
		8,
		nil,
	}
	anim.Next = anim
	e.Components["sprite"] = &sprite.Sprite{
		img,
		make(map[string]*sprite.Animation, 1),
		64, 64,
		anim,
		0,
		0,
	}
	g.Run()
}

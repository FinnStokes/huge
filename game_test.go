package huge

import (
	"testing"

	"./entity"
	"./system"
)

func TestInit(t *testing.T) {
	g := NewGame()
	g.Quit()
	g.Run()
}

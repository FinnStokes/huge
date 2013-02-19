package sprite

import (
	"time"

	"../entity"

	"github.com/go-gl/gl"
)

// Manager is a placeholder system that draws sprites for all entities with a position component
type Manager struct {
}

// NewManager returns an initialised sprite manager
func NewManager() *Manager {
	return new(Manager)
}

// Update moves animated sprites on to the next frame when appropriate and performs any required
// operations once the animation is complete, such as advancing to the follow-up animation.
func (m *Manager) Update(dt time.Duration, entities *entity.Manager) {
}

// Draw draws the current frame of all entities with sprite components at the position given by the
// pos component.
func (m *Manager) Draw(entities *entity.Manager) {
	gl.Enable(gl.BLEND)
	gl.Enable(gl.POINT_SMOOTH)
	gl.Enable(gl.LINE_SMOOTH)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Color4f(0.0, 0.0, 0.0, 0.1)
	gl.Begin(gl.QUADS)
	for _, e := range entities.All() {
		if pos, ok := e.Components["pos"].(*entity.Position); ok {
			gl.Vertex2i(pos.X, pos.Y)
			gl.Vertex2i(pos.X+10, pos.Y)
			gl.Vertex2i(pos.X+10, pos.Y+10)
			gl.Vertex2i(pos.X, pos.Y+10)
		}
	}
}

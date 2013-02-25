package sprite

import (
	"image"
	"time"

	"github.com/FinnStokes/huge/camera"
	"github.com/FinnStokes/huge/entity"

	"github.com/go-gl/gl"
)

// Manager is a placeholder system that draws sprites for all entities with a position component
type Manager struct {
	textures map[image.Image]gl.Texture
}

// NewManager returns an initialised sprite manager
func NewManager() *Manager {
	m := new(Manager)
	m.textures = make(map[image.Image]gl.Texture)
	return m
}

// Update moves animated sprites on to the next frame when appropriate and performs any required
// operations once the animation is complete, such as advancing to the follow-up animation.
func (m *Manager) Update(dt time.Duration, entities *entity.Manager) {
	for _, e := range entities.All() {
		if sprite, ok := e.Components["sprite"].(*Sprite); ok {
			sprite.FrameTime += dt
			for sprite.FrameTime*time.Duration(sprite.CurrentAnimation.Fps) >= time.Second {
				sprite.FrameTime -= time.Second / time.Duration(sprite.CurrentAnimation.Fps)
				sprite.CurrentFrame++
			}
			for sprite.CurrentFrame >= len(sprite.CurrentAnimation.Frames) {
				sprite.CurrentFrame -= len(sprite.CurrentAnimation.Frames)
				sprite.CurrentAnimation = sprite.CurrentAnimation.Next
			}
		}
	}
}

// Draw draws the current frame of all entities with sprite components at the position given by the
// pos component.
func (m *Manager) Draw(c *camera.Camera, entities *entity.Manager) {
	gl.Enable(gl.BLEND)
	gl.Disable(gl.LIGHTING)
	gl.TexEnvf(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.REPLACE)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Color4f(1.0, 1.0, 1.0, 1.0)
	gl.Enable(gl.TEXTURE_2D)
	for _, e := range entities.All() {
		if pos, ok := e.Components["pos"].(*entity.Position); ok {
			if sprite, ok := e.Components["sprite"].(*Sprite); ok {
				if c.World.Intersects(&camera.Rectangle{
					pos.X, pos.Y,
					float32(sprite.Width) * c.World.Width / float32(c.Screen.Width),
					float32(sprite.Height) * c.World.Height / float32(c.Screen.Height),
				}) {
					tex, ok := m.textures[sprite.Image]
					if !ok {
						tex = loadImageAsTexture(sprite.Image)
						m.textures[sprite.Image] = tex
					}

					tw := float32(sprite.Width) / float32(sprite.Image.Bounds().Dx())
					th := float32(sprite.Height) / float32(sprite.Image.Bounds().Dy())
					n := int(1.0 / tw)
					f := sprite.CurrentAnimation.Frames[sprite.CurrentFrame]
					tx := tw * float32(f%n)
					ty := th * float32(f/n)

					x := int((pos.X - c.World.X) * float32(c.Screen.Width) / c.World.Width)
					y := int((pos.Y - c.World.Y) * float32(c.Screen.Height) / c.World.Height)
					w := sprite.Width
					h := sprite.Height

					tex.Bind(gl.TEXTURE_2D)
					gl.Begin(gl.QUADS)

					gl.TexCoord2f(tx, ty)
					gl.Vertex2i(x, y)

					gl.TexCoord2f(tx+tw, ty)
					gl.Vertex2i(x+w, y)

					gl.TexCoord2f(tx+tw, ty+th)
					gl.Vertex2i(x+w, y+h)

					gl.TexCoord2f(tx, ty+th)
					gl.Vertex2i(x, y+h)

					gl.End()
					tex.Unbind(gl.TEXTURE_2D)
				}
			}
		}
	}
}

type rgba struct {
	r, g, b, a uint16
}

func loadImageAsTexture(img image.Image) gl.Texture {
	tex := gl.GenTexture()
	tex.Bind(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	data := make([]rgba, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			data[x+y*w].r = uint16(r)
			data[x+y*w].g = uint16(g)
			data[x+y*w].b = uint16(b)
			data[x+y*w].a = uint16(a)
		}
	}
	gl.TexImage2D(
		gl.TEXTURE_2D,     // target
		0,                 // level, 0 = base, no mipmap,
		gl.RGBA,           // internal format
		w,                 // width
		h,                 // height
		0,                 // border
		gl.RGBA,           // format
		gl.UNSIGNED_SHORT, // type
		data,              // image
	)
	tex.Unbind(gl.TEXTURE_2D)
	return tex
}

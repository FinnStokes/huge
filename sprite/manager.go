package sprite

import (
	"image"
	"log"
	"time"

	"../entity"

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
func (m *Manager) Draw(entities *entity.Manager) {
	gl.Enable(gl.BLEND)
	gl.Disable(gl.LIGHTING)
	gl.TexEnvf(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.REPLACE)
	//gl.Enable(gl.POINT_SMOOTH)
	//gl.Enable(gl.LINE_SMOOTH)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Color4f(1.0, 1.0, 1.0, 1.0)
	gl.Enable(gl.TEXTURE_2D)
	//gl.ShadeModel(gl.SMOOTH)
	for _, e := range entities.All() {
		if pos, ok := e.Components["pos"].(*entity.Position); ok {
			if sprite, ok := e.Components["sprite"].(*Sprite); ok {
				tex, ok := m.textures[sprite.Image]
				if !ok {
					tex = loadImageAsTexture(sprite.Image)
					m.textures[sprite.Image] = tex
				}

				w := float32(sprite.Width) / float32(sprite.Image.Bounds().Dx())
				h := float32(sprite.Height) / float32(sprite.Image.Bounds().Dy())
				n := int(1.0 / w)
				f := sprite.CurrentAnimation.Frames[sprite.CurrentFrame]
				x := w * float32(f%n)
				y := h * float32(f/n)

				tex.Bind(gl.TEXTURE_2D)
				gl.Begin(gl.QUADS)

				gl.TexCoord2f(x, y)
				gl.Vertex2i(pos.X, pos.Y)

				gl.TexCoord2f(x+w, y)
				gl.Vertex2i(pos.X+sprite.Width, pos.Y)

				gl.TexCoord2f(x+w, y+h)
				gl.Vertex2i(pos.X+sprite.Width, pos.Y+sprite.Height)

				gl.TexCoord2f(x, y+h)
				gl.Vertex2i(pos.X, pos.Y+sprite.Height)

				gl.End()
				tex.Unbind(gl.TEXTURE_2D)
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
	log.Println(img.At(0, 0).RGBA())
	log.Println(img.At(32, 16).RGBA())
	log.Println(img.At(32, 32).RGBA())
	log.Println(img.At(63, 63).RGBA())
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

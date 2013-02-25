package huge

import (
	"log"
	"time"

	"github.com/FinnStokes/huge/camera"
	"github.com/FinnStokes/huge/entity"
	"github.com/FinnStokes/huge/resource"
	"github.com/FinnStokes/huge/system"

	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
)

type Game struct {
	Camera    *camera.Camera
	Entities  *entity.Manager
	Resources *resource.Manager
	Systems   *system.Manager
	ticker    [system.NumSpeeds]*time.Ticker
	oldTime   [system.NumSpeeds]time.Time
	running   bool
	quitting  bool
}

func NewGame() *Game {
	g := new(Game)
	g.Camera = new(camera.Camera)
	g.Entities = entity.NewManager()
	g.Resources = resource.NewManager()
	g.Systems = system.NewManager()
	g.ticker[system.Slow] = time.NewTicker(time.Second)
	g.ticker[system.Normal] = time.NewTicker(20 * time.Millisecond)
	g.ticker[system.Fast] = time.NewTicker(10 * time.Millisecond)
	return g
}

func (g *Game) Run() {
	defer g.terminate()

	var err error
	if err = glfw.Init(); err != nil {
		log.Fatalf("%v\n", err)
		return
	}

	defer glfw.Terminate()

	if err = glfw.OpenWindow(640, 480, 16, 16, 16, 16, 0, 0, glfw.Windowed); err != nil {
		log.Fatalf("%v\n", err)
		return
	}

	defer glfw.CloseWindow()

	g.Camera.Screen.Width, g.Camera.Screen.Height = 640, 480
	g.Camera.World.X, g.Camera.World.Y = 0, 0
	g.Camera.World.Width, g.Camera.World.Height = 640, 480

	glfw.SetWindowTitle("Draw")
	glfw.SetSwapInterval(1)
	glfw.SetWindowSizeCallback(func(w, h int) { g.onResize(w, h) })

	g.running = true
	for !g.quitting && glfw.WindowParam(glfw.Opened) == 1 {
		select {
		case t := <-g.ticker[system.Slow].C:
			if !g.oldTime[system.Slow].IsZero() {
				g.Systems.Update(system.Slow, t.Sub(g.oldTime[system.Slow]), g.Entities)
			}
			g.oldTime[system.Slow] = t
		case t := <-g.ticker[system.Normal].C:
			if !g.oldTime[system.Normal].IsZero() {
				g.Systems.Update(system.Normal, t.Sub(g.oldTime[system.Normal]), g.Entities)
			}
			g.oldTime[system.Normal] = t
			gl.ClearColor(1, 1, 1, 1)
			gl.Clear(gl.COLOR_BUFFER_BIT)
			g.Systems.Draw(g.Camera, g.Entities)
			glfw.SwapBuffers()
		case t := <-g.ticker[system.Fast].C:
			if !g.oldTime[system.Fast].IsZero() {
				g.Systems.Update(system.Fast, t.Sub(g.oldTime[system.Fast]), g.Entities)
			}
			g.oldTime[system.Fast] = t
		}
	}
}

func (g *Game) onResize(w, h int) {
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Viewport(0, 0, w, h)
	gl.Ortho(0, float64(w), float64(h), 0, -1.0, 1.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	g.Camera.Screen.Width, g.Camera.Screen.Height = w, h
	g.Camera.World.Width, g.Camera.World.Height = float32(w), float32(h)
}

func (g *Game) terminate() {
	g.running = false
	g.quitting = false
}

func (g *Game) Quit() {
	g.quitting = true
}

func (g *Game) SetSpeed(speed system.Speed, duration time.Duration) {
	g.ticker[speed].Stop()
	g.ticker[speed] = time.NewTicker(duration)
}

package huge

import (
	"time"

	"./entity"
	"./system"
)

type Game struct {
	Entities *entity.Manager
	Systems  *system.Manager
	ticker   [system.NumSpeeds]*time.Ticker
	oldTime  [system.NumSpeeds]time.Time
	running  bool
	quitting bool
}

func NewGame() *Game {
	g := new(Game)
	g.Entities = entity.NewManager()
	g.Systems = system.NewManager()
	g.ticker[system.Slow] = time.NewTicker(time.Second)
	g.ticker[system.Normal] = time.NewTicker(20 * time.Millisecond)
	g.ticker[system.Fast] = time.NewTicker(time.Millisecond)
	return g
}

func (g *Game) Run() {
	defer g.terminate()

	g.running = true
	for !g.quitting {
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
			g.Systems.Draw(g.Entities)
		case t := <-g.ticker[system.Fast].C:
			if !g.oldTime[system.Fast].IsZero() {
				g.Systems.Update(system.Fast, t.Sub(g.oldTime[system.Fast]), g.Entities)
			}
			g.oldTime[system.Fast] = t
		}
	}
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

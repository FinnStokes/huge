package gauge

type Game struct {
	Entities EntityManager
}

func NewGame() *Game {
	em := NewEntityManager()
	return &Game{em}
}

package battle

type computerPlayer struct {
	world *worldState
}

func newComputerPlayer(world *worldState) *computerPlayer {
	return &computerPlayer{world: world}
}

func (p *computerPlayer) Update(delta float64) {}

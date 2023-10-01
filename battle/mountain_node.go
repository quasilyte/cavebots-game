package battle

import (
	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type lootKind int

const (
	lootUnknown lootKind = iota
	lootNone
	lootExtraStones
	lootEasyDig
	lootIronDeposit
	lootLargeIronDeposit
	lootLavaCell
	lootFlatCell
	lootBotHarvester
	lootBotVanguard
	lootBotPatrol
	lootBotTitan
	lootBotGenerator
)

type mountainNode struct {
	world  *worldState
	pos    gmath.Vec
	sprite *ge.Sprite

	outer bool

	loot lootKind
}

func newMountainNode(world *worldState, pos gmath.Vec) *mountainNode {
	return &mountainNode{
		pos:   pos,
		world: world,
	}
}

func (m *mountainNode) Init(scene *ge.Scene) {
	s := scene.NewSprite(assets.ImageMountains)
	s.FlipHorizontal = scene.Rand().Bool()
	s.Pos.Base = &m.pos
	s.FrameOffset.X = float64(scene.Rand().IntRange(0, 1)) * s.FrameWidth
	m.world.stage.AddSprite(s)
	m.sprite = s
	if m.outer {
		m.sprite.SetColorScale(ge.ColorScale{R: 1.35, G: 1.2, B: 0.8, A: 1})
	} else {
		m.sprite.SetColorScale(ge.ColorScale{R: 0.5, G: 0.5, B: 0.5, A: 1})
	}
}

func (m *mountainNode) IsDisposed() bool {
	return m.sprite.IsDisposed()
}

func (m *mountainNode) Dispose() {
	m.sprite.Dispose()
}

func (m *mountainNode) Update(delta float64) {}

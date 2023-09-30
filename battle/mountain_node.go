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
	lootIronDeposit
	lootLargeIronDeposit
	lootLavaCell
	lootFlatCell
	lootBotHarvester
	lootBotGuard
	lootBotPatrol
)

type mountainNode struct {
	pos    gmath.Vec
	sprite *ge.Sprite

	outer bool

	loot lootKind
}

func newMountainNode(pos gmath.Vec) *mountainNode {
	return &mountainNode{
		pos: pos,
	}
}

func (m *mountainNode) Init(scene *ge.Scene) {
	s := scene.NewSprite(assets.ImageMountains)
	s.FlipHorizontal = scene.Rand().Bool()
	s.Pos.Base = &m.pos
	s.FrameOffset.X = float64(scene.Rand().IntRange(0, 1)) * s.FrameWidth
	scene.AddGraphics(s)
	m.sprite = s
}

func (m *mountainNode) IsDisposed() bool {
	return m.sprite.IsDisposed()
}

func (m *mountainNode) Dispose() {
	m.sprite.Dispose()
}

func (m *mountainNode) Update(delta float64) {}

package battle

import (
	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type hardTerrainNode struct {
	sprite *ge.Sprite

	pos gmath.Vec
}

func newHardTerrainNode(pos gmath.Vec) *hardTerrainNode {
	return &hardTerrainNode{
		pos: pos,
	}
}

func (t *hardTerrainNode) Init(scene *ge.Scene) {
	t.sprite = scene.NewSprite(assets.ImageHardTerrain)
	t.sprite.Pos.Base = &t.pos
	t.sprite.FlipHorizontal = scene.Rand().Bool()
	scene.AddGraphicsBelow(t.sprite, 1)
}

func (t *hardTerrainNode) IsDisposed() bool { return false }

func (t *hardTerrainNode) Update(delta float64) {}

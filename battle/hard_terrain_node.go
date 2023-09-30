package battle

import (
	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type hardTerrainNode struct {
	world *worldState

	sprite *ge.Sprite

	buildOptions [2]*unitStats

	building *unitNode

	pos gmath.Vec
}

func newHardTerrainNode(world *worldState, pos gmath.Vec, buildOptions [2]*unitStats) *hardTerrainNode {
	return &hardTerrainNode{
		world:        world,
		pos:          pos,
		buildOptions: buildOptions,
	}
}

func (t *hardTerrainNode) Init(scene *ge.Scene) {
	t.sprite = scene.NewSprite(assets.ImageHardTerrain)
	t.sprite.Pos.Base = &t.pos
	t.sprite.FlipHorizontal = scene.Rand().Bool()
	t.world.stage.AddSpriteBelow(t.sprite)
}

func (t *hardTerrainNode) IsDisposed() bool { return false }

func (t *hardTerrainNode) Update(delta float64) {
}

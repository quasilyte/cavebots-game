package battle

import (
	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type coreUnit struct {
	world  *worldState
	sprite *ge.Sprite

	pos gmath.Vec
}

func newCoreUnit(world *worldState) *coreUnit {
	return &coreUnit{world: world}
}

func (u *coreUnit) Init(scene *ge.Scene) {
	u.sprite = scene.NewSprite(assets.ImageDroneCore)
	u.sprite.Pos.Base = &u.pos
	scene.AddGraphics(u.sprite)
}

func (u *coreUnit) IsDisposed() bool {
	return u.sprite.IsDisposed()
}

func (u *coreUnit) Update(delta float64) {

}

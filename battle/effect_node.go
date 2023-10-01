package battle

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/gmath"
)

type effectNode struct {
	world *worldState

	pos     gmath.Vec
	image   resource.ImageID
	anim    *ge.Animation
	above   bool
	rotates bool
	noFlip  bool

	rotation gmath.Rad

	EventCompleted gesignal.Event[gesignal.Void]
}

func newEffectNode(world *worldState, pos gmath.Vec, above bool, image resource.ImageID) *effectNode {
	return &effectNode{
		world: world,
		pos:   pos,
		image: image,
		above: above,
	}
}

func (e *effectNode) Init(scene *ge.Scene) {
	var sprite *ge.Sprite
	if e.anim == nil {
		sprite = scene.NewSprite(e.image)
		sprite.Pos.Base = &e.pos
	} else {
		sprite = e.anim.Sprite()
	}
	sprite.Rotation = &e.rotation
	if !e.noFlip {
		sprite.FlipHorizontal = e.world.rand.Bool()
	}
	if e.above {
		e.world.stage.AddSpriteAbove(sprite)
	} else {
		e.world.stage.AddSprite(sprite)
	}
	if e.anim == nil {
		e.anim = ge.NewAnimation(sprite, -1)
		e.anim.SetSecondsPerFrame(0.05)
	}
}

func (e *effectNode) IsDisposed() bool {
	return e.anim.IsDisposed()
}

func (e *effectNode) Dispose() {
	e.anim.Sprite().Dispose()
}

func (e *effectNode) Update(delta float64) {
	if e.anim.Tick(delta) {
		e.EventCompleted.Emit(gesignal.Void{})
		e.Dispose()
		return
	}
	if e.rotates {
		e.rotation += gmath.Rad(delta * 2)
	}
}

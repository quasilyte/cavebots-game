package battle

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type projectileNode struct {
	world    *worldState
	attacker *unitNode
	target   *unitNode

	pos      gmath.Vec
	toPos    gmath.Vec
	rotation gmath.Rad

	sprite *ge.Sprite
}

type projectileNodeConfig struct {
	attacker  *unitNode
	target    *unitNode
	targetPos gmath.Vec
}

func newProjectileNode(config projectileNodeConfig) *projectileNode {
	return &projectileNode{
		world:    config.attacker.world,
		attacker: config.attacker,
		target:   config.target,
		toPos:    config.targetPos,
	}
}

func (p *projectileNode) Init(scene *ge.Scene) {
	p.pos = p.attacker.pos

	p.rotation = p.pos.AngleToPoint(p.toPos)

	p.sprite = scene.NewSprite(p.attacker.stats.weapon.projectileImage)
	p.sprite.Pos.Base = &p.pos
	p.sprite.Rotation = &p.rotation
	p.world.stage.AddSpriteAbove(p.sprite)
	p.sprite.Visible = false
}

func (p *projectileNode) IsDisposed() bool {
	return p.sprite.IsDisposed()
}

func (p *projectileNode) Update(delta float64) {
	weapon := p.attacker.stats.weapon
	travelled := weapon.projectileSpeed * delta

	if p.pos.DistanceTo(p.toPos) <= travelled {
		p.detonate()
		return
	}
	p.sprite.Visible = true
	p.pos = p.pos.MoveTowards(p.toPos, travelled)
}

func (p *projectileNode) detonate() {
	if p.IsDisposed() {
		return
	}

	p.sprite.Dispose()

	weapon := p.attacker.stats.weapon

	if p.toPos.DistanceSquaredTo(p.target.pos) > (20 * 20) {
		return
	}

	p.target.OnDamage(weapon.damage)
}

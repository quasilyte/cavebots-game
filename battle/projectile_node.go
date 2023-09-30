package battle

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type projectileNode struct {
	world     *worldState
	attacker  *unitNode
	target    *unitNode
	fireDelay float64

	pos      gmath.Vec
	toPos    gmath.Vec
	rotation gmath.Rad

	arcProgressionScaling float64
	arcProgression        float64
	arcStart              gmath.Vec
	arcFrom               gmath.Vec
	arcTo                 gmath.Vec

	sprite *ge.Sprite
}

type projectileNodeConfig struct {
	attacker  *unitNode
	target    *unitNode
	targetPos gmath.Vec
	fireDelay float64
}

func newProjectileNode(config projectileNodeConfig) *projectileNode {
	return &projectileNode{
		world:     config.attacker.world,
		attacker:  config.attacker,
		target:    config.target,
		toPos:     config.targetPos,
		fireDelay: config.fireDelay,
	}
}

func (p *projectileNode) Init(scene *ge.Scene) {
	p.pos = p.attacker.pos

	weapon := p.attacker.stats.weapon

	if weapon.arcPower == 0 {
		p.rotation = p.pos.AngleToPoint(p.toPos)
	} else {
		arcPower := weapon.arcPower
		speed := weapon.projectileSpeed
		p.rotation = -math.Pi / 2
		if p.toPos.Y >= p.pos.Y {
			arcPower *= 0.6
			speed *= 1.2
		}
		dist := p.pos.DistanceTo(p.toPos)
		t := dist / speed
		p.arcProgressionScaling = 1.0 / t
		power := gmath.Vec{Y: dist * arcPower}
		p.arcFrom = p.pos.Add(power)
		p.arcTo = p.toPos.Add(power)
		p.arcStart = p.pos
	}

	p.sprite = scene.NewSprite(weapon.projectileImage)
	p.sprite.Pos.Base = &p.pos
	p.sprite.Rotation = &p.rotation
	p.sprite.Visible = false
	p.world.stage.AddSpriteAbove(p.sprite)
}

func (p *projectileNode) IsDisposed() bool {
	return p.sprite.IsDisposed()
}

func (p *projectileNode) Update(delta float64) {
	weapon := p.attacker.stats.weapon

	if p.fireDelay > 0 {
		if p.attacker.IsDisposed() {
			p.dispose()
			return
		}
		p.fireDelay -= delta
		if p.fireDelay <= 0 {
			p.sprite.Visible = true
			p.arcStart = p.pos
		} else {
			return
		}
	}

	if p.arcProgressionScaling == 0 {
		travelled := weapon.projectileSpeed * delta
		if p.pos.DistanceTo(p.toPos) <= travelled {
			p.detonate()
			return
		}
		p.sprite.Visible = true
		p.pos = p.pos.MoveTowards(p.toPos, travelled)
		return
	}

	p.arcProgression += delta * p.arcProgressionScaling
	if p.arcProgression >= 1 {
		p.detonate()
		return
	}
	newPos := p.arcStart.CubicInterpolate(p.arcFrom, p.toPos, p.arcTo, p.arcProgression)
	p.rotation = p.pos.AngleToPoint(newPos)
	p.pos = newPos
	p.sprite.Visible = true
}

func (p *projectileNode) dispose() {
	p.sprite.Dispose()
}

func (p *projectileNode) detonate() {
	if p.IsDisposed() {
		return
	}

	p.dispose()

	weapon := p.attacker.stats.weapon

	if p.toPos.DistanceSquaredTo(p.target.pos) > (20 * 20) {
		return
	}

	p.target.OnDamage(weapon.damage, p.attacker)
}

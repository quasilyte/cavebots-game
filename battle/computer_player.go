package battle

import (
	"github.com/quasilyte/gmath"
)

type computerPlayer struct {
	world *worldState

	attackGroup []*unitNode

	attackDelay float64
}

func newComputerPlayer(world *worldState) *computerPlayer {
	return &computerPlayer{
		world:       world,
		attackGroup: make([]*unitNode, 0, 16),
		attackDelay: world.rand.FloatRange(15, 30),
	}
}

func (p *computerPlayer) Update(delta float64) {
	if p.world.creepBase == nil || p.world.creepBase.IsDisposed() {
		return
	}

	p.attackDelay = gmath.ClampMin(p.attackDelay-delta, 0)
	if p.attackDelay == 0 {
		if delay := p.maybeDoAttack(); delay != 0 {
			p.attackDelay = delay
		} else {
			p.attackDelay = p.world.rand.FloatRange(2.5, 8.5)
		}
	}
}

func (p *computerPlayer) maybeDoAttack() float64 {
	attackBudget := p.world.creepBaseAttackBudget

	p.attackGroup = p.attackGroup[:0]

	rng := p.world.rand

	randIterate(p.world.rand, p.world.creeps, func(creep *unitNode) bool {
		if creep.stats.building {
			return false
		}
		if creep.order != orderNone {
			return false
		}
		if creep.stats.score > attackBudget {
			return false
		}
		attackBudget -= creep.stats.score
		p.attackGroup = append(p.attackGroup, creep)
		return attackBudget <= 0
	})
	if len(p.attackGroup) == 0 {
		return rng.FloatRange(8, 30)
	}

	// Every 2.2 minutes give +1 budget growth.
	extraBudgetGrowth := int(p.world.creepBaseLevel / (2.2 * 60.0))

	var bestTarget *unitNode
	var bestScore float64
	for _, u := range p.world.playerUnits {
		unitScore := u.stats.botPriority * rng.FloatRange(0.8, 1.2)
		if u.health < u.stats.maxHealth*0.5 {
			unitScore *= 1.2
		}
		if u.pos.DistanceSquaredTo(p.world.creepBase.pos) < (200 * 200) {
			unitScore *= 1.5
		}
		if unitScore > bestScore {
			bestScore = unitScore
			bestTarget = u
		}
	}
	if bestScore < 10 && rng.Chance(0.4) {
		p.world.creepBaseAttackBudget += 2 + extraBudgetGrowth
		return rng.FloatRange(15, 30)
	}
	if bestTarget == nil {
		return 40
	}

	for _, u := range p.attackGroup {
		dst := bestTarget.pos
		if rng.Chance(0.4) {
			dst = dst.Add(rng.Offset(-64, 64))
		}
		u.SendAttacking(dst)
	}

	p.world.creepBaseAttackBudget += 1 + extraBudgetGrowth

	return rng.FloatRange(20, 60)
}

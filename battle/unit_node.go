package battle

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/pathing"
)

type unitOrder int

const (
	orderNone unitOrder = iota
	orderDig
	orderHarvestResource
	orderDeliverResource
	orderPatrolMove
)

type unitNode struct {
	world  *worldState
	sprite *ge.Sprite
	anim   *ge.Animation

	stats *unitStats

	health float64

	scene *ge.Scene

	path        pathing.GridPath
	pathDest    gmath.Vec
	waypoint    gmath.Vec
	order       unitOrder
	orderTarget any

	specialDelay float64

	pos gmath.Vec

	cargo int

	offline bool

	EventDisposed gsignal.Event[*unitNode]
}

func newUnitNode(world *worldState, stats *unitStats) *unitNode {
	return &unitNode{
		world: world,
		stats: stats,
	}
}

func (u *unitNode) Init(scene *ge.Scene) {
	u.scene = scene

	u.health = u.stats.maxHealth

	u.sprite = scene.NewSprite(u.stats.img)
	u.sprite.Pos.Base = &u.pos
	scene.AddGraphics(u.sprite)

	if u.sprite.FrameWidth != u.sprite.ImageWidth() {
		u.anim = ge.NewRepeatedAnimation(u.sprite, -1)
		u.anim.Tick(scene.Rand().FloatRange(0.1, 4.6))
		u.anim.SetAnimationSpan(1)
	}
}

func (u *unitNode) IsDisposed() bool {
	return u.sprite.IsDisposed()
}

func (u *unitNode) Dispose() {
	if u.IsDisposed() {
		return
	}
	u.EventDisposed.Emit(u)

	u.sprite.Dispose()
}

func (u *unitNode) SendTo(pos gmath.Vec) {
	u.sendTo(pos)
	u.order = orderNone
}

func (u *unitNode) SendDigging(pos gmath.Vec) {
	u.sendTo(pos)
	u.order = orderDig
}

func (u *unitNode) sendTo(pos gmath.Vec) {
	from := u.world.grid.PosToCoord(u.pos.X, u.pos.Y)
	to := u.world.grid.PosToCoord(pos.X, pos.Y)
	result := u.world.astar.BuildPath(u.world.grid, from, to, normalLayer)
	u.path = result.Steps

	u.pathDest = makePos(u.world.grid.CoordToPos(result.Finish))
	u.waypoint = makePos(u.world.grid.AlignPos(u.pos.X, u.pos.Y))
}

func (u *unitNode) Update(delta float64) {
	if u.anim != nil {
		u.anim.Tick(delta)
	}

	if !u.waypoint.IsZero() {
		newPos, reached := moveTowardsWithSpeed(u.pos, u.waypoint, delta, u.movementSpeed())
		u.pos = newPos
		if reached {
			if u.path.HasNext() {
				nextPos := nextPathWaypoint(u.world, u.pos, &u.path, normalLayer)
				u.waypoint = nextPos.Add(u.world.rand.Offset(-2, 2))
				return
			}
			order := u.order
			u.order = orderNone
			u.waypoint = gmath.Vec{}
			u.completeOrder(order)
		}
	}

	switch u.stats {
	case droneHarvesterStats:
		u.updateHarvester(delta)
	case dronePatrolStats:
		u.updatePatrol(delta)
	case droneGeneratorStats:
		u.updateGenerator(delta)
	}
}

func (u *unitNode) completeOrder(order unitOrder) {
	switch order {
	case orderDig:
		u.completeDig()
	case orderHarvestResource:
		u.completeHarvestResource()
	case orderDeliverResource:
		u.completeDeliverResource()
	case orderPatrolMove:
		u.specialDelay = u.scene.Rand().FloatRange(7, 10)
	}
}

func (u *unitNode) completeDeliverResource() {
	target := u.orderTarget.(*unitNode)
	if target.IsDisposed() {
		u.order = orderDeliverResource
		return
	}
	if target.pos.DistanceSquaredTo(u.pos) > (64 * 64) {
		u.order = orderDeliverResource
		return
	}
	u.world.AddIron(u.cargo)
	u.cargo = 0
}

func (u *unitNode) completeHarvestResource() {
	res := u.orderTarget.(*resourceNode)
	if res.IsDisposed() || res.amount <= 0 {
		return
	}

	res.amount--
	if res.amount <= 0 {
		res.Dispose()
	}
	u.cargo = 1
	u.order = orderDeliverResource
}

func (u *unitNode) updateGenerator(delta float64) {
	if !u.waypoint.IsZero() {
		return
	}

	u.sendTo(randomSectorPos(u.scene.Rand(), u.world.diggedRect))
}

func (u *unitNode) updatePatrol(delta float64) {
	if !u.waypoint.IsZero() {
		return
	}

	u.specialDelay = gmath.ClampMin(u.specialDelay-delta, 0)
	if u.specialDelay != 0 {
		return
	}
	if u.scene.Rand().Chance(0.2) {
		u.specialDelay = u.scene.Rand().FloatRange(3, 5)
		return
	}

	u.sendTo(randomSectorPos(u.scene.Rand(), u.world.diggedRect))
	u.order = orderPatrolMove
}

func (u *unitNode) updateHarvester(delta float64) {
	switch u.order {
	case orderHarvestResource:
		// Moving towards the resource.

	case orderDeliverResource:
		u.specialDelay = gmath.ClampMin(u.specialDelay-delta, 0)
		if u.specialDelay != 0 {
			return
		}
		if u.cargo == 0 {
			u.order = orderNone
			return
		}
		if !u.waypoint.IsZero() {
			// Already delivering it somewhere.
			return
		}
		var closestPlace *unitNode
		closestDistSqr := math.MaxFloat64
		for _, other := range u.world.playerUnits {
			distMultiplier := 1.0
			switch other.stats {
			case droneCoreStats:
				distMultiplier = 1.5
			case buildingSmelter:
				// OK
			default:
				continue
			}
			distSqr := other.pos.DistanceSquaredTo(u.pos) * distMultiplier
			if distSqr < closestDistSqr {
				closestDistSqr = distSqr
				closestPlace = other
			}
		}
		if closestPlace != nil {
			if closestPlace.pos.DistanceSquaredTo(u.pos) < 16 {
				u.specialDelay = 2
				return
			}
			u.orderTarget = closestPlace
			u.sendTo(closestPlace.pos)
		}

	default:
		u.specialDelay = gmath.ClampMin(u.specialDelay-delta, 0)
		if u.specialDelay != 0 {
			return
		}
		closestDistSqr := math.MaxFloat64
		var closestResource *resourceNode
		for _, res := range u.world.resourceNodes {
			distSqr := res.pos.DistanceSquaredTo(u.pos)
			if distSqr >= (512 * 512) {
				continue
			}
			if distSqr < closestDistSqr {
				closestDistSqr = distSqr
				closestResource = res
			}
		}
		if closestResource != nil {
			u.order = orderHarvestResource
			u.orderTarget = closestResource
			u.sendTo(closestResource.pos)
		} else {
			u.specialDelay = u.scene.Rand().FloatRange(2, 7)
		}
	}

}

func (u *unitNode) movementSpeed() float64 {
	multiplier := 1.0
	if u.cargo != 0 {
		multiplier = 0.25
	}
	return u.stats.speed * multiplier
}

func (u *unitNode) completeDig() {
	m := u.orderTarget.(*mountainNode)
	if m.IsDisposed() {
		return
	}

	if u.world.energy < digEnergyCost {
		u.scene.AddObject(newFloatingTextNode(m.pos, "Error: not enough energy"))
		return
	}

	u.world.AddEnergy(-digEnergyCost)
	u.scene.AddObject(newFloatingTextNode(m.pos, "Status: dig complete"))
	u.world.AddStones(1)

	switch loot := u.world.PeekLoot(m); loot {
	case lootExtraStones:
		u.world.AddStones(2)
	case lootIronDeposit, lootLargeIronDeposit:
		minAmount := 4
		maxAmount := 8
		if loot == lootLargeIronDeposit {
			minAmount = 12
			maxAmount = 30
		}
		iron := u.world.NewResourceNode(m.pos, ironResourceStats, u.scene.Rand().IntRange(minAmount, maxAmount))
		u.scene.AddObjectBelow(iron, 1)
	case lootBotHarvester:
		u.scene.AddObject(u.world.NewUnitNode(m.pos, droneHarvesterStats))
	case lootBotPatrol:
		u.scene.AddObject(u.world.NewUnitNode(m.pos, dronePatrolStats))
	case lootBotGenerator:
		u.scene.AddObject(u.world.NewUnitNode(m.pos, droneGeneratorStats))
	case lootFlatCell:
		u.scene.AddObject(u.world.NewHardTerrainNode(m.pos))
	}

	// TODO: could be a plain tile.
	u.world.GrowDiggedRect(m.pos)
	u.world.grid.SetCellTile(u.world.grid.PosToCoord(m.pos.X, m.pos.Y), tileCaveMud)

	m.Dispose()
	delete(u.world.mountainByCoord, u.world.grid.PackCoord(u.world.grid.PosToCoord(m.pos.X, m.pos.Y)))
}

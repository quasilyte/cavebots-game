package battle

import (
	"math"

	"github.com/quasilyte/cavebots-game/assets"
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
	orderRepairOther
	orderVanguardFollow
	orderMakeUnit
	orderCreepAttack
	orderCreepAfterAttack
)

type unitNode struct {
	world  *worldState
	sprite *ge.Sprite
	anim   *ge.Animation

	stats *unitStats

	reload float64
	health float64

	scene *ge.Scene

	path        pathing.GridPath
	pathDest    gmath.Vec
	waypoint    gmath.Vec
	order       unitOrder
	orderTarget any

	specialDelay float64
	chargeDelay  float64

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
	switch u.stats.layer {
	case 0:
		u.world.stage.AddSpriteBelow(u.sprite)
	case 1:
		u.world.stage.AddSprite(u.sprite)
	case 2:
		u.world.stage.AddSpriteSlightlyAbove(u.sprite)
	case 3:
		u.world.stage.AddSpriteAbove(u.sprite)
	}

	if u.sprite.FrameWidth != u.sprite.ImageWidth() {
		u.anim = ge.NewRepeatedAnimation(u.sprite, -1)
		u.anim.Tick(scene.Rand().FloatRange(0.1, 4.6))
		u.anim.SetAnimationSpan(1)
	}

	switch u.stats {
	case creepMutantBase:
		u.reload = scene.Rand().FloatRange(10, 20)
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

func (u *unitNode) SendAttacking(pos gmath.Vec) {
	u.sendTo(pos)
	u.order = orderCreepAttack
}

func (u *unitNode) sendTo(pos gmath.Vec) {
	from := u.world.grid.PosToCoord(u.pos.X, u.pos.Y)
	to := u.world.grid.PosToCoord(pos.X, pos.Y)
	result := u.world.astar.BuildPath(u.world.grid, from, to, normalLayer)
	u.path = result.Steps

	u.pathDest = makePos(u.world.grid.AlignPos(pos.X, pos.Y))
	u.waypoint = makePos(u.world.grid.AlignPos(u.pos.X, u.pos.Y))
}

func (u *unitNode) OnDamage(amount float64, attacker *unitNode) {
	u.health = gmath.ClampMin(u.health-amount, 0)
	if u.health == 0 {
		u.destroy()
		return
	}

	switch u.order {
	case orderCreepAfterAttack:
		if u.waypoint.IsZero() && u.scene.Rand().Chance(0.7) {
			u.sendTo(attacker.pos)
		}
	}
}

func (u *unitNode) destroy() {
	if u.IsDisposed() {
		return
	}
	explodes := u.stats.allied || u.stats.building
	if explodes {
		playSound(u.world, assets.AudioExplosion1, u.pos)
		u.scene.AddObject(newEffectNode(u.world, u.pos, true, assets.ImageEffectExplosion))
	} else if u.stats.tiny {
		// Only mutant infantry is considered to be tiny.
		u.scene.AddObject(newEffectNode(u.world, u.pos, true, assets.ImageEffectMutantImpact))
	}
	u.Dispose()
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
				var offset gmath.Vec
				if u.stats.tiny {
					offset = u.world.rand.Offset(-10, 10)
				} else {
					offset = u.world.rand.Offset(-2, 2)
				}
				u.waypoint = nextPos.Add(offset)
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
	case droneVanguardStats:
		u.updateVanguard(delta)
	case droneTitanStats:
		u.updateTitan(delta)
	case droneGeneratorStats:
		u.updateGenerator(delta)
	case droneRepairStats:
		u.updateRepairbot(delta)
	case buildingFactory:
		u.updateFactory(delta)
	case buildingTurret:
		u.updateTurret(delta)
	case creepMutantBase:
		u.updateCreepBase(delta)
	case creepMutantWarrior:
		u.updateMutantWarrior(delta)
	case creepMutantHunter:
		u.updateMutantHunter(delta)
	case creepMutantGunner:
		u.updateMutantGunner(delta)
	case creepMutantWarlord:
		u.updateMutantWarlord(delta)
	case creepJeep:
		u.updateJeep(delta)
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
	case orderRepairOther:
		u.completeRepairOther()
	case orderCreepAttack, orderCreepAfterAttack:
		u.order = orderCreepAfterAttack
	case orderPatrolMove:
		u.specialDelay = u.scene.Rand().FloatRange(7, 10)
	case orderVanguardFollow:
		u.specialDelay = u.scene.Rand().FloatRange(3, 9)
	}
}

func (u *unitNode) completeRepairOther() {
	other := u.orderTarget.(*unitNode)
	if other.IsDisposed() {
		return
	}

	if other.pos.DistanceTo(u.pos) > 64 {
		if u.world.diggedRect.Contains(other.pos) && u.world.rand.Chance(0.6) {
			u.sendTo(other.pos)
			u.order = orderRepairOther
			return
		}
		u.reload = u.world.rand.FloatRange(1, 5)
		return
	}

	u.reload = u.world.rand.FloatRange(10, 30)
	other.health = gmath.ClampMax(other.health+15, other.stats.maxHealth)
	playSound(u.world, assets.AudioRepair, other.pos)
	u.scene.AddObject(newEffectNode(u.world, other.pos.Sub(gmath.Vec{Y: 8}), true, assets.ImageEffectRepair))
}

func (u *unitNode) completeDeliverResource() {
	target := u.orderTarget.(*unitNode)
	if target.IsDisposed() {
		u.order = orderDeliverResource
		return
	}
	if target.stats == droneCoreStats {
		if target.pos.DistanceSquaredTo(u.pos) > (64 * 64) {
			u.order = orderDeliverResource
			return
		}
	} else {
		if target.pos.DistanceSquaredTo(u.pos) > (10 * 10) {
			u.waypoint = target.pos.Add(u.scene.Rand().Offset(-2, 2))
			u.order = orderDeliverResource
			return
		}
	}

	if u.cargo == 0 {
		return
	}

	u.world.AddIron(u.cargo)
	u.cargo = 0
	playGlobalSound(u.world, assets.AudioResourceAdded)
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
	u.specialDelay = u.scene.Rand().FloatRange(0.4, 0.8)
}

func (u *unitNode) maybeCharge(delta, maxDist float64) {
	u.chargeDelay -= delta
	if u.chargeDelay > 0 {
		return
	}

	if u.world.caveRect.Contains(u.pos) {
		maxDist *= 3
	}

	target := gmath.RandElem(u.world.rand, u.world.playerUnits)
	if target == nil {
		u.chargeDelay = u.scene.Rand().FloatRange(8, 12)
		return
	}
	if target.pos.DistanceTo(u.pos) > maxDist {
		u.chargeDelay = u.scene.Rand().FloatRange(1, 3.5)
		return
	}
	u.chargeDelay = u.scene.Rand().FloatRange(4, 15)
	u.sendTo(target.pos)
}

func (u *unitNode) updateMutantWarrior(delta float64) {
	u.processWeapon(delta)
	if u.waypoint.IsZero() {
		u.maybeCharge(delta, 140)
	}
}

func (u *unitNode) updateMutantHunter(delta float64) {
	u.processWeapon(delta)
}

func (u *unitNode) updateMutantGunner(delta float64) {
	u.processWeapon(delta)
	if u.waypoint.IsZero() && u.order == orderCreepAfterAttack {
		u.maybeCharge(delta, 640)
	}
}

func (u *unitNode) updateMutantWarlord(delta float64) {
	u.processWeapon(delta)
	if u.waypoint.IsZero() {
		u.maybeCharge(delta, 200)
	}
}

func (u *unitNode) updateJeep(delta float64) {
	u.processWeapon(delta)
}

func (u *unitNode) processWeapon(delta float64) {
	u.reload = gmath.ClampMin(u.reload-delta, 0)
	if u.reload != 0 {
		return
	}

	attackDistSqr := u.stats.weapon.attackRange * u.stats.weapon.attackRange
	isMelee := attackDistSqr == 0
	if isMelee {
		attackDistSqr = 34.0 * 34.0
	}

	u.world.tmpTargetsSlice = u.world.tmpTargetsSlice[:0]

	targetSlice := u.world.playerUnits
	if u.stats.allied {
		targetSlice = u.world.creeps
	}

	randIterate(u.world.rand, targetSlice, func(other *unitNode) bool {
		if other.pos.DistanceSquaredTo(u.pos) > attackDistSqr {
			return false
		}
		if !isMelee && !hasLineOfFire(u.world, u.pos, other.pos) {
			return false
		}
		u.world.tmpTargetsSlice = append(u.world.tmpTargetsSlice, other)
		return len(u.world.tmpTargetsSlice) >= u.stats.weapon.maxTargets
	})

	if len(u.world.tmpTargetsSlice) == 0 {
		if isMelee {
			u.reload = u.scene.Rand().FloatRange(0.4, 1.1)
		} else {
			u.reload = u.scene.Rand().FloatRange(1.5, 3.5)
		}
		return
	}

	u.reload = u.stats.weapon.reload * u.scene.Rand().FloatRange(0.8, 1.2)
	if isMelee {
		u.world.tmpTargetsSlice[0].OnDamage(u.stats.weapon.damage, u)
		playSound(u.world, u.stats.weapon.impactSound, u.world.tmpTargetsSlice[0].pos)
		return
	}

	for _, target := range u.world.tmpTargetsSlice {
		if u.stats.weapon.projectileImage != assets.ImageNone {
			for i := 0; i < u.stats.weapon.burstSize; i++ {
				fireDelay := float64(i) * u.stats.weapon.burstDelay
				projectile := newProjectileNode(projectileNodeConfig{
					attacker:  u,
					target:    target,
					targetPos: target.pos.Add(u.scene.Rand().Offset(-6, 6)),
					fireDelay: fireDelay,
				})
				u.scene.AddObject(projectile)
			}
		} else {
			from := ge.Pos{Base: &u.pos}
			to := ge.Pos{Base: &target.pos}
			line := newLineNode(u.world, from, to, u.stats.weapon.beamColor)
			u.scene.AddObject(line)
			target.OnDamage(u.stats.weapon.damage, u)
		}
	}
	playSound(u.world, u.stats.weapon.fireSound, u.pos)
}

func (u *unitNode) updateCreepBase(delta float64) {
	u.world.creepBaseLevel += delta * u.world.creepsEvolutionRate

	u.reload = gmath.ClampMin(u.reload-delta, 0)
	if u.reload != 0 {
		return
	}

	if len(u.world.creeps) > 200 {
		u.reload = 15
		return
	}
	if len(u.world.creeps) > 120 && u.scene.Rand().Chance(0.6) {
		u.reload = u.scene.Rand().FloatRange(5, 10)
		return
	}

	u.reload = u.scene.Rand().FloatRange(20, 30)
	waypoint := u.pos.Add(gmath.Vec{
		X: u.scene.Rand().FloatRange(-32, 32),
		Y: u.scene.Rand().FloatRange(32, 160),
	})
	maxWarriors := 5
	minWarriors := 2
	if u.world.difficulty == 2 {
		minWarriors = 3
		maxWarriors = 6
	}
	numWarriors := u.scene.Rand().IntRange(minWarriors, maxWarriors)
	for i := 0; i < numWarriors; i++ {
		stats := creepMutantWarrior
		// Every minute gives +8% warlord chance.
		warlardChance := 0.0
		if u.world.creepBaseLevel > 60 {
			warlardChance = 0.08 * (u.world.creepBaseLevel / 60)
		}
		if warlardChance != 0 && u.scene.Rand().Chance(warlardChance) {
			stats = creepMutantWarlord
		}
		newUnit := u.world.NewUnitNode(u.pos.Add(u.scene.Rand().Offset(-8, 8)), stats)
		u.scene.AddObject(newUnit)
		newUnit.SendTo(waypoint)
	}

	maxArchers := u.world.difficulty + 2
	numArchers := gmath.Clamp(int(u.world.creepBaseLevel/100.0), 0, maxArchers)
	for i := 0; i < numArchers; i++ {
		newUnit := u.world.NewUnitNode(u.pos.Add(u.scene.Rand().Offset(-8, 8)), creepMutantHunter)
		u.scene.AddObject(newUnit)
		newUnit.SendTo(waypoint)
	}

	if u.world.creepBaseLevel > 300 {
		maxGunners := u.world.difficulty + 2
		numGunners := gmath.Clamp(int(u.world.creepBaseLevel/300.0), 0, maxGunners)
		for i := 0; i < numGunners; i++ {
			newUnit := u.world.NewUnitNode(u.pos.Add(u.scene.Rand().Offset(-8, 8)), creepMutantGunner)
			u.scene.AddObject(newUnit)
			newUnit.SendTo(waypoint)
		}
	}
}

func (u *unitNode) updateTurret(delta float64) {
	u.processWeapon(delta)
}

func (u *unitNode) updateFactory(delta float64) {
	if u.order != orderMakeUnit {
		return
	}

	u.specialDelay = gmath.ClampMin(u.specialDelay-delta, 0)
	if u.specialDelay == 0 {
		u.specialDelay = u.world.rand.FloatRange(0.7, 1.4)
		u.scene.AddObject(newEffectNode(u.world, u.pos.Sub(gmath.Vec{Y: 16}), true, assets.ImageEffectSmokeUp))
	}

	u.reload -= delta
	if u.reload > 0 {
		return
	}

	u.reload = 0
	stats := u.orderTarget.(*unitStats)
	u.orderTarget = nil
	u.order = orderNone

	u.scene.AddObject(u.world.NewUnitNode(u.pos, stats))
	u.world.results.BotsCreated++
	playGlobalSound(u.world, assets.AudioUnitReady)
	u.world.EventTooltipUpdateRequest.Emit(gsignal.Void{})
}

func (u *unitNode) updateRepairbot(delta float64) {
	u.reload = gmath.ClampMin(u.reload-delta, 0)
	if u.order != orderNone {
		return
	}
	if !u.waypoint.IsZero() {
		return
	}
	var target *unitNode
	if u.reload == 0 {
		target = randIterate(u.world.rand, u.world.playerUnits, func(other *unitNode) bool {
			if other == u {
				return false
			}
			if other.stats.building {
				return false
			}
			if other.health >= other.stats.maxHealth {
				return false
			}
			if !u.world.diggedRect.Contains(other.pos) {
				return false
			}
			return true
		})
	}
	if target == nil {
		u.sendTo(randomSectorPos(u.scene.Rand(), u.world.diggedRect))
		return
	}
	u.order = orderRepairOther
	u.orderTarget = target
	u.sendTo(target.pos)
}

func (u *unitNode) updateGenerator(delta float64) {
	if !u.waypoint.IsZero() {
		return
	}

	u.sendTo(randomSectorPos(u.scene.Rand(), u.world.diggedRect))
}

func (u *unitNode) updateTitan(delta float64) {
	u.processWeapon(delta)

	if !u.waypoint.IsZero() {
		return
	}
	if u.world.creepBase == nil || u.world.creepBase.IsDisposed() {
		return
	}

	u.SendTo(u.world.creepBase.pos.Add(u.scene.Rand().Offset(-140, 140)))
}

func (u *unitNode) updateVanguard(delta float64) {
	u.processWeapon(delta)

	if !u.waypoint.IsZero() {
		return
	}
	if u.world.core == nil || u.world.core.IsDisposed() {
		return
	}

	u.specialDelay = gmath.ClampMin(u.specialDelay-delta, 0)
	if u.specialDelay != 0 {
		return
	}

	u.sendTo(u.world.core.pos.Add(u.scene.Rand().Offset(-80, 80)))
	u.order = orderVanguardFollow
}

func (u *unitNode) updatePatrol(delta float64) {
	u.processWeapon(delta)

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
				distMultiplier = 1.6
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
			return
		}
		if u.scene.Rand().Chance(0.6) {
			u.specialDelay = u.scene.Rand().FloatRange(2, 7)
		} else {
			u.SendTo(randomSectorPos(u.scene.Rand(), u.world.diggedRect))
		}
	}
}

func (u *unitNode) movementSpeed() float64 {
	multiplier := 1.0
	if u.cargo != 0 {
		multiplier = 0.25
	}
	if u.order == orderRepairOther {
		multiplier += 0.75
	}
	return u.stats.speed * multiplier
}

func (u *unitNode) completeDig() {
	m := u.orderTarget.(*mountainNode)
	if m.IsDisposed() {
		return
	}
	if u.pos.DistanceTo(u.pathDest) > 40 {
		return
	}

	energyCost := float64(digEnergyCost)
	if m.loot == lootEasyDig {
		energyCost = 0
	}

	if u.world.energy < energyCost {
		playGlobalSound(u.world, assets.AudioError)
		u.scene.AddObject(newFloatingTextNode(m.world, m.pos, "Error: not enough energy"))
		u.world.notEnoughEnergy += 1
		return
	}

	u.world.AddEnergy(-energyCost)
	u.scene.AddObject(newFloatingTextNode(m.world, m.pos, "Status: dig complete"))
	u.world.AddStones(1)

	tileType := uint8(tileCaveMud)

	createdBot := false

	switch m.loot {
	case lootExtraStones:
		u.world.AddStones(2)
	case lootIronDeposit, lootLargeIronDeposit:
		minAmount := 4
		maxAmount := 8
		if m.loot == lootLargeIronDeposit {
			minAmount = 12
			maxAmount = 30
		}
		iron := u.world.NewResourceNode(m.pos, ironResourceStats, u.scene.Rand().IntRange(minAmount, maxAmount))
		u.scene.AddObjectBelow(iron, 1)

	case lootBotHarvester:
		createdBot = true
		u.scene.AddObject(u.world.NewUnitNode(m.pos, droneHarvesterStats))
	case lootBotPatrol:
		createdBot = true
		u.scene.AddObject(u.world.NewUnitNode(m.pos, dronePatrolStats))
	case lootBotGenerator:
		createdBot = true
		u.scene.AddObject(u.world.NewUnitNode(m.pos, droneGeneratorStats))
	case lootBotRepair:
		createdBot = true
		u.scene.AddObject(u.world.NewUnitNode(m.pos, droneRepairStats))
	case lootBotVanguard:
		createdBot = true
		u.scene.AddObject(u.world.NewUnitNode(m.pos, droneVanguardStats))
	case lootBotTitan:
		createdBot = true
		u.scene.AddObject(u.world.NewUnitNode(m.pos, droneTitanStats))

	case lootFlatCell:
		u.scene.AddObject(u.world.NewHardTerrainNode(m.pos))
		tileType = tileCaveFlat
	}

	if createdBot {
		u.world.results.BotsCreated++
		playGlobalSound(u.world, assets.AudioUnitReady)
	}

	u.world.GrowDiggedRect(m.pos)
	u.world.grid.SetCellTile(u.world.grid.PosToCoord(m.pos.X, m.pos.Y), tileType)

	playSound(u.world, assets.AudioDig, m.pos)
	m.Dispose()
	delete(u.world.mountainByCoord, u.world.grid.PackCoord(u.world.grid.PosToCoord(m.pos.X, m.pos.Y)))
	u.world.RevealNeighbors(m.pos)
	u.world.results.Digs++
	u.world.EventTooltipUpdateRequest.Emit(gsignal.Void{})
}

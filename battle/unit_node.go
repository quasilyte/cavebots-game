package battle

import (
	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/pathing"
)

type unitOrder int

const (
	orderNone unitOrder = iota
	orderDig
)

type unitNode struct {
	world  *worldState
	sprite *ge.Sprite

	stats *unitStats

	scene *ge.Scene

	path        pathing.GridPath
	pathDest    gmath.Vec
	waypoint    gmath.Vec
	order       unitOrder
	orderTarget any

	pos gmath.Vec
}

func newUnitNode(world *worldState, stats *unitStats) *unitNode {
	return &unitNode{
		world: world,
		stats: stats,
	}
}

func (u *unitNode) Init(scene *ge.Scene) {
	u.scene = scene

	u.sprite = scene.NewSprite(assets.ImageDroneCore)
	u.sprite.Pos.Base = &u.pos
	scene.AddGraphics(u.sprite)
}

func (u *unitNode) IsDisposed() bool {
	return u.sprite.IsDisposed()
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
	if !u.waypoint.IsZero() {
		newPos, reached := moveTowardsWithSpeed(u.pos, u.waypoint, delta, u.stats.speed)
		u.pos = newPos
		if reached {
			if u.path.HasNext() {
				nextPos := nextPathWaypoint(u.world, u.pos, &u.path, normalLayer)
				u.waypoint = nextPos.Add(u.world.rand.Offset(-2, 2))
				return
			}
			order := u.order
			u.order = orderNone
			if order == orderDig {
				u.completeDig()
			}
		}
	}
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
	u.world.AddStones(1)
	u.scene.AddObject(newFloatingTextNode(m.pos, "Status: dig complete"))
	// TODO: could be a plain tile.
	u.world.grid.SetCellTile(u.world.grid.PosToCoord(m.pos.X, m.pos.Y), tileCaveMud)
	m.Dispose()
	delete(u.world.mountainByCoord, u.world.grid.PackCoord(u.world.grid.PosToCoord(m.pos.X, m.pos.Y)))
}

package battle

import (
	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/pathing"
)

type unitNode struct {
	world  *worldState
	sprite *ge.Sprite

	stats *unitStats

	path     pathing.GridPath
	pathDest gmath.Vec
	waypoint gmath.Vec

	pos gmath.Vec
}

func newUnitNode(world *worldState, stats *unitStats) *unitNode {
	return &unitNode{
		world: world,
		stats: stats,
	}
}

func (u *unitNode) Init(scene *ge.Scene) {
	u.sprite = scene.NewSprite(assets.ImageDroneCore)
	u.sprite.Pos.Base = &u.pos
	scene.AddGraphics(u.sprite)
}

func (u *unitNode) IsDisposed() bool {
	return u.sprite.IsDisposed()
}

func (u *unitNode) SendTo(pos gmath.Vec) {
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
		}
	}
}

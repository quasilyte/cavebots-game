package battle

import (
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/quasilyte/cavebots-game/assets"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/pathing"
	"golang.org/x/image/font"
)

var firstBuildingList = []*unitStats{
	buildingPowerGenerator,
	buildingSmelter,
}

var secondBuildingList = []*unitStats{
	buildingBarricate,
}

var buildingHotkeys = []string{
	"Q",
	"W",
}

func randomSectorPos(rng *gmath.Rand, sector gmath.Rect) gmath.Vec {
	return gmath.Vec{
		X: rng.FloatRange(sector.Min.X, sector.Max.X),
		Y: rng.FloatRange(sector.Min.Y, sector.Max.Y),
	}
}

var cellNeighborOffsets = []pathing.GridCoord{
	{X: 1},
	{Y: 1},
	{X: -1},
	{Y: -1},
}

type tinyCoord struct {
	data uint8
}

func makeTinyCoord(x, y int) tinyCoord {
	xsign := uint8(x & (1 << 7))
	xvalue := uint8(x&1) << 6
	ysign := uint8(y&(1<<7)) >> 2
	yvalue := uint8(y&1) << 4
	data := (xsign | xvalue | ysign | yvalue)
	return tinyCoord{data: data}
}

func (c tinyCoord) X() int8 {
	return int8(c.data&0b11000000) >> 6
}

func (c tinyCoord) Y() int8 {
	return int8((c.data<<2)&0b11000000) >> 6
}

func (c tinyCoord) IsZero() bool { return c.data == 0 }

var diagonalMoveTab = [32]tinyCoord{
	pathing.DirRight | (pathing.DirUp << 2):   makeTinyCoord(1, -1),
	pathing.DirRight | (pathing.DirDown << 2): makeTinyCoord(1, 1),
	pathing.DirDown | (pathing.DirRight << 2): makeTinyCoord(1, 1),
	pathing.DirDown | (pathing.DirLeft << 2):  makeTinyCoord(-1, 1),
	pathing.DirLeft | (pathing.DirDown << 2):  makeTinyCoord(-1, 1),
	pathing.DirLeft | (pathing.DirUp << 2):    makeTinyCoord(-1, -1),
	pathing.DirUp | (pathing.DirLeft << 2):    makeTinyCoord(-1, -1),
	pathing.DirUp | (pathing.DirRight << 2):   makeTinyCoord(1, -1),
}

func makePos(x, y float64) gmath.Vec {
	return gmath.Vec{X: x, Y: y}
}

func moveTowardsWithSpeed(from, to gmath.Vec, delta, speed float64) (gmath.Vec, bool) {
	travelled := speed * delta
	result := from.MoveTowards(to, travelled)
	return result, result == to
}

func nextPathWaypoint(world *worldState, pos gmath.Vec, p *pathing.GridPath, l pathing.GridLayer) gmath.Vec {
	cell := world.grid.PosToCoord(pos.X, pos.Y)
	d1, d2 := p.Peek2()
	diagOffset := diagonalMoveTab[(d1|(d2<<2))&0b11111]
	if !diagOffset.IsZero() {
		// Can make a diagonal move, if second cell is free too.
		otherCell := cell.Move(d2)
		if world.grid.GetCellCost(otherCell, l) != 0 {
			p.Skip(2)
			return makePos(world.grid.CoordToPos(cell.Move(d1).Move(d2)))
		}
	}
	p.Skip(1)
	return makePos(world.grid.CoordToPos(cell.Move(d1)))
}

func estimateMessageBounds(f font.Face, s string, xpadding float64) (width, height float64) {
	bounds := text.BoundString(f, s)
	width = (float64(bounds.Dx()) + 16) + xpadding
	height = (float64(bounds.Dy()) + 16)
	return width, height
}

func playGlobalSound(world *worldState, id resource.AudioID) {
	numSamples := assets.NumSamples(id)
	if numSamples == 1 {
		world.scene.Audio().PlaySound(id)
	} else {
		soundIndex := world.rand.IntRange(0, numSamples-1)
		sound := resource.AudioID(int(id) + soundIndex)
		world.scene.Audio().PlaySound(sound)
	}
}

func playSound(world *worldState, id resource.AudioID, pos gmath.Vec) {
	if world.camera.ContainsPos(pos) {
		playGlobalSound(world, id)
	}
}

func randIterate[T any](rand *gmath.Rand, slice []T, f func(x T) bool) T {
	var result T
	if len(slice) == 0 {
		return result
	}
	if len(slice) == 1 {
		// Don't use rand() if there is only 1 element.
		x := slice[0]
		if f(x) {
			result = x
		}
		return result
	}

	var slider gmath.Slider
	slider.SetBounds(0, len(slice)-1)
	slider.TrySetValue(rand.IntRange(0, len(slice)-1))
	inc := rand.Bool()
	for i := 0; i < len(slice); i++ {
		x := slice[slider.Value()]
		if inc {
			slider.Inc()
		} else {
			slider.Dec()
		}
		if f(x) {
			result = x
			break
		}
	}
	return result
}

func hasLineOfFire(world *worldState, from, to gmath.Vec) bool {
	dist := from.DistanceTo(to)
	if dist <= 40 {
		return true
	}
	pos := from
	for dist > 30 {
		pos = pos.MoveTowards(to, 28)
		coord := world.grid.PosToCoord(pos.X, pos.Y)
		if world.grid.GetCellTile(coord) == tileBlocked {
			return false
		}
		dist -= 30
	}
	return true
}

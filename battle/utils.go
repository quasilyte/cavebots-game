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

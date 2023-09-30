package battle

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/pathing"
)

type worldState struct {
	width  float64
	height float64

	scene *ge.Scene
	rand  *gmath.Rand

	rect          gmath.Rect
	caveRect      gmath.Rect
	innerCaveRect gmath.Rect
	diggedRect    gmath.Rect

	grid  *pathing.Grid
	astar *pathing.AStar

	energy float64
	iron   int
	stones int

	lootSeq int

	caveWidth float64

	core        *unitNode
	playerUnits []*unitNode
	hardTerrain []*hardTerrainNode

	resourceNodes []*resourceNode

	mountainByCoord map[uint32]*mountainNode

	EventResourcesUpdated gsignal.Event[gsignal.Void]
}

func (w *worldState) Init() {
	w.energy = 15

	w.rect = gmath.Rect{
		Max: gmath.Vec{
			X: w.width,
			Y: w.height,
		},
	}
	w.caveRect = gmath.Rect{
		Max: gmath.Vec{
			X: w.caveWidth,
			Y: w.rect.Max.Y,
		},
	}
	w.innerCaveRect = gmath.Rect{
		Min: gmath.Vec{X: 32, Y: 32},
		Max: w.caveRect.Max.Sub(gmath.Vec{X: 32, Y: 32}),
	}
	w.mountainByCoord = make(map[uint32]*mountainNode)

	w.grid = pathing.NewGrid(pathing.GridConfig{
		WorldWidth:  uint(w.width),
		WorldHeight: uint(w.height),
		CellWidth:   32,
		CellHeight:  32,
	})
	w.astar = pathing.NewAStar(pathing.AStarConfig{
		NumCols: uint(w.grid.NumCols()),
		NumRows: uint(w.grid.NumRows()),
	})
}

func (w *worldState) GrowDiggedRect(pos gmath.Vec) {
	if w.diggedRect.Min.X > pos.X {
		w.diggedRect.Min.X = pos.X
	}
	if w.diggedRect.Max.X < pos.X {
		w.diggedRect.Max.X = pos.X
	}

	if w.diggedRect.Min.Y > pos.Y {
		w.diggedRect.Min.Y = pos.Y
	}
	if w.diggedRect.Max.Y < pos.Y {
		w.diggedRect.Max.Y = pos.Y
	}
}

func (w *worldState) NewHardTerrainNode(pos gmath.Vec) *hardTerrainNode {
	n := newHardTerrainNode(pos)
	w.hardTerrain = append(w.hardTerrain, n)
	return n
}

func (w *worldState) NewUnitNode(pos gmath.Vec, stats *unitStats) *unitNode {
	n := newUnitNode(w, stats)
	n.pos = pos
	n.EventDisposed.Connect(nil, func(n *unitNode) {
		if n.stats.allied {
			w.playerUnits = xslices.Remove(w.playerUnits, n)
		}
	})
	if stats.allied {
		w.playerUnits = append(w.playerUnits, n)
	}
	return n
}

func (w *worldState) NewResourceNode(pos gmath.Vec, stats *resourceStats, amount int) *resourceNode {
	n := newResourceNode(pos, stats, amount)
	n.EventDisposed.Connect(nil, func(n *resourceNode) {
		w.resourceNodes = xslices.Remove(w.resourceNodes, n)
	})
	w.resourceNodes = append(w.resourceNodes, n)
	return n
}

func (w *worldState) MountainAt(pos gmath.Vec) *mountainNode {
	coord := w.grid.PosToCoord(pos.X, pos.Y)
	packedCoord := w.grid.PackCoord(coord)
	if m, ok := w.mountainByCoord[packedCoord]; ok {
		return m
	}
	return nil
}

func (w *worldState) PeekLoot(m *mountainNode) lootKind {
	if m.loot == lootUnknown {
		w.AssignLoot(m)
	}
	return m.loot
}

func (w *worldState) AssignLoot(m *mountainNode) {
	m.loot = w.selectLootKind()
	w.lootSeq++
}

func (w *worldState) selectLootKind() lootKind {
	// Since this is a game prototype and our random loot
	// generation is not that reliable, hardcode
	// some of the items here.
	switch w.lootSeq {
	case 0:
		return lootBotHarvester
	case 3:
		return lootFlatCell
	case 8:
		return lootBotPatrol
	}

	if w.lootSeq%5 == 0 {
		roll := w.rand.Float()
		if roll <= 0.4 {
			return lootBotHarvester
		}
		if roll <= 0.6 {
			return lootBotPatrol
		}
	}

	if w.lootSeq%3 == 0 {
		if w.rand.Chance(0.5) {
			return lootIronDeposit
		}
		if w.lootSeq > 25 && w.rand.Chance(0.25) {
			return lootLargeIronDeposit
		}
	}

	if w.lootSeq%4 == 0 {
		if w.rand.Chance(0.4) {
			return lootFlatCell
		}
	}

	if w.rand.Chance(0.1) {
		return lootExtraStones
	}

	return lootNone
}

func (w *worldState) CanDig(m *mountainNode) bool {
	if m.outer {
		return false
	}
	coord := w.grid.PosToCoord(m.pos.X, m.pos.Y)
	for _, offset := range cellNeighborOffsets {
		probe := coord.Add(offset)
		if probe.X < 0 || probe.Y < 0 {
			continue
		}
		if w.grid.GetCellTile(probe) != tileBlocked {
			return true
		}
	}
	return false
}

func (w *worldState) CalcEnergyUpkeep() float64 {
	v := 0.0
	for _, u := range w.playerUnits {
		if u.offline {
			continue
		}
		v += u.stats.energyUpkeep
	}
	return v
}

func (w *worldState) CalcEnergyRegen() float64 {
	regen := 1.0 // Base regen (core-provided)
	return regen - w.CalcEnergyUpkeep()
}

func (w *worldState) AddEnergy(delta float64) {
	if delta == 0 {
		return
	}
	w.energy += delta
	w.EventResourcesUpdated.Emit(gsignal.Void{})
}

func (w *worldState) AddIron(delta int) {
	if delta == 0 {
		return
	}
	w.iron += delta
	w.EventResourcesUpdated.Emit(gsignal.Void{})
}

func (w *worldState) AddStones(delta int) {
	if delta == 0 {
		return
	}
	w.stones += delta
	w.EventResourcesUpdated.Emit(gsignal.Void{})
}

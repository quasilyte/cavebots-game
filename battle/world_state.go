package battle

import (
	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/cavebots-game/viewport"
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

	stage  *viewport.Stage
	camera *viewport.Camera

	rect          gmath.Rect
	caveRect      gmath.Rect
	innerCaveRect gmath.Rect
	diggedRect    gmath.Rect

	grid  *pathing.Grid
	astar *pathing.AStar

	energy float64
	iron   int
	stones int

	lootSeq     int
	buildingSeq int

	caveWidth float64

	results *Results

	core        *unitNode
	playerUnits []*unitNode
	creeps      []*unitNode
	hardTerrain []*hardTerrainNode

	tmpTargetsSlice []*unitNode

	creepBase             *unitNode
	creepBaseLevel        float64
	creepBaseAttackBudget int
	creepsEvolutionRate   float64
	difficulty            int

	tutorial        bool
	notEnoughEnergy float64

	resourceNodes []*resourceNode

	mountainByCoord map[uint32]*mountainNode

	EventResourcesUpdated     gsignal.Event[gsignal.Void]
	EventDefeat               gsignal.Event[gmath.Vec]
	EventVictory              gsignal.Event[gmath.Vec]
	EventTooltipUpdateRequest gsignal.Event[gsignal.Void]
}

func (w *worldState) Init() {
	w.results = &Results{}

	switch w.difficulty {
	case 0:
		w.creepsEvolutionRate = 0.8
	case 1:
		w.creepsEvolutionRate = 1.0
	case 2:
		w.creepsEvolutionRate = 1.25
	}

	w.tmpTargetsSlice = make([]*unitNode, 0, 10)

	w.energy = 20

	w.creepBaseAttackBudget = 5

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
		DefaultTile: tileGrass,
	})
	w.astar = pathing.NewAStar(pathing.AStarConfig{
		NumCols: uint(w.grid.NumCols()),
		NumRows: uint(w.grid.NumRows()),
	})

	w.stage = viewport.NewStage()
	w.camera = viewport.NewCamera(w.stage, gmath.Rect{
		Max: gmath.Vec{
			X: 1920,
			Y: 1080,
		},
	}, 1920.0/2, 1080.0/2)
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
	var buildOptions [2]*unitStats
	switch w.buildingSeq {
	case 0:
		buildOptions[0] = buildingPowerGenerator
		buildOptions[1] = buildingBarricate
	case 2:
		buildOptions[0] = buildingPowerGenerator
		buildOptions[1] = buildingSmelter
	default:
		buildOptions[0] = gmath.RandElem(w.rand, firstBuildingList)
		buildOptions[1] = gmath.RandElem(w.rand, secondBuildingList)
	}
	w.buildingSeq++

	n := newHardTerrainNode(w, pos, buildOptions)
	w.hardTerrain = append(w.hardTerrain, n)
	return n
}

func (w *worldState) TryBuy(stats *unitStats, pos gmath.Vec) bool {
	if w.energy < float64(stats.energyCost) {
		playGlobalSound(w, assets.AudioError)
		w.scene.AddObject(newFloatingTextNode(w, pos, "Error: not enough energy"))
		w.notEnoughEnergy += 1
		return false
	}
	if w.iron < stats.ironCost {
		playGlobalSound(w, assets.AudioError)
		w.scene.AddObject(newFloatingTextNode(w, pos, "Error: not enough iron"))
		return false
	}
	if w.stones < stats.stoneCost {
		playGlobalSound(w, assets.AudioError)
		w.scene.AddObject(newFloatingTextNode(w, pos, "Error: not enough stone"))
		return false
	}

	w.AddEnergy(-float64(stats.energyCost))
	w.AddIron(-stats.ironCost)
	w.AddStones(-stats.stoneCost)
	return true
}

func (w *worldState) NewUnitNode(pos gmath.Vec, stats *unitStats) *unitNode {
	n := newUnitNode(w, stats)
	n.pos = pos
	if stats.building {
		w.grid.SetCellTile(w.grid.PosToCoord(pos.X, pos.Y), tileBlocked)
	}
	n.EventDisposed.Connect(nil, func(n *unitNode) {
		if stats.building {
			w.grid.SetCellTile(w.grid.PosToCoord(pos.X, pos.Y), tileCaveFlat)
		}
		if n.stats.allied {
			if n == w.core {
				w.core = nil
				w.EventDefeat.Emit(n.pos)
			}
			w.playerUnits = xslices.Remove(w.playerUnits, n)
		} else {
			if n == w.creepBase {
				w.EventVictory.Emit(n.pos)
			}
			w.creeps = xslices.Remove(w.creeps, n)
		}
	})
	if stats.allied {
		w.playerUnits = append(w.playerUnits, n)
	} else {
		w.creeps = append(w.creeps, n)
	}
	return n
}

func (w *worldState) NewResourceNode(pos gmath.Vec, stats *resourceStats, amount int) *resourceNode {
	n := newResourceNode(w, pos, stats, amount)
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
	case 6:
		return lootIronDeposit
	case 8:
		return lootBotPatrol
	case 11:
		return lootBotGenerator
	case 20:
		return lootBotVanguard
	case 77:
		return lootBotTitan
	case 89:
		return lootBotRepair
	}

	if w.lootSeq%5 == 0 {
		roll := w.rand.Float()
		if roll <= 0.4 { // 40%
			generatorChance := 0.0
			switch {
			case w.lootSeq > 50:
				generatorChance = 0.6
			case w.lootSeq > 30:
				generatorChance = 0.3
			}
			if generatorChance != 0 && w.rand.Chance(generatorChance) {
				return lootBotGenerator
			}
			return lootBotHarvester
		}
		if roll <= 0.6 { // 20%
			return lootBotPatrol
		}
		if roll <= 0.7 { // 10%
			return lootBotGenerator
		}
		if w.lootSeq > 32 && roll < 0.75 { // 5%
			return lootBotRepair
		}
	}

	if w.lootSeq%3 == 0 {
		if w.lootSeq > 25 && w.rand.Chance(0.4) {
			return lootLargeIronDeposit
		}
		// A second roll for the late game.
		if w.lootSeq > 40 && w.rand.Chance(0.3) {
			return lootLargeIronDeposit
		}
		if w.rand.Chance(0.5) {
			return lootIronDeposit
		}
	}

	if w.lootSeq%4 == 0 {
		if w.rand.Chance(0.55) {
			return lootFlatCell
		}
	}

	if w.rand.Chance(0.15) {
		return lootExtraStones
	}

	if w.rand.Chance(0.45) {
		return lootEasyDig
	}

	return lootNone
}

func (w *worldState) RevealNeighbors(pos gmath.Vec) {
	coord := w.grid.PosToCoord(pos.X, pos.Y)
	for _, offset := range cellNeighborOffsets {
		probe := coord.Add(offset)
		packedCoord := w.grid.PackCoord(probe)
		mountain := w.mountainByCoord[packedCoord]
		if mountain == nil {
			continue
		}
		w.Reveal(mountain)
	}
}

func (w *worldState) Reveal(m *mountainNode) {
	if m.loot != lootUnknown || !w.CanDig(m) {
		return
	}

	w.AssignLoot(m)

	m.sprite.SetColorScale(ge.ColorScale{R: 1, G: 1, B: 1, A: 1})

	switch m.loot {
	case lootFlatCell:
		m.sprite.SetImage(w.scene.LoadImage(assets.ImageHardMountains))
	case lootIronDeposit, lootLargeIronDeposit:
		m.sprite.SetImage(w.scene.LoadImage(assets.ImageIronMountains))
	case lootEasyDig:
		m.sprite.SetImage(w.scene.LoadImage(assets.ImageWeakMountains))
	case lootBotGenerator, lootBotPatrol, lootBotVanguard, lootBotHarvester, lootBotRepair, lootBotTitan:
		m.sprite.SetImage(w.scene.LoadImage(assets.ImageUnitMountains))
	case lootExtraStones:
		m.sprite.SetImage(w.scene.LoadImage(assets.ImageRockyMountains))
	}
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
	generatorMultiplier := 1.0
	generatorBotMultiplier := 1.0
	for _, u := range w.playerUnits {
		switch u.stats {
		case buildingPowerGenerator:
			regen += 1.1 * generatorMultiplier
			generatorMultiplier = gmath.ClampMin(generatorMultiplier-0.1, 0.3)
		case droneGeneratorStats:
			regen += 0.6 * generatorBotMultiplier
			generatorBotMultiplier = gmath.ClampMin(generatorBotMultiplier-0.1, 0.1)
		}
	}
	return regen - w.CalcEnergyUpkeep()
}

func (w *worldState) AddEnergy(delta float64) {
	if delta == 0 {
		return
	}
	w.energy = gmath.ClampMin(w.energy+delta, 0)
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

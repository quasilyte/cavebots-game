package battle

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/cavebots-game/controls"
	"github.com/quasilyte/cavebots-game/session"
	"github.com/quasilyte/cavebots-game/styles"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/pathing"
)

type Runner struct {
	scene *ge.Scene

	state *session.State

	world *worldState

	energyRegenDelay float64

	labelsRect  *ge.Rect
	energyLabel *ge.Label
	ironLabel   *ge.Label
	stonesLabel *ge.Label

	computer *computerPlayer

	startTime time.Time

	stillTime      float64
	hoverTriggered bool
	hoverPos       gmath.Vec
	ttm            *tooltipManager
	tutorial       *tutorialManager

	cameraPanSpeed    float64
	cameraPanBoundary float64

	cellSelector *ge.Sprite

	EventBattleCompleted gsignal.Event[*Results]
}

type Results struct {
	Victory bool

	BotsCreated int
	Digs        int

	Duration time.Duration

	Difficulty int
}

func NewRunner(state *session.State) *Runner {
	return &Runner{
		state:             state,
		cameraPanSpeed:    8,
		cameraPanBoundary: 8,
	}
}

func (r *Runner) Init(scene *ge.Scene) {
	r.scene = scene

	r.startTime = time.Now()

	r.world = &worldState{
		tutorial:   r.state.Settings.FirstTime,
		width:      1920,
		height:     32 * numCaveVerticalCells,
		caveWidth:  float64(32 * numCaveHorizontalCells),
		scene:      scene,
		rand:       scene.Rand(),
		difficulty: r.state.Settings.Difficulty,
	}
	r.world.Init()

	r.world.EventTooltipUpdateRequest.Connect(nil, func(gsignal.Void) {
		r.ttm.ForceUpdate()
	})

	r.world.EventDefeat.Connect(nil, func(pos gmath.Vec) {
		r.scene.AddObject(newBigFloatingTextNode(r.world, pos, "Defeat!"))
		r.world.results.Victory = false
		r.world.results.Duration = time.Since(r.startTime)
		r.world.results.Difficulty = r.world.difficulty
		r.EventBattleCompleted.Emit(r.world.results)
	})

	r.world.EventVictory.Connect(nil, func(pos gmath.Vec) {
		r.scene.AddObject(newBigFloatingTextNode(r.world, pos, "Victory!"))
		r.world.results.Victory = true
		r.world.results.Duration = time.Since(r.startTime)
		r.world.results.Difficulty = r.world.difficulty
		r.EventBattleCompleted.Emit(r.world.results)
	})

	r.computer = newComputerPlayer(r.world)

	r.world.EventResourcesUpdated.Connect(nil, func(gsignal.Void) {
		r.updateLabels()
	})

	spawnPos := gmath.Vec{
		X: float64((numCaveHorizontalCells - 1) * 32),
		Y: float64(32 * r.scene.Rand().IntRange(8, numCaveVerticalCells-8)),
	}

	r.initMap(spawnPos)

	r.cellSelector = scene.NewSprite(assets.ImageCellSelector)
	r.cellSelector.Visible = false
	r.world.stage.AddSpriteAbove(r.cellSelector)

	r.world.core = r.world.NewUnitNode(spawnPos.Add(gmath.Vec{X: 16, Y: 16}).Sub(gmath.Vec{X: 32}), droneCoreStats)
	scene.AddObject(r.world.core)

	{
		res := r.world.NewResourceNode(r.world.core.pos.Sub(gmath.Vec{X: 64}), ironResourceStats, 4)
		scene.AddObject(res)
	}

	r.world.camera.CenterOn(spawnPos)

	r.placeCreeps()

	r.world.diggedRect = gmath.Rect{
		Min: r.world.core.pos.Sub(gmath.Vec{X: 15, Y: 15}),
		Max: r.world.core.pos.Add(gmath.Vec{X: 15, Y: 15}),
	}

	r.labelsRect = ge.NewRect(scene.Context(), 1920.0/2, 56)
	r.labelsRect.Centered = false
	r.labelsRect.FillColorScale.SetColor(styles.BgSuperDark)
	r.labelsRect.Pos.Offset.Y = (1080 / 2) - 56
	r.world.camera.UI.AddGraphicsAbove(r.labelsRect)

	r.energyLabel = scene.NewLabel(assets.FontNormal)
	r.energyLabel.ColorScale.SetColor(styles.ButtonTextColor)
	r.energyLabel.Pos.Offset.Y = (1080 / 2) - 56
	r.energyLabel.Pos.Offset.X = 16
	r.energyLabel.Height = 56
	r.energyLabel.AlignVertical = ge.AlignVerticalCenter
	r.world.camera.UI.AddGraphicsAbove(r.energyLabel)

	r.ironLabel = scene.NewLabel(assets.FontNormal)
	r.ironLabel.ColorScale.SetColor(styles.ButtonTextColor)
	r.ironLabel.Pos.Offset.Y = (1080 / 2) - 56
	r.ironLabel.Pos.Offset.X = 16 + 320
	r.ironLabel.Height = 56
	r.ironLabel.AlignVertical = ge.AlignVerticalCenter
	r.world.camera.UI.AddGraphicsAbove(r.ironLabel)

	r.stonesLabel = scene.NewLabel(assets.FontNormal)
	r.stonesLabel.ColorScale.SetColor(styles.ButtonTextColor)
	r.stonesLabel.Pos.Offset.Y = (1080 / 2) - 56
	r.stonesLabel.Pos.Offset.X = 16 + 320 + 320
	r.stonesLabel.Height = 56
	r.stonesLabel.AlignVertical = ge.AlignVerticalCenter
	r.world.camera.UI.AddGraphicsAbove(r.stonesLabel)

	r.energyRegenDelay = 10

	r.updateLabels()

	r.ttm = newTooltipManager(r.world)
	scene.AddObject(r.ttm)

	scene.AddGraphics(r.world.camera)

	if r.world.tutorial {
		r.tutorial = newTutorialManager(r.world)
		r.tutorial.Init(scene)
	}
}

func (r *Runner) placeCreeps() {
	scene := r.scene

	for i := 0; i < 3; i++ {
		pos := r.world.core.pos.Add(gmath.Vec{X: 32 * 13})
		creep := r.world.NewUnitNode(pos.Add(scene.Rand().Offset(-12, 12)), creepMutantWarrior)
		scene.AddObject(creep)
	}

	for i := 0; i < 1; i++ {
		pos := r.world.core.pos.Add(gmath.Vec{X: 32 * 18, Y: -128})
		creep := r.world.NewUnitNode(pos.Add(scene.Rand().Offset(-12, 12)), creepMutantWarlord)
		scene.AddObject(creep)
	}
	if r.world.difficulty >= 1 {
		for i := 0; i < 2; i++ {
			pos := r.world.core.pos.Add(gmath.Vec{X: 32 * 17, Y: 32})
			creep := r.world.NewUnitNode(pos.Add(scene.Rand().Offset(-12, 12)), creepMutantHunter)
			scene.AddObject(creep)
		}
	}
	if r.world.difficulty >= 2 {
		for i := 0; i < 1; i++ {
			pos := r.world.core.pos.Add(gmath.Vec{X: 32 * 18, Y: (-32 * 5)})
			creep := r.world.NewUnitNode(pos.Add(scene.Rand().Offset(-12, 12)), creepMutantGunner)
			scene.AddObject(creep)
		}
	}

	{
		pos := r.world.core.pos.Add(gmath.Vec{X: 32 * 17, Y: 96})
		boss := r.world.NewUnitNode(pos, creepJeep)
		scene.AddObject(boss)
	}
	{
		pos := r.world.core.pos.Add(gmath.Vec{X: 32 * 15, Y: -32 * 6})
		boss := r.world.NewUnitNode(pos, creepJeep)
		scene.AddObject(boss)
	}

	{
		creep := r.world.NewUnitNode(r.world.core.pos.Add(gmath.Vec{X: 32 * 18, Y: -96}), creepMutantBase)
		scene.AddObject(creep)
		r.world.creepBase = creep
	}
}

func (r *Runner) initMap(spawnPos gmath.Vec) {
	r.initBackground(spawnPos)

	initialTunnel := make([]gmath.Vec, 0, 8)

	for i := 0; i < 8; i++ {
		initialTunnel = append(initialTunnel, gmath.Vec{
			X: float64((numCaveHorizontalCells-1-i)*32) + 16,
			Y: spawnPos.Y + 16,
		})
	}

	for y := 0; y < numCaveVerticalCells; y++ {
		for x := 0; x < numCaveHorizontalCells; x++ {
			pos := gmath.Vec{X: float64(x*32) + 16, Y: float64(y*32) + 16}
			coord := pathing.GridCoord{X: x, Y: y}
			if xslices.Contains(initialTunnel, pos) {
				r.world.grid.SetCellTile(coord, tileCaveMud)
				continue
			}
			m := newMountainNode(r.world, pos)
			if !r.world.innerCaveRect.Contains(pos) {
				m.outer = true
			}
			r.scene.AddObject(m)
			r.world.grid.SetCellTile(coord, tileBlocked)
			packedCoord := r.world.grid.PackCoord(coord)
			r.world.mountainByCoord[packedCoord] = m
		}
	}

	randIterate(r.world.rand, initialTunnel, func(pos gmath.Vec) bool {
		r.world.RevealNeighbors(pos)
		return false
	})
}

func (r *Runner) initBackground(spawnPos gmath.Vec) {
	wholeBg := ebiten.NewImage(1920, numCaveVerticalCells*32)
	{
		bg := ge.NewTiledBackground(r.scene.Context())
		bg.LoadTileset(r.scene.Context(), r.world.caveWidth, 32*numCaveVerticalCells, assets.ImageCaveTiles, assets.RawCaveTileset)
		bg.Draw(wholeBg)
	}
	{
		bg := ge.NewTiledBackground(r.scene.Context())
		bg.LoadTileset(r.scene.Context(), 1920-r.world.caveWidth, 32*numCaveVerticalCells, assets.ImageForestTiles, assets.RawCaveTileset)
		bg.Pos.Offset.X = r.world.caveWidth
		bg.Draw(wholeBg)
	}
	{
		s := r.scene.NewSprite(assets.ImageBiomeTransition)
		s.Pos.Offset = spawnPos
		s.Centered = false
		s.Draw(wholeBg)
	}
	s := ge.NewSprite(r.scene.Context())
	s.Centered = false
	s.SetImage(resource.Image{Data: wholeBg})
	r.world.stage.AddSpriteBelow(s)
}

func (r *Runner) Update(delta float64) {
	r.handleInput(delta)
	r.computer.Update(delta)
	r.world.stage.Update()
	if r.tutorial != nil {
		r.tutorial.Update(delta)
	}

	r.energyRegenDelay -= delta
	if r.energyRegenDelay <= 0 {
		r.energyRegenDelay = 10
		r.world.AddEnergy(r.world.CalcEnergyRegen())
	}
}

func (r *Runner) updateLabels() {
	energyIncome := r.world.CalcEnergyRegen()
	r.energyLabel.Text = fmt.Sprintf("Energy: %d (%+.1f)", int(r.world.energy), energyIncome)
	r.ironLabel.Text = fmt.Sprintf("Iron: %d", r.world.iron)
	r.stonesLabel.Text = fmt.Sprintf("Stone: %d", r.world.stones)
}

func (r *Runner) handleInput(delta float64) {
	cursorPos := r.state.Input.CursorPos()
	cursorWorldPos := r.world.camera.AbsPos(cursorPos)

	r.handleHover(cursorWorldPos, delta)

	handler := r.state.Input

	if r.world.rect.Contains(cursorWorldPos) {
		r.cellSelector.Visible = true
		r.cellSelector.Pos.Offset = makePos(r.world.grid.AlignPos(cursorWorldPos.X, cursorWorldPos.Y))
	} else {
		r.cellSelector.Visible = false
	}

	if r.world.core != nil && handler.ActionIsJustPressed(controls.ActionSendUnit) {
		playGlobalSound(r.world, assets.AudioUnitAck1)
		r.world.core.SendTo(cursorWorldPos)
		r.scene.AddObject(newFloatingTextNode(r.world, cursorWorldPos, "Order: move here"))
		return
	}

	if r.world.core != nil && handler.ActionIsJustPressed(controls.ActionInteract) {
		m := r.world.MountainAt(cursorWorldPos)
		if m != nil {
			if !r.world.CanDig(m) {
				playGlobalSound(r.world, assets.AudioError)
				r.scene.AddObject(newFloatingTextNode(r.world, cursorWorldPos, "Error: can't dig here"))
				return
			}
			if r.world.energy < digEnergyCost && m.loot != lootEasyDig {
				playGlobalSound(r.world, assets.AudioError)
				r.scene.AddObject(newFloatingTextNode(r.world, cursorWorldPos, "Error: not enough energy"))
				r.world.notEnoughEnergy += 1
				return
			}
			r.scene.AddObject(newFloatingTextNode(r.world, m.pos, "Order: dig here"))
			playGlobalSound(r.world, assets.AudioUnitAck1)
			r.world.core.SendDigging(cursorWorldPos)
			r.world.core.orderTarget = m
			return
		}
	}

	if handler.ActionIsJustPressed(controls.ActionBuild1) {
		r.doBuildAction(cursorWorldPos, 0)
		return
	}
	if handler.ActionIsJustPressed(controls.ActionBuild2) {
		r.doBuildAction(cursorWorldPos, 1)
		return
	}
	if handler.ActionIsJustPressed(controls.ActionBuild3) {
		r.doBuildAction(cursorWorldPos, 2)
		return
	}
	if handler.ActionIsJustPressed(controls.ActionBuild4) {
		r.doBuildAction(cursorWorldPos, 3)
		return
	}

	var cameraPan gmath.Vec
	if handler.ActionIsPressed(controls.ActionPanRight) {
		cameraPan.X += r.cameraPanSpeed
	}
	if handler.ActionIsPressed(controls.ActionPanDown) {
		cameraPan.Y += r.cameraPanSpeed
	}
	if handler.ActionIsPressed(controls.ActionPanLeft) {
		cameraPan.X -= r.cameraPanSpeed
	}
	if handler.ActionIsPressed(controls.ActionPanUp) {
		cameraPan.Y -= r.cameraPanSpeed
	}
	if cameraPan.IsZero() && r.cameraPanBoundary != 0 {
		// Mouse cursor can pan the camera too.
		cursor := handler.CursorPos()
		if cursor.X > r.world.camera.Rect.Width()-r.cameraPanBoundary {
			cameraPan.X += r.cameraPanSpeed
		}
		if cursor.Y > r.world.camera.Rect.Height()-r.cameraPanBoundary {
			cameraPan.Y += r.cameraPanSpeed
		}
		if cursor.X < r.cameraPanBoundary {
			cameraPan.X -= r.cameraPanSpeed
		}
		if cursor.Y < r.cameraPanBoundary {
			cameraPan.Y -= r.cameraPanSpeed
		}
	}
	if !cameraPan.IsZero() {
		r.world.camera.Pan(cameraPan)
	}
}

func (r *Runner) doBuildBuilding(coord pathing.GridCoord, cursorPos gmath.Vec, i int) {
	if i > 1 {
		return
	}

	var buildingSpot *hardTerrainNode
	for _, tile := range r.world.hardTerrain {
		if tile.building != nil {
			continue
		}
		tileCoord := r.world.grid.PosToCoord(tile.pos.X, tile.pos.Y)
		if tileCoord == coord {
			buildingSpot = tile
			break
		}
	}
	if buildingSpot == nil {
		return
	}

	buildingStats := buildingSpot.buildOptions[i]
	if !r.world.TryBuy(buildingStats, cursorPos) {
		return
	}

	newBuilding := r.world.NewUnitNode(buildingSpot.pos, buildingStats)
	r.scene.AddObject(newBuilding)
	buildingSpot.building = newBuilding
	newBuilding.EventDisposed.Connect(nil, func(*unitNode) {
		buildingSpot.building = nil
	})
	r.ttm.ForceUpdate()
	playGlobalSound(r.world, assets.AudioBuildingPlaced)
}

func (r *Runner) doBuildAction(cursorPos gmath.Vec, i int) {
	coord := r.world.grid.PosToCoord(cursorPos.X, cursorPos.Y)
	if coord.X < 0 || coord.Y < 0 {
		return
	}

	if r.world.grid.GetCellTile(coord) == tileCaveFlat {
		r.doBuildBuilding(coord, cursorPos, i)
		return
	}

	var factory *unitNode
	for _, u := range r.world.playerUnits {
		if u.stats != buildingFactory {
			continue
		}
		if u.order != orderNone {
			continue
		}
		factoryCoord := r.world.grid.PosToCoord(u.pos.X, u.pos.Y)
		if factoryCoord == coord {
			factory = u
			break
		}
	}
	if factory == nil {
		return
	}

	var stats *unitStats
	switch i {
	case 0:
		stats = droneHarvesterStats
	case 1:
		stats = dronePatrolStats
	case 2:
		stats = droneVanguardStats
	case 3:
		stats = droneTitanStats
	}
	if !r.world.TryBuy(stats, factory.pos) {
		return
	}

	factory.orderTarget = stats
	factory.order = orderMakeUnit
	factory.reload = r.scene.Rand().FloatRange(15, 20)
	r.ttm.ForceUpdate()
	playGlobalSound(r.world, assets.AudioProductionStarted)
}

func (r *Runner) stopHover() {
	r.ttm.OnStopHover()
}

func (r *Runner) hover(pos gmath.Vec) {
	r.ttm.OnHover(pos)
}

func (r *Runner) handleHover(pos gmath.Vec, delta float64) {
	maxDistSqr := 6.0 * 6.0
	if r.hoverPos.IsZero() {
		r.hoverPos = pos
	}
	distSqr := pos.DistanceSquaredTo(r.hoverPos)
	if distSqr < maxDistSqr {
		if !r.hoverTriggered {
			r.stillTime += delta
			if r.hoverPos.IsZero() && r.stillTime > 0.1 {
				r.hoverPos = pos
			}
			if r.stillTime > 0.2 {
				r.hoverTriggered = true
				r.hover(r.hoverPos)
			}
		}
	} else {
		if r.hoverTriggered && r.stillTime > 0 {
			r.hoverTriggered = false
			r.stopHover()
		}
		r.stillTime = 0
		r.hoverPos = gmath.Vec{}
	}
}

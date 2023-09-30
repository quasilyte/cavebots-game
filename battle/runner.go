package battle

import (
	"fmt"

	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/cavebots-game/controls"
	"github.com/quasilyte/cavebots-game/session"
	"github.com/quasilyte/cavebots-game/styles"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/pathing"
)

type Runner struct {
	scene *ge.Scene

	state *session.State

	core *unitNode

	world *worldState

	energyRegenDelay float64

	energyLabel *ge.Label
	ironLabel   *ge.Label
	stonesLabel *ge.Label

	stillTime      float64
	hoverTriggered bool
	hoverPos       gmath.Vec
	ttm            *tooltipManager

	cellSelector *ge.Sprite
}

func NewRunner(state *session.State) *Runner {
	return &Runner{state: state}
}

func (r *Runner) Init(scene *ge.Scene) {
	r.scene = scene

	r.world = &worldState{
		width:     1920,
		height:    32 * numCaveVerticalCells,
		caveWidth: float64(32 * numCaveHorizontalCells),
		scene:     scene,
		rand:      scene.Rand(),
	}
	r.world.Init()

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
	r.scene.AddGraphics(r.cellSelector)

	r.core = newUnitNode(r.world, droneCoreStats)
	r.core.pos = spawnPos.Add(gmath.Vec{X: 16, Y: 16})
	r.world.playerUnits = append(r.world.playerUnits, r.core)
	scene.AddObject(r.core)

	r.energyLabel = scene.NewLabel(assets.FontNormal)
	r.energyLabel.ColorScale.SetColor(styles.ButtonTextColor)
	r.energyLabel.Pos.Offset.Y = r.world.height
	r.energyLabel.Pos.Offset.X = 16
	r.energyLabel.Height = 1080 - r.world.height
	r.energyLabel.AlignVertical = ge.AlignVerticalCenter
	scene.AddGraphics(r.energyLabel)

	r.ironLabel = scene.NewLabel(assets.FontNormal)
	r.ironLabel.ColorScale.SetColor(styles.ButtonTextColor)
	r.ironLabel.Pos.Offset.Y = r.world.height
	r.ironLabel.Pos.Offset.X = 16 + 320
	r.ironLabel.Height = 1080 - r.world.height
	r.ironLabel.AlignVertical = ge.AlignVerticalCenter
	scene.AddGraphics(r.ironLabel)

	r.stonesLabel = scene.NewLabel(assets.FontNormal)
	r.stonesLabel.ColorScale.SetColor(styles.ButtonTextColor)
	r.stonesLabel.Pos.Offset.Y = r.world.height
	r.stonesLabel.Pos.Offset.X = 16 + 320 + 320
	r.stonesLabel.Height = 1080 - r.world.height
	r.stonesLabel.AlignVertical = ge.AlignVerticalCenter
	scene.AddGraphics(r.stonesLabel)

	r.energyRegenDelay = 10

	r.updateLabels()

	r.ttm = newTooltipManager(r.world)
	scene.AddObject(r.ttm)
}

func (r *Runner) initMap(spawnPos gmath.Vec) {
	caveWidth := r.world.caveWidth

	{
		bg := ge.NewTiledBackground(r.scene.Context())
		bg.LoadTileset(r.scene.Context(), caveWidth, 32*numCaveVerticalCells, assets.ImageCaveTiles, assets.RawCaveTileset)
		r.scene.AddGraphicsBelow(bg, 1)
	}
	{
		bg := ge.NewTiledBackground(r.scene.Context())
		bg.LoadTileset(r.scene.Context(), 1920-caveWidth, 32*numCaveVerticalCells, assets.ImageForestTiles, assets.RawCaveTileset)
		bg.Pos.Offset.X = caveWidth
		r.scene.AddGraphicsBelow(bg, 1)
	}

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
			if xslices.Contains(initialTunnel, pos) {
				continue
			}
			m := newMountainNode(pos)
			r.scene.AddObject(m)
			coord := pathing.GridCoord{X: x, Y: y}
			r.world.grid.SetCellTile(coord, tileBlocked)
			packedCoord := r.world.grid.PackCoord(coord)
			r.world.mountainByCoord[packedCoord] = m
		}
	}
}

func (r *Runner) Update(delta float64) {
	r.handleInput(delta)

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

	r.handleHover(cursorPos, delta)

	if r.world.rect.Contains(cursorPos) {
		r.cellSelector.Visible = true
		r.cellSelector.Pos.Offset = gmath.Vec{
			X: float64((int(cursorPos.X)/32)*32) + 16,
			Y: float64((int(cursorPos.Y)/32)*32) + 16,
		}
	} else {
		r.cellSelector.Visible = false
	}

	if info, ok := r.state.Input.JustPressedActionInfo(controls.ActionSendUnit); ok {
		r.core.SendTo(info.Pos)
		r.scene.AddObject(newFloatingTextNode(info.Pos, "Order: move here"))
		return
	}

	if info, ok := r.state.Input.JustPressedActionInfo(controls.ActionInteract); ok {
		m := r.world.MountainAt(info.Pos)
		if m != nil {
			if !r.world.CanDig(m) {
				r.scene.AddObject(newFloatingTextNode(info.Pos, "Error: can't dig here"))
				return
			}
			if r.world.energy < digEnergyCost {
				r.scene.AddObject(newFloatingTextNode(info.Pos, "Error: not enough energy"))
				return
			}
			r.scene.AddObject(newFloatingTextNode(m.pos, "Order: dig here"))
			r.core.SendDigging(info.Pos)
			r.core.orderTarget = m
			return
		}
	}
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
			if r.hoverPos.IsZero() && r.stillTime > 0.15 {
				r.hoverPos = pos
			}
			if r.stillTime > 0.3 {
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

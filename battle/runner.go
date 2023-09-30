package battle

import (
	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/cavebots-game/controls"
	"github.com/quasilyte/cavebots-game/session"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/pathing"
)

type Runner struct {
	scene *ge.Scene

	state *session.State

	core *unitNode

	world *worldState

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
	scene.AddObject(r.core)
}

func (r *Runner) initMap(spawnPos gmath.Vec) {
	caveWidth := r.world.caveWidth

	{
		bg := ge.NewTiledBackground(r.scene.Context())
		bg.LoadTileset(r.scene.Context(), caveWidth, 32*numCaveVerticalCells, assets.ImageCaveTiles, assets.RawCaveTileset)
		r.scene.AddGraphics(bg)
	}
	{
		bg := ge.NewTiledBackground(r.scene.Context())
		bg.LoadTileset(r.scene.Context(), 1920-caveWidth, 32*numCaveVerticalCells, assets.ImageForestTiles, assets.RawCaveTileset)
		bg.Pos.Offset.X = caveWidth
		r.scene.AddGraphics(bg)
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
	r.handleInput()
}

func (r *Runner) handleInput() {
	cursorPos := r.state.Input.CursorPos()
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
		return
	}

	if info, ok := r.state.Input.JustPressedActionInfo(controls.ActionInteract); ok {
		m := r.world.MountainAt(info.Pos)
		if m != nil && r.world.CanDig(m) {
			// TODO: could be a plain tile.
			r.world.grid.SetCellTile(r.world.grid.PosToCoord(m.pos.X, m.pos.Y), tileCaveMud)
			m.Dispose()
			return
		}
	}
}

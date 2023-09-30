package battle

import (
	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/cavebots-game/session"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
)

type Runner struct {
	scene *ge.Scene

	state *session.State

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
	}
	r.world.Init()

	r.initMap()

	r.cellSelector = scene.NewSprite(assets.ImageCellSelector)
	r.cellSelector.Visible = false
	r.scene.AddGraphics(r.cellSelector)
}

func (r *Runner) initMap() {
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

	caveEntranceY := r.scene.Rand().IntRange(8, numCaveVerticalCells-8)

	initialTunnel := make([]gmath.Vec, 0, 8)

	for i := 0; i < 8; i++ {
		initialTunnel = append(initialTunnel, gmath.Vec{
			X: float64((numCaveHorizontalCells - 1 - i) * 32),
			Y: float64((caveEntranceY - 1) * 32),
		})
	}

	for y := 0; y < numCaveVerticalCells; y++ {
		for x := 0; x < numCaveHorizontalCells; x++ {
			pos := gmath.Vec{X: float64(x * 32), Y: float64(y * 32)}
			if xslices.Contains(initialTunnel, pos) {
				continue
			}
			pos = pos.Sub(gmath.Vec{X: 8, Y: 8})
			m := newMountainNode(pos)
			r.scene.AddObject(m)
		}
	}
}

func (r *Runner) Update(delta float64) {
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
}

package scenes

import (
	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/cavebots-game/session"
	"github.com/quasilyte/ge"
)

type BattleController struct {
	scene *ge.Scene
	state *session.State
}

func NewBattleController(state *session.State) *BattleController {
	return &BattleController{state: state}
}

func (c *BattleController) Init(scene *ge.Scene) {
	c.scene = scene

	caveWidth := float64(32 * 40)

	{
		bg := ge.NewTiledBackground(scene.Context())
		bg.LoadTileset(scene.Context(), caveWidth, 32*32, assets.ImageCaveTiles, assets.RawCaveTileset)
		scene.AddGraphics(bg)
	}
	{
		bg := ge.NewTiledBackground(scene.Context())
		bg.LoadTileset(scene.Context(), 1920-caveWidth, 32*32, assets.ImageForestTiles, assets.RawCaveTileset)
		bg.Pos.Offset.X = caveWidth
		scene.AddGraphics(bg)
	}
}

func (c *BattleController) Update(delta float64) {}

package scenes

import (
	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/cavebots-game/session"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
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

	for y := 32.0; y < 32*5; y += 32 {
		for x := 32.0; x < 32*6; x += 32 {
			s := scene.NewSprite(assets.ImageMountains)
			s.Centered = false
			s.Pos.Offset = gmath.Vec{X: x, Y: y}
			s.FrameOffset.X = float64(scene.Rand().IntRange(0, 1)) * s.FrameWidth
			s.FlipHorizontal = scene.Rand().Bool()
			scene.AddGraphics(s)
		}
	}

	for y := 32.0; y < 32*5; y += 32 {
		for x := float64(32 * 7); x < 32*10; x += 32 {
			s := scene.NewSprite(assets.ImageMountains)
			s.Pos.Offset = gmath.Vec{X: x, Y: y}
			s.Centered = false
			s.FrameOffset.X = float64(scene.Rand().IntRange(0, 1)) * s.FrameWidth
			s.FlipHorizontal = scene.Rand().Bool()
			scene.AddGraphics(s)
		}
	}
}

func (c *BattleController) Update(delta float64) {}

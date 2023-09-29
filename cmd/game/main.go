package main

import (
	"time"

	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/cavebots-game/scenes"
	"github.com/quasilyte/ge"
)

func main() {
	ctx := ge.NewContext(ge.ContextConfig{})
	ctx.Rand.SetSeed(time.Now().Unix())
	ctx.GameName = "cavebots"
	ctx.WindowTitle = "Cavebots"
	ctx.WindowWidth = 1920
	ctx.WindowHeight = 1080
	ctx.FullScreen = true

	ctx.Loader.OpenAssetFunc = assets.MakeOpenAssetFunc(ctx)
	assets.RegisterResources(ctx)

	if err := ge.RunGame(ctx, scenes.NewMainMenuController()); err != nil {
		panic(err)
	}
}

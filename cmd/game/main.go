package main

import (
	"time"

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

	if err := ge.RunGame(ctx, scenes.NewMainMenuController()); err != nil {
		panic(err)
	}
}

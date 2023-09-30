package main

import (
	"time"

	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/cavebots-game/controls"
	"github.com/quasilyte/cavebots-game/eui"
	"github.com/quasilyte/cavebots-game/scenes"
	"github.com/quasilyte/cavebots-game/session"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
)

func main() {
	ctx := ge.NewContext(ge.ContextConfig{})
	ctx.Rand.SetSeed(time.Now().Unix())
	ctx.GameName = "cavebots"
	ctx.WindowTitle = "Cavebots"
	ctx.WindowWidth = 1920 / 2
	ctx.WindowHeight = 1080 / 2
	ctx.FullScreen = true

	ctx.Loader.OpenAssetFunc = assets.MakeOpenAssetFunc(ctx)
	assets.RegisterResources(ctx)

	state := &session.State{
		UIResources: eui.PrepareResources(ctx.Loader),
	}

	keymap := input.Keymap{
		controls.ActionSendUnit: {input.KeyMouseRight},
		controls.ActionInteract: {input.KeyMouseLeft},
		controls.ActionBuild1:   {input.KeyQ},
		controls.ActionBuild2:   {input.KeyW},
		controls.ActionPanRight: {input.KeyRight},
		controls.ActionPanDown:  {input.KeyDown},
		controls.ActionPanLeft:  {input.KeyLeft},
		controls.ActionPanUp:    {input.KeyUp},
	}
	state.Input = ctx.Input.NewHandler(0, keymap)

	if err := ge.RunGame(ctx, scenes.NewMainMenuController(state)); err != nil {
		panic(err)
	}
}

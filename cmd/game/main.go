package main

import (
	"fmt"
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
		Settings:    getDefaultSettings(),
	}

	keymap := input.Keymap{
		controls.ActionBack:     {input.KeyEscape},
		controls.ActionSendUnit: {input.KeyMouseRight},
		controls.ActionInteract: {input.KeyMouseLeft},
		controls.ActionBuild1:   {input.KeyQ},
		controls.ActionBuild2:   {input.KeyW},
		controls.ActionBuild3:   {input.KeyE},
		controls.ActionBuild4:   {input.KeyR},
		controls.ActionPanRight: {input.KeyRight},
		controls.ActionPanDown:  {input.KeyDown},
		controls.ActionPanLeft:  {input.KeyLeft},
		controls.ActionPanUp:    {input.KeyUp},
	}
	state.Input = ctx.Input.NewHandler(0, keymap)

	if err := ctx.LoadGameData("save", &state.Settings); err != nil {
		fmt.Printf("can't load game data: %v", err)
		state.Settings = getDefaultSettings()
		ctx.SaveGameData("save", state.Settings)
	}

	if err := ge.RunGame(ctx, scenes.NewMainMenuController(state)); err != nil {
		panic(err)
	}
}

func getDefaultSettings() session.Settings {
	return session.Settings{
		Difficulty: 0,
		SoundLevel: 2,
		MusicLevel: 2,
		FirstTime:  true,
	}
}

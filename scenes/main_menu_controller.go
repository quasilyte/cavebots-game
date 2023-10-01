package scenes

import (
	"fmt"
	"os"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/cavebots-game/eui"
	"github.com/quasilyte/cavebots-game/session"
	"github.com/quasilyte/cavebots-game/styles"
	"github.com/quasilyte/ge"
)

type MainMenuController struct {
	state *session.State
}

func NewMainMenuController(state *session.State) *MainMenuController {
	return &MainMenuController{state: state}
}

func (c *MainMenuController) Init(scene *ge.Scene) {
	scene.Audio().SetGroupVolume(assets.SoundGroupEffect, assets.VolumeMultiplier(c.state.Settings.SoundLevel))

	c.initUI(scene)
}

func (c *MainMenuController) initUI(scene *ge.Scene) {
	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	rowContainer := eui.NewRowLayoutContainerWithMinWidth(400, 16, nil)
	root.AddChild(rowContainer)

	bigFont := scene.Context().Loader.LoadFont(assets.FontBig).Face
	tinyFont := scene.Context().Loader.LoadFont(assets.FontTiny).Face

	rowContainer.AddChild(eui.NewCenteredLabel("CaveBots", bigFont))

	rowContainer.AddChild(eui.NewSeparator(nil, styles.TransparentColor))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, "PLAY", func() {
		scene.Context().ChangeScene(NewBattleController(c.state))
	}))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, "SETTINGS", func() {
		scene.Context().ChangeScene(NewSettingsController(c.state))
	}))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, "CREDITS", func() {
		scene.Context().ChangeScene(NewCreditsController(c.state))
	}))

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, "EXIT", func() {
		os.Exit(0)
	}))

	rowContainer.AddChild(eui.NewSeparator(nil, styles.TransparentColor))
	rowContainer.AddChild(eui.NewCenteredLabel(fmt.Sprintf("LD54 build %d", currentBuild), tinyFont))

	initUI(scene, root)
}

func (c *MainMenuController) Update(delta float64) {}

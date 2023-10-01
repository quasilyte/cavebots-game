package scenes

import (
	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/cavebots-game/controls"
	"github.com/quasilyte/cavebots-game/eui"
	"github.com/quasilyte/cavebots-game/session"
	"github.com/quasilyte/cavebots-game/styles"
	"github.com/quasilyte/ge"
)

type settingsController struct {
	state *session.State
	scene *ge.Scene
}

func NewSettingsController(state *session.State) *settingsController {
	return &settingsController{
		state: state,
	}
}

func (c *settingsController) Init(scene *ge.Scene) {
	c.scene = scene
	c.initUI(scene)
}

func (c *settingsController) initUI(scene *ge.Scene) {
	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	rowContainer := eui.NewRowLayoutContainerWithMinWidth(400, 16, nil)
	root.AddChild(rowContainer)

	bigFont := scene.Context().Loader.LoadFont(assets.FontBig).Face

	rowContainer.AddChild(eui.NewCenteredLabel("Credits", bigFont))
	rowContainer.AddChild(eui.NewSeparator(nil, styles.TransparentColor))

	rowContainer.AddChild(eui.NewSelectButton(eui.SelectButtonConfig{
		Resources:  c.state.UIResources,
		Input:      c.state.Input,
		Value:      &c.state.Settings.SoundLevel,
		Label:      "Volume level",
		ValueNames: []string{"0", "1", "2", "3", "4", "5"},
		OnPressed: func() {
			if c.state.Settings.SoundLevel != 0 {
				scene.Audio().SetGroupVolume(assets.SoundGroupEffect, assets.VolumeMultiplier(c.state.Settings.SoundLevel))
				scene.Audio().PlaySound(assets.AudioUnitAck2)
			}
		},
	}))

	rowContainer.AddChild(eui.NewSeparator(nil, styles.TransparentColor))
	rowContainer.AddChild(eui.NewButton(c.state.UIResources, "OK", func() {
		c.leave()
	}))

	initUI(scene, root)
}

func (c *settingsController) Update(delta float64) {
	if c.state.Input.ActionIsJustPressed(controls.ActionBack) {
		c.leave()
	}
}

func (c *settingsController) leave() {
	c.scene.Context().SaveGameData("save", c.state.Settings)
	c.scene.Context().ChangeScene(NewMainMenuController(c.state))
}

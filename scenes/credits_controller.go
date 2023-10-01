package scenes

import (
	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/cavebots-game/controls"
	"github.com/quasilyte/cavebots-game/eui"
	"github.com/quasilyte/cavebots-game/session"
	"github.com/quasilyte/ge"
)

type creditsController struct {
	state *session.State

	scene *ge.Scene
}

func NewCreditsController(state *session.State) *creditsController {
	return &creditsController{
		state: state,
	}
}

func (c *creditsController) Init(scene *ge.Scene) {
	c.scene = scene
	c.initUI(scene)
}

func (c *creditsController) initUI(scene *ge.Scene) {
	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	rowContainer := eui.NewRowLayoutContainerWithMinWidth(400, 16, nil)
	root.AddChild(rowContainer)

	bigFont := scene.Context().Loader.LoadFont(assets.FontBig).Face
	smallFont := scene.Context().Loader.LoadFont(assets.FontSmall).Face

	rowContainer.AddChild(eui.NewCenteredLabel("Credits", bigFont))

	text := `A game by quasilyte
Created during a compo LD54

Made with Ebitengine`

	panel := eui.NewPanelWithPadding(c.state.UIResources, 0, 0, widget.NewInsetsSimple(24))
	content := eui.NewLabel(text, smallFont)
	panel.AddChild(content)
	rowContainer.AddChild(panel)

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, "OK", func() {
		scene.Context().ChangeScene(NewMainMenuController(c.state))
	}))

	initUI(scene, root)
}

func (c *creditsController) Update(delta float64) {
	if c.state.Input.ActionIsJustPressed(controls.ActionBack) {
		c.scene.Context().ChangeScene(NewMainMenuController(c.state))
	}
}

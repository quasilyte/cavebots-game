package scenes

import (
	"strconv"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/cavebots-game/battle"
	"github.com/quasilyte/cavebots-game/eui"
	"github.com/quasilyte/cavebots-game/session"
	"github.com/quasilyte/ge"
)

type resultsController struct {
	state   *session.State
	results *battle.Results
}

func NewResultscontroller(state *session.State, results *battle.Results) *resultsController {
	return &resultsController{
		state:   state,
		results: results,
	}
}

func (c *resultsController) Init(scene *ge.Scene) {
	c.initUI(scene)
}

func (c *resultsController) initUI(scene *ge.Scene) {
	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	rowContainer := eui.NewRowLayoutContainerWithMinWidth(400, 16, nil)
	root.AddChild(rowContainer)

	bigFont := scene.Context().Loader.LoadFont(assets.FontBig).Face
	smallFont := scene.Context().Loader.LoadFont(assets.FontSmall).Face

	title := "Victory!"
	if !c.results.Victory {
		title = "Defeat!"
	}
	rowContainer.AddChild(eui.NewCenteredLabel(title, bigFont))

	panel := eui.NewPanelWithPadding(c.state.UIResources, 0, 0, widget.NewInsetsSimple(16))
	table := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Stretch(nil, nil),
			widget.GridLayoutOpts.Spacing(24, 8))))
	panel.AddChild(table)
	rowContainer.AddChild(panel)

	rows := [][2]string{
		{"Game duration", formatDurationCompact(c.results.Duration)},
		{"Bots created", strconv.Itoa(c.results.BotsCreated)},
		{"Tiles excavated", strconv.Itoa(c.results.Digs)},
	}
	for _, row := range rows {
		table.AddChild(eui.NewLabel(row[0], smallFont))
		table.AddChild(eui.NewLabel(row[1], smallFont))
	}

	rowContainer.AddChild(eui.NewButton(c.state.UIResources, "OK", func() {
		scene.Context().ChangeScene(NewMainMenuController(c.state))
	}))

	initUI(scene, root)
}

func (c *resultsController) Update(delta float64) {}

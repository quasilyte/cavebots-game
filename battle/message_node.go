package battle

import (
	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/cavebots-game/styles"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type messageNode struct {
	rect  *ge.Rect
	label *ge.Label

	world *worldState

	pos  gmath.Vec
	text string
}

func newMessageNode(world *worldState, pos gmath.Vec, text string) *messageNode {
	return &messageNode{
		world: world,
		text:  text,
		pos:   pos,
	}
}

func (m *messageNode) Init(scene *ge.Scene) {
	w, h := estimateMessageBounds(scene.Context().Loader.LoadFont(assets.FontTiny).Face, m.text, 8)

	m.label = scene.NewLabel(assets.FontTiny)
	m.label.Text = m.text
	m.label.Pos.Base = &m.pos
	m.label.Pos.Offset.X = 8
	m.label.Width = w
	m.label.Height = h
	m.label.AlignVertical = ge.AlignVerticalCenter
	m.label.ColorScale.SetColor(styles.ButtonTextColor)

	m.rect = ge.NewRect(scene.Context(), w+16, h)
	m.rect.Centered = false
	m.rect.FillColorScale.SetColor(styles.BgDark)
	m.rect.FillColorScale.A = 0.7
	m.rect.OutlineWidth = 1
	m.rect.OutlineColorScale.SetColor(styles.OutlineLight)
	m.rect.Pos.Base = &m.pos

	m.world.camera.UI.AddGraphicsAbove(m.rect)
	m.world.camera.UI.AddGraphicsAbove(m.label)
}

func (m *messageNode) Update(delta float64) {}

func (m *messageNode) Dispose() {
	m.rect.Dispose()
	m.label.Dispose()
}

func (m *messageNode) IsDisposed() bool {
	return m.rect.IsDisposed()
}

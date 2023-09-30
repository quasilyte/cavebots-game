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

	pos  gmath.Vec
	text string
}

func newMessageNode(pos gmath.Vec, text string) *messageNode {
	return &messageNode{
		pos:  pos,
		text: text,
	}
}

func (m *messageNode) Init(scene *ge.Scene) {
	w, h := estimateMessageBounds(scene.Context().Loader.LoadFont(assets.FontSmall).Face, m.text, 16)

	m.label = scene.NewLabel(assets.FontSmall)
	m.label.Text = m.text
	m.label.Pos.Base = &m.pos
	m.label.Pos.Offset.X = 8
	m.label.Width = w
	m.label.Height = h
	m.label.AlignHorizontal = ge.AlignHorizontalCenter
	m.label.AlignVertical = ge.AlignVerticalCenter
	m.label.ColorScale.SetColor(styles.ButtonTextColor)

	m.rect = ge.NewRect(scene.Context(), w+8, h)
	m.rect.Centered = false
	m.rect.FillColorScale.SetColor(styles.BgDark)
	m.rect.FillColorScale.A = 0.5
	m.rect.OutlineColorScale.SetColor(styles.OutlineLight)
	m.rect.Pos.Base = &m.pos

	scene.AddGraphics(m.rect)
	scene.AddGraphics(m.label)
}

func (m *messageNode) Update(delta float64) {}

func (m *messageNode) Dispose() {
	m.rect.Dispose()
	m.label.Dispose()
}

func (m *messageNode) IsDisposed() bool {
	return m.rect.IsDisposed()
}

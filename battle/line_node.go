package battle

import (
	"image/color"

	"github.com/quasilyte/ge"
)

type lineNode struct {
	from  ge.Pos
	to    ge.Pos
	world *worldState

	color      color.RGBA
	width      float64
	line       *ge.Line
	opaqueTime float64
}

func newLineNode(world *worldState, from, to ge.Pos, clr color.RGBA) *lineNode {
	return &lineNode{
		world: world,
		from:  from,
		to:    to,
		width: 1,
		color: clr,
	}
}

func (l *lineNode) Init(scene *ge.Scene) {
	l.line = ge.NewLine(l.from, l.to)
	var c ge.ColorScale
	c.SetColor(l.color)
	l.line.SetColorScale(c)
	l.line.Width = l.width
	l.world.stage.AddGraphicsAbove(l.line)
}

func (l *lineNode) IsDisposed() bool {
	return l.line.IsDisposed()
}

func (l *lineNode) Update(delta float64) {
	if l.opaqueTime > 0 {
		l.opaqueTime -= delta
		return
	}

	if l.line.GetAlpha() < 0.1 {
		l.line.Dispose()
		return
	}
	l.line.SetAlpha(l.line.GetAlpha() - float32(delta*4))
}

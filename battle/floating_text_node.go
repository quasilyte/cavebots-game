package battle

import (
	"math"

	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/cavebots-game/styles"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type floatingTextNode struct {
	label     *ge.Label
	rect      *ge.Rect
	world     *worldState
	pos       gmath.Vec
	spritePos gmath.Vec
	text      string
	hp        float64
}

func newFloatingTextNode(world *worldState, pos gmath.Vec, text string) *floatingTextNode {
	return &floatingTextNode{
		world: world,
		pos:   pos,
		text:  text,
		hp:    1.5,
	}
}

func (t *floatingTextNode) Init(scene *ge.Scene) {
	w, h := estimateMessageBounds(scene.Context().Loader.LoadFont(assets.FontTiny).Face, t.text, 16)

	t.pos = t.pos.Sub(t.world.camera.Offset)

	t.label = scene.NewLabel(assets.FontTiny)
	t.label.Text = t.text
	t.label.Pos.Base = &t.spritePos
	t.label.Pos.Offset.X = 4
	t.label.Width = w
	t.label.Height = h
	t.label.AlignHorizontal = ge.AlignHorizontalCenter
	t.label.AlignVertical = ge.AlignVerticalCenter
	t.label.ColorScale.SetColor(styles.ButtonTextColor)

	t.rect = ge.NewRect(scene.Context(), w+8, h)
	t.rect.Centered = false
	t.rect.FillColorScale.SetColor(styles.BgDark)
	t.rect.OutlineColorScale.SetColor(styles.OutlineLight)
	t.rect.Pos.Base = &t.spritePos

	// TODO: handle out-of-screen.
	t.pos = t.pos.Sub(gmath.Vec{
		X: w / 2,
		Y: h / 2,
	})
	t.pos.Y -= 54
	t.spritePos = t.pos

	t.world.camera.UI.AddGraphics(t.rect)
	t.world.camera.UI.AddGraphics(t.label)
}

func (t *floatingTextNode) IsDisposed() bool {
	return t.label.IsDisposed()
}

func (t *floatingTextNode) Update(delta float64) {
	t.hp -= delta
	if t.hp <= 0 {
		t.label.Dispose()
		t.rect.Dispose()
		return
	}

	t.pos.Y -= delta * 12
	t.spritePos.X = t.pos.X
	t.spritePos.Y = math.Round(t.pos.Y)
	t.label.ColorScale.A -= float32(delta / 2)
	t.rect.FillColorScale.A -= float32(delta / 2)
	t.rect.OutlineColorScale.A -= float32(delta / 2)
}

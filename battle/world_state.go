package battle

import "github.com/quasilyte/gmath"

type worldState struct {
	width  float64
	height float64

	rect gmath.Rect

	caveWidth float64
}

func (w *worldState) Init() {
	w.rect = gmath.Rect{
		Max: gmath.Vec{
			X: w.width,
			Y: w.height,
		},
	}
}

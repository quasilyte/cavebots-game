package battle

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/pathing"
)

type worldState struct {
	width  float64
	height float64

	scene *ge.Scene
	rand  *gmath.Rand

	rect gmath.Rect

	grid  *pathing.Grid
	astar *pathing.AStar

	caveWidth float64
}

func (w *worldState) Init() {
	w.rect = gmath.Rect{
		Max: gmath.Vec{
			X: w.width,
			Y: w.height,
		},
	}

	w.grid = pathing.NewGrid(pathing.GridConfig{
		WorldWidth:  uint(w.width),
		WorldHeight: uint(w.height),
		CellWidth:   32,
		CellHeight:  32,
	})
	w.astar = pathing.NewAStar(pathing.AStarConfig{
		NumCols: uint(w.grid.NumCols()),
		NumRows: uint(w.grid.NumRows()),
	})
}

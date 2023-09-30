package battle

import (
	"github.com/quasilyte/pathing"
)

const (
	tileCaveMud = iota
	tileCaveFlat
	tileGrass
	tileBlocked
)

var (
	normalLayer = pathing.MakeGridLayer([4]uint8{
		tileCaveMud:  1,
		tileCaveFlat: 1,
		tileGrass:    1,
		tileBlocked:  0,
	})
)

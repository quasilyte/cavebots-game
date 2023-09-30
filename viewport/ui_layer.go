package viewport

import (
	"github.com/quasilyte/ge"
)

type UserInterfaceLayer struct {
	belowObjects []ge.SceneGraphics
	objects      []ge.SceneGraphics
	aboveObjects []ge.SceneGraphics

	Visible bool
}

func (l *UserInterfaceLayer) AddGraphicsBelow(o ge.SceneGraphics) {
	l.belowObjects = append(l.belowObjects, o)
}

func (l *UserInterfaceLayer) AddGraphics(o ge.SceneGraphics) {
	l.objects = append(l.objects, o)
}

func (l *UserInterfaceLayer) AddGraphicsAbove(o ge.SceneGraphics) {
	l.aboveObjects = append(l.aboveObjects, o)
}

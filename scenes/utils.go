package scenes

import (
	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/cavebots-game/eui"
	"github.com/quasilyte/ge"
)

func initUI(scene *ge.Scene, root *widget.Container) {
	// TODO: add menu background image?

	uiObject := eui.NewSceneObject(root)
	scene.AddGraphics(uiObject)
	scene.AddObject(uiObject)
}

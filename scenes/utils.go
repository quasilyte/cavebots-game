package scenes

import (
	"fmt"
	"time"

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

func formatDurationCompact(d time.Duration) string {
	d = d.Round(time.Second)
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	d -= minutes * time.Minute
	seconds := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

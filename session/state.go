package session

import (
	"github.com/quasilyte/cavebots-game/eui"
	"github.com/quasilyte/ge/input"
)

type State struct {
	UIResources *eui.Resources

	Input *input.Handler
}

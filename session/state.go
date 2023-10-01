package session

import (
	"github.com/quasilyte/cavebots-game/eui"
	"github.com/quasilyte/ge/input"
)

type State struct {
	UIResources *eui.Resources

	Settings Settings

	Input *input.Handler
}

type Settings struct {
	SoundLevel int
	MusicLevel int
	Difficulty int
	FirstTime  bool
}

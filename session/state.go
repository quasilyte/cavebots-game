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
}

func VolumeMultiplier(level int) float64 {
	switch level {
	case 1:
		return 0.01
	case 2:
		return 0.05
	case 3:
		return 0.10
	case 4:
		return 0.3
	case 5:
		return 0.55
	case 6:
		return 0.8
	case 7:
		return 1.0
	default:
		return 0
	}
}

package controls

import (
	"github.com/quasilyte/ge/input"
)

const (
	ActionUnknown input.Action = iota

	ActionBack

	ActionSendUnit
	ActionInteract

	ActionBuild1
	ActionBuild2
	ActionBuild3

	ActionPanRight
	ActionPanDown
	ActionPanLeft
	ActionPanUp
)

package scenes

import (
	"github.com/quasilyte/cavebots-game/battle"
	"github.com/quasilyte/cavebots-game/session"
	"github.com/quasilyte/ge"
)

type BattleController struct {
	scene *ge.Scene
	state *session.State

	runner     *battle.Runner
	transition bool
}

func NewBattleController(state *session.State) *BattleController {
	return &BattleController{state: state}
}

func (c *BattleController) Init(scene *ge.Scene) {
	c.scene = scene

	c.runner = battle.NewRunner(c.state)
	c.runner.Init(scene)

	c.runner.EventBattleCompleted.Connect(nil, func(results *battle.Results) {
		c.transition = true
		scene.DelayedCall(4, func() {
			scene.Context().ChangeScene(NewResultsController(c.state, results))
		})
	})
}

func (c *BattleController) Update(delta float64) {
	c.runner.Update(delta)
}

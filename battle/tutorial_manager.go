package battle

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type tutorialManager struct {
	world *worldState
	scene *ge.Scene

	step int

	contextTimer float64

	timer   float64
	message *messageNode

	notEnoughEnergyHint bool
}

func newTutorialManager(world *worldState) *tutorialManager {
	return &tutorialManager{
		world: world,
	}
}

func (m *tutorialManager) Init(scene *ge.Scene) {
	m.scene = scene
}

func (m *tutorialManager) Update(delta float64) {
	m.world.notEnoughEnergy = gmath.ClampMin(m.world.notEnoughEnergy-delta, 0)

	if m.message == nil {
		if !m.notEnoughEnergyHint && m.maybeHintAboutEnergy() {
			m.notEnoughEnergyHint = true
			return
		}
	}

	if m.contextTimer != 0 {
		m.contextTimer = gmath.ClampMin(m.contextTimer-delta, 0)
		if m.contextTimer == 0 {
			m.message.Dispose()
			m.message = nil
		}
		return
	}

	m.timer += delta

	switch m.step {
	case 0:
		msg := `Welcome to CaveBots!

Hover over things to get contextual hints.`
		m.message = newMessageNode(m.world, gmath.Vec{X: 64, Y: 64}, msg)
		m.scene.AddObject(m.message)
		m.step++

	case 1:
		if m.timer >= 20 {
			m.finishMessageStep()
		}

	case 2:

	}
}

func (m *tutorialManager) finishMessageStep() {
	m.message.Dispose()
	m.message = nil
	m.step++

	m.world.notEnoughEnergy = 0
}

func (m *tutorialManager) maybeHintAboutEnergy() bool {
	if m.world.notEnoughEnergy < 2.1 {
		return false
	}

	msg := `To get more energy, build Generators.
Keep in mind that bots gradually consume the energy.`
	m.message = newMessageNode(m.world, gmath.Vec{X: 64, Y: 64}, msg)
	m.scene.AddObject(m.message)
	m.contextTimer = 15

	return true
}

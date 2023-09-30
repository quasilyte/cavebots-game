package battle

import (
	"fmt"

	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type tooltipManager struct {
	world   *worldState
	message *messageNode

	tooltipTime float64
}

func newTooltipManager(world *worldState) *tooltipManager {
	return &tooltipManager{
		world: world,
	}
}

func (m *tooltipManager) Init(scene *ge.Scene) {

}

func (m *tooltipManager) IsDisposed() bool {
	return false
}

func (m *tooltipManager) Update(delta float64) {
	if m.tooltipTime > 0 {
		m.tooltipTime = gmath.ClampMin(m.tooltipTime-delta, 0)
		if m.tooltipTime == 0 {
			m.OnStopHover()
		}
	}
}

func (m *tooltipManager) OnStopHover() {
	m.removeTooltip()
}

func (m *tooltipManager) OnHover(pos gmath.Vec) {
	mountain := m.world.MountainAt(pos)
	if mountain != nil {
		s := "Mountain"
		if m.world.CanDig(mountain) {
			s = m.formatMountainInfo(mountain)
		}
		m.createTooltip(pos, s)
		return
	}

	for _, res := range m.world.resourceNodes {
		if res.pos.DistanceSquaredTo(pos) < (22 * 22) {
			s := fmt.Sprintf("Iron resource (%d)", res.amount)
			m.createTooltip(pos, s)
			return
		}
	}
}

func (m *tooltipManager) formatMountainInfo(mountain *mountainNode) string {
	terrain := "normal"
	loot := ""
	switch m.world.PeekLoot(mountain) {
	case lootLavaCell:
		terrain = "lava"
	case lootFlatCell:
		terrain = "flat"

	case lootIronDeposit:
		loot = "iron deposit"
	case lootLargeIronDeposit:
		loot = "large iron deposit"

	case lootBotHarvester:
		loot = "harvester bot"
	case lootBotPatrol:
		loot = "patrol bot"
	case lootBotGuard:
		loot = "guard bot"
	}

	if loot == "" {
		return "Diggable mountain\nTerrain: " + terrain
	}
	return "Diggable mountain\nTerrain: " + terrain + "\nExtra: " + loot
}

func (m *tooltipManager) removeTooltip() {
	m.tooltipTime = 0
	if m.message != nil {
		m.message.Dispose()
		m.message = nil
	}
}

func (m *tooltipManager) createTooltip(pos gmath.Vec, s string) {
	if m.message != nil {
		m.removeTooltip()
	}

	m.tooltipTime = 5

	w, h := estimateMessageBounds(m.world.scene.Context().Loader.LoadFont(assets.FontSmall).Face, s, 0)
	if w+pos.X+26 > 1920 {
		pos.X -= w
	}
	if h+pos.Y+26 > 1080 {
		pos.Y -= h + 26
	} else {
		pos.Y += 26
	}

	m.message = newMessageNode(pos, s)
	m.world.scene.AddObject(m.message)
}

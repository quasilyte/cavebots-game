package battle

import (
	"fmt"
	"math"
	"strconv"
	"strings"

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
		s := "Inner block"
		if mountain.outer {
			s = "Outer block"
		} else if m.world.CanDig(mountain) {
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

	for _, u := range m.world.playerUnits {
		if u.pos.DistanceSquaredTo(pos) < (18 * 18) {
			s := "Core"
			health := strconv.Itoa(int(100*math.Ceil(u.health/u.stats.maxHealth))) + "%"
			if u.stats != droneCoreStats {
				if u.stats.building {
					s = u.stats.name
				} else {
					status := "online"
					if u.offline {
						status = "offline"
					}
					s = fmt.Sprintf("%s drone (%s)", u.stats.name, status)
				}
			}
			m.createTooltip(pos, s+"\n"+health)
			return
		}
	}

	for _, t := range m.world.hardTerrain {
		if t.building != nil {
			continue
		}
		if t.pos.DistanceSquaredTo(pos) < (20 * 20) {
			// TODO: show building options.
			parts := []string{"Hard terrain, can build here:"}
			for i, option := range t.buildOptions {
				priceParts := make([]string, 0, 3)
				if option.energyCost != 0 {
					priceParts = append(priceParts, strconv.Itoa(option.energyCost)+" energy")
				}
				if option.ironCost != 0 {
					priceParts = append(priceParts, strconv.Itoa(option.ironCost)+" iron")
				}
				if option.stoneCost != 0 {
					priceParts = append(priceParts, strconv.Itoa(option.stoneCost)+" stone")
				}
				price := strings.Join(priceParts, " and ")
				parts = append(parts, fmt.Sprintf("[%s] %s - %s", buildingHotkeys[i], option.name, price))
			}
			m.createTooltip(pos, strings.Join(parts, "\n"))
			return
		}
	}
}

func (m *tooltipManager) formatMountainInfo(mountain *mountainNode) string {
	terrain := "soft"
	loot := ""
	switch m.world.PeekLoot(mountain) {
	case lootLavaCell:
		terrain = "lava"
	case lootFlatCell:
		terrain = "hard"

	case lootIronDeposit:
		loot = "iron deposit"
	case lootLargeIronDeposit:
		loot = "large iron deposit"

	case lootExtraStones:
		loot = "stone-rich"

	case lootBotHarvester:
		loot = "harvester bot"
	case lootBotPatrol:
		loot = "patrol bot"
	case lootBotGuard:
		loot = "guard bot"
	}

	if loot == "" {
		return "Diggable block\nTerrain: " + terrain
	}
	return "Diggable block\nTerrain: " + terrain + "\nExtra: " + loot
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

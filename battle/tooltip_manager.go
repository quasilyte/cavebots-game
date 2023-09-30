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

func (m *tooltipManager) formatPrice(u *unitStats) string {
	priceParts := make([]string, 0, 3)
	if u.energyCost != 0 {
		priceParts = append(priceParts, strconv.Itoa(u.energyCost)+" energy")
	}
	if u.ironCost != 0 {
		priceParts = append(priceParts, strconv.Itoa(u.ironCost)+" iron")
	}
	if u.stoneCost != 0 {
		priceParts = append(priceParts, strconv.Itoa(u.stoneCost)+" stone")
	}
	return strings.Join(priceParts, " / ")
}

func (m *tooltipManager) OnHover(pos gmath.Vec) {
	screenPos := pos.Sub(m.world.camera.Offset)
	if screenPos.Y >= ((1080.0 / 2) - 56) {
		// TODO: resource hints.
		return
	}

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
			s := fmt.Sprintf("Iron resource (%d)\nHarvesters collect it", res.amount)
			m.createTooltip(pos, s)
			return
		}
	}

	for _, u := range m.world.playerUnits {
		if u.pos.DistanceSquaredTo(pos) < (18 * 18) {
			s := "Core"
			health := strconv.Itoa(int(math.Ceil(100*u.health/u.stats.maxHealth))) + "%"
			if u.stats != droneCoreStats {
				if u.stats.building {
					s = u.stats.name
				} else {
					status := "online"
					if u.offline {
						status = "offline"
					}
					s = fmt.Sprintf("%s bot (%s)", u.stats.name, status)
				}
			}
			if u.stats == buildingFactory {
				if u.order == orderMakeUnit {
					s += "\n" + u.orderTarget.(*unitStats).name + " bot is being produced..."
				} else {
					s += "\n[Q] Harvester - " + m.formatPrice(droneHarvesterStats)
					s += "\n[W] Patrol - " + m.formatPrice(dronePatrolStats)
					s += "\n[E] Vanguard - " + m.formatPrice(droneVanguardStats)
				}
			}
			m.createTooltip(pos, s+"\nHP: "+health)
			return
		}
	}

	for _, t := range m.world.hardTerrain {
		if t.building != nil {
			continue
		}
		if t.pos.DistanceSquaredTo(pos) < (20 * 20) {
			parts := []string{"Hard terrain, can build here:"}
			for i, option := range t.buildOptions {
				price := m.formatPrice(option)
				parts = append(parts, fmt.Sprintf("[%s] %s - %s", buildingHotkeys[i], option.name, price))
			}
			m.createTooltip(pos, strings.Join(parts, "\n"))
			return
		}
	}

	for _, u := range m.world.creeps {
		if u.pos.DistanceSquaredTo(pos) < (16 * 16) {
			m.createTooltip(pos, u.stats.name+" (hostile)")
			return
		}
	}

	if !m.world.rect.Contains(pos) {
		return
	}
	coord := m.world.grid.PosToCoord(pos.X, pos.Y)
	if coord.X < 0 || coord.Y < 0 {
		return
	}
	cellType := m.world.grid.GetCellTile(coord)
	switch cellType {
	case tileCaveMud:
		m.createTooltip(pos, "Cave area\n[RMB to move here]")
	case tileGrass:
		m.createTooltip(pos, "Forest area\n[RMB to move here]")
	}
}

func (m *tooltipManager) formatMountainInfo(mountain *mountainNode) string {
	var extra string
	switch m.world.PeekLoot(mountain) {
	case lootLavaCell:
		extra = "Lava terrain"
	case lootFlatCell:
		extra = "Hard terrain (can build)"

	case lootIronDeposit, lootLargeIronDeposit:
		extra = "Contains iron deposit"

	case lootExtraStones:
		extra = "Grants extra stones"
	case lootEasyDig:
		extra = "Can be dug for free"

	case lootBotHarvester:
		extra = "Contains harvester bot"
	case lootBotPatrol:
		extra = "Contains patrol bot"
	case lootBotVanguard:
		extra = "Contains vanguard bot"
	case lootBotGenerator:
		extra = "Contains generator bot"
	}

	if extra == "" {
		return "Diggable block\n[LMB to dig here]"
	}
	return "Diggable block\n" + extra + "\n[LMB to dig here]"
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
	screepPos := pos.Sub(m.world.camera.Offset)

	w, h := estimateMessageBounds(m.world.scene.Context().Loader.LoadFont(assets.FontSmall).Face, s, 0)
	if w+screepPos.X+26 > 1920.0/2 {
		screepPos.X -= w
	}
	if h+screepPos.Y+26 > 1080.0/2 {
		screepPos.Y -= h + 26
	} else {
		screepPos.Y += 26
	}

	m.message = newMessageNode(m.world, screepPos, s)
	m.world.scene.AddObject(m.message)
}

package battle

import (
	"github.com/quasilyte/cavebots-game/assets"
	resource "github.com/quasilyte/ebitengine-resource"
)

type unitStats struct {
	name         string
	speed        float64
	maxHealth    float64
	img          resource.ImageID
	allied       bool
	building     bool
	energyUpkeep float64

	energyCost int
	ironCost   int
	stoneCost  int
}

var droneCoreStats = &unitStats{
	name:         "Core",
	speed:        80,
	maxHealth:    100,
	img:          assets.ImageDroneCore,
	allied:       true,
	energyUpkeep: 0,
}

var droneHarvesterStats = &unitStats{
	name:         "Harvester",
	speed:        96,
	maxHealth:    30,
	img:          assets.ImageDroneHarvester,
	allied:       true,
	energyUpkeep: 0.2,
}

var dronePatrolStats = &unitStats{
	name:         "Patrol",
	speed:        110,
	maxHealth:    60,
	img:          assets.ImageDronePatrol,
	allied:       true,
	energyUpkeep: 0.3,
}

var buildingPowerGenerator = &unitStats{
	name:      "Generator",
	maxHealth: 75,
	img:       assets.ImageBuildingGenerator,
	allied:    true,
	building:  true,
	ironCost:  1,
	stoneCost: 6,
}

var buildingBarricate = &unitStats{
	name:       "Barricade",
	maxHealth:  90,
	img:        assets.ImageBuildingBarricade,
	allied:     true,
	building:   true,
	energyCost: 3,
	ironCost:   1,
}

var buildingSmelter = &unitStats{
	name:       "Smelter",
	maxHealth:  130,
	img:        assets.ImageBuildingSmelter,
	allied:     true,
	building:   true,
	energyCost: 3,
	ironCost:   2,
	stoneCost:  5,
}

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
	tiny         bool
	energyUpkeep float64
	layer        int

	energyCost int
	ironCost   int
	stoneCost  int
}

var droneCoreStats = &unitStats{
	name:         "Core",
	layer:        2,
	speed:        80,
	maxHealth:    100,
	img:          assets.ImageDroneCore,
	allied:       true,
	energyUpkeep: 0,
}

var droneHarvesterStats = &unitStats{
	name:         "Harvester",
	layer:        2,
	speed:        96,
	maxHealth:    30,
	img:          assets.ImageDroneHarvester,
	allied:       true,
	energyUpkeep: 0.2,
}

var dronePatrolStats = &unitStats{
	name:         "Patrol",
	layer:        2,
	speed:        110,
	maxHealth:    60,
	img:          assets.ImageDronePatrol,
	allied:       true,
	energyUpkeep: 0.3,
}

var droneGeneratorStats = &unitStats{
	name:      "Mobile generator",
	layer:     2,
	speed:     30,
	maxHealth: 60,
	img:       assets.ImageDroneGenerator,
	allied:    true,
}

var buildingPowerGenerator = &unitStats{
	name:      "Generator",
	layer:     2,
	maxHealth: 75,
	img:       assets.ImageBuildingGenerator,
	allied:    true,
	building:  true,
	ironCost:  1,
	stoneCost: 6,
}

var buildingBarricate = &unitStats{
	name:       "Barricade",
	layer:      2,
	maxHealth:  90,
	img:        assets.ImageBuildingBarricade,
	allied:     true,
	building:   true,
	energyCost: 3,
	ironCost:   1,
}

var buildingSmelter = &unitStats{
	name:       "Smelter",
	layer:      2,
	maxHealth:  130,
	img:        assets.ImageBuildingSmelter,
	allied:     true,
	building:   true,
	energyCost: 3,
	ironCost:   2,
	stoneCost:  5,
}

var creepMutantBase = &unitStats{
	name:      "Mutant base",
	layer:     2,
	maxHealth: 300,
	img:       assets.ImageBuildingMutantBase,
	building:  true,
}

var creepMutantWarrior = &unitStats{
	name:      "Mutant warrior",
	layer:     1,
	maxHealth: 15,
	img:       assets.ImageMutantWarrior,
	speed:     15,
	tiny:      true,
}

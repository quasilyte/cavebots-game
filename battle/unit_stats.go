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

	weapon *weaponStats

	score       int
	botPriority float64

	energyCost int
	ironCost   int
	stoneCost  int
}

type weaponStats struct {
	damage          float64
	attackRange     float64
	reload          float64
	projectileSpeed float64
	projectileImage resource.ImageID
	fireSound       resource.AudioID
	impactSound     resource.AudioID
}

var droneCoreStats = &unitStats{
	name:         "Core",
	botPriority:  10,
	layer:        2,
	speed:        80,
	maxHealth:    100,
	img:          assets.ImageDroneCore,
	allied:       true,
	energyUpkeep: 0,
}

var droneHarvesterStats = &unitStats{
	name:         "Harvester",
	botPriority:  3,
	layer:        2,
	speed:        96,
	maxHealth:    30,
	img:          assets.ImageDroneHarvester,
	allied:       true,
	energyUpkeep: 0.2,
}

var dronePatrolStats = &unitStats{
	name:         "Patrol",
	botPriority:  0.5,
	layer:        2,
	speed:        110,
	maxHealth:    60,
	img:          assets.ImageDronePatrol,
	allied:       true,
	energyUpkeep: 0.3,
	weapon: &weaponStats{
		damage:          8,
		reload:          1.8,
		attackRange:     160,
		projectileSpeed: 550,
		projectileImage: assets.ImagePatrolLaserProjectile,
		fireSound:       assets.AudioPatrolLaser1,
	},
}

var droneGeneratorStats = &unitStats{
	name:        "Mobile generator",
	botPriority: 5.0,
	layer:       2,
	speed:       30,
	maxHealth:   60,
	img:         assets.ImageDroneGenerator,
	allied:      true,
}

var buildingPowerGenerator = &unitStats{
	name:        "Generator",
	botPriority: 11.0,
	layer:       2,
	maxHealth:   75,
	img:         assets.ImageBuildingGenerator,
	allied:      true,
	building:    true,
	ironCost:    1,
	stoneCost:   6,
}

var buildingBarricate = &unitStats{
	name:        "Barricade",
	botPriority: 1,
	layer:       2,
	maxHealth:   90,
	img:         assets.ImageBuildingBarricade,
	allied:      true,
	building:    true,
	energyCost:  3,
	ironCost:    1,
}

var buildingSmelter = &unitStats{
	name:        "Smelter",
	botPriority: 6,
	layer:       2,
	maxHealth:   130,
	img:         assets.ImageBuildingSmelter,
	allied:      true,
	building:    true,
	energyCost:  3,
	ironCost:    2,
	stoneCost:   5,
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
	score:     2,
	layer:     1,
	maxHealth: 15,
	img:       assets.ImageMutantWarrior,
	speed:     15,
	tiny:      true,
	weapon: &weaponStats{
		damage:      3,
		reload:      1.2,
		impactSound: assets.AudioWarriorHit1,
	},
}

var creepMutantWarlord = &unitStats{
	name:      "Mutant warlord",
	score:     3,
	layer:     1,
	maxHealth: 25,
	img:       assets.ImageMutantWarlord,
	speed:     20,
	tiny:      true,
	weapon: &weaponStats{
		damage:      8,
		reload:      1.4,
		impactSound: assets.AudioWarriorHit1,
	},
}

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
	arcPower        float64
	burstSize       int
	burstDelay      float64
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
	energyCost:   2,
	ironCost:     2,
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
	ironCost:     8,
	weapon: &weaponStats{
		damage:          8,
		reload:          1.8,
		attackRange:     140,
		projectileSpeed: 550,
		burstSize:       1,
		projectileImage: assets.ImagePatrolLaserProjectile,
		fireSound:       assets.AudioPatrolLaser1,
	},
}

var droneVanguardStats = &unitStats{
	name:         "Vanguard",
	botPriority:  0.5,
	layer:        2,
	speed:        90,
	maxHealth:    50,
	img:          assets.ImageDroneVanguard,
	allied:       true,
	energyUpkeep: 0.6,
	energyCost:   3,
	ironCost:     9,
	weapon: &weaponStats{
		damage:          10,
		reload:          2.2,
		attackRange:     165,
		projectileSpeed: 600,
		burstSize:       1,
		projectileImage: assets.ImageVanguardProjectile,
		fireSound:       assets.AudioVanguardShot1,
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

var buildingFactory = &unitStats{
	name:        "Factory",
	botPriority: 14,
	layer:       3,
	maxHealth:   120,
	img:         assets.ImageBuildingFactory,
	allied:      true,
	building:    true,
	stoneCost:   14,
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

var creepMutantHunter = &unitStats{
	name:      "Mutant hunter",
	score:     4,
	layer:     1,
	maxHealth: 14,
	img:       assets.ImageMutantHunter,
	speed:     25,
	tiny:      true,
	weapon: &weaponStats{
		damage:          5,
		reload:          2.0,
		attackRange:     120,
		projectileSpeed: 350,
		burstSize:       1,
		projectileImage: assets.ImageArrowProjectile,
		fireSound:       assets.AudioBowShot1,
		arcPower:        2.5,
	},
}

var creepMutantWarlord = &unitStats{
	name:      "Mutant warlord",
	score:     4,
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

var creepJeep = &unitStats{
	name:      "Monster",
	layer:     2,
	maxHealth: 140,
	img:       assets.ImageJeep,
	building:  true,
	weapon: &weaponStats{
		damage:          3,
		reload:          1.9,
		attackRange:     230,
		burstSize:       3,
		burstDelay:      0.1,
		projectileSpeed: 600,
		projectileImage: assets.ImageGatlingProjectile,
		fireSound:       assets.AudioGatlingShot,
	},
}

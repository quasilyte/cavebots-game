package battle

import (
	"github.com/quasilyte/cavebots-game/assets"
	resource "github.com/quasilyte/ebitengine-resource"
)

type unitStats struct {
	speed        float64
	img          resource.ImageID
	allied       bool
	energyUpkeep float64
}

var droneCoreStats = &unitStats{
	speed:        80,
	img:          assets.ImageDroneCore,
	allied:       true,
	energyUpkeep: 0,
}

var droneHarvesterStats = &unitStats{
	speed:        96,
	img:          assets.ImageDroneHarvester,
	allied:       true,
	energyUpkeep: 0.2,
}

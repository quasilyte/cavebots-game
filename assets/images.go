package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"

	_ "image/png"
)

func registerImageResources(ctx *ge.Context) {
	imageResources := map[resource.ImageID]resource.ImageInfo{
		ImageUIButtonDisabled:      {Path: "images/ebitenui/button-disabled.png"},
		ImageUIButtonIdle:          {Path: "images/ebitenui/button-idle.png"},
		ImageUIButtonHover:         {Path: "images/ebitenui/button-hover.png"},
		ImageUIButtonPressed:       {Path: "images/ebitenui/button-pressed.png"},
		ImageUISelectButtonIdle:    {Path: "images/ebitenui/select-button-idle.png"},
		ImageUISelectButtonHover:   {Path: "images/ebitenui/select-button-hover.png"},
		ImageUISelectButtonPressed: {Path: "images/ebitenui/select-button-pressed.png"},
		ImageUIPanelIdle:           {Path: "images/ebitenui/panel-idle.png"},

		ImageCaveTiles:   {Path: "images/cave_tiles.png"},
		ImageForestTiles: {Path: "images/forest_tiles.png"},

		ImageCellSelector: {Path: "images/cell_selector.png"},

		ImageDroneCore:      {Path: "images/drone_core.png"},
		ImageDroneHarvester: {Path: "images/drone_harvester.png"},
		ImageDronePatrol:    {Path: "images/drone_patrol.png"},
		ImageDroneGenerator: {Path: "images/drone_mobile_generator.png"},
		ImageDroneVanguard:  {Path: "images/drone_vanguard.png"},
		ImageDroneRepair:    {Path: "images/drone_repair.png"},
		ImageDroneTitan:     {Path: "images/drone_titan.png"},

		ImageBuildingGenerator: {Path: "images/building_generator.png"},
		ImageBuildingBarricade: {Path: "images/building_wall.png"},
		ImageBuildingSmelter:   {Path: "images/building_smelter.png"},
		ImageBuildingFactory:   {Path: "images/building_factory.png"},
		ImageBuildingTurret:    {Path: "images/building_turret.png"},

		ImageBuildingMutantBase: {Path: "images/mutant_base.png", FrameWidth: 34},

		ImageMutantWarrior: {Path: "images/mutant_warrior.png", FrameWidth: 32},
		ImageMutantHunter:  {Path: "images/mutant_hunter.png", FrameWidth: 32},
		ImageMutantWarlord: {Path: "images/mutant_warlord.png", FrameWidth: 32},
		ImageJeep:          {Path: "images/jeep.png", FrameWidth: 32},

		ImageIronResource: {Path: "images/iron_resource.png"},

		ImageMountains:       {Path: "images/mountains.png", FrameWidth: 48},
		ImageHardMountains:   {Path: "images/hard_mountains.png", FrameWidth: 48},
		ImageIronMountains:   {Path: "images/iron_mountains.png", FrameWidth: 48},
		ImageWeakMountains:   {Path: "images/weak_mountains.png", FrameWidth: 48},
		ImageUnitMountains:   {Path: "images/unit_mountains.png", FrameWidth: 48},
		ImageRockyMountains:  {Path: "images/rocky_mountains.png", FrameWidth: 48},
		ImageHardTerrain:     {Path: "images/hard_terrain.png", FrameWidth: 32},
		ImageBiomeTransition: {Path: "images/transition.png"},

		ImagePatrolLaserProjectile: {Path: "images/patrol_laser_projectile.png"},
		ImageVanguardProjectile:    {Path: "images/vanguard_projectile.png"},
		ImageArrowProjectile:       {Path: "images/arrow_projectile.png"},
		ImageGatlingProjectile:     {Path: "images/gatling_projectile.png"},

		ImageEffectSmokeUp:      {Path: "images/smoke_up.png", FrameWidth: 8},
		ImageEffectRepair:       {Path: "images/repair_effect.png", FrameWidth: 16},
		ImageEffectExplosion:    {Path: "images/explosion_effect.png", FrameWidth: 16},
		ImageEffectMutantImpact: {Path: "images/mutant_impact.png", FrameWidth: 16},
	}

	for id, res := range imageResources {
		ctx.Loader.ImageRegistry.Set(id, res)
		ctx.Loader.LoadImage(id)
	}
}

const (
	ImageNone resource.ImageID = iota

	ImageUIButtonDisabled
	ImageUIButtonIdle
	ImageUIButtonHover
	ImageUIButtonPressed
	ImageUISelectButtonIdle
	ImageUISelectButtonHover
	ImageUISelectButtonPressed
	ImageUIPanelIdle

	ImageCaveTiles
	ImageForestTiles

	ImageCellSelector

	ImageDroneCore
	ImageDroneHarvester
	ImageDronePatrol
	ImageDroneGenerator
	ImageDroneVanguard
	ImageDroneRepair
	ImageDroneTitan

	ImageMutantWarrior
	ImageMutantHunter
	ImageMutantWarlord
	ImageJeep

	ImageBuildingBarricade
	ImageBuildingGenerator
	ImageBuildingSmelter
	ImageBuildingFactory
	ImageBuildingTurret

	ImageBuildingMutantBase

	ImageIronResource

	ImageMountains
	ImageHardMountains
	ImageIronMountains
	ImageWeakMountains
	ImageUnitMountains
	ImageRockyMountains
	ImageHardTerrain
	ImageBiomeTransition

	ImagePatrolLaserProjectile
	ImageVanguardProjectile
	ImageArrowProjectile
	ImageGatlingProjectile

	ImageEffectSmokeUp
	ImageEffectRepair
	ImageEffectExplosion
	ImageEffectMutantImpact
)

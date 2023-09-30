package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"

	_ "image/png"
)

func registerImageResources(ctx *ge.Context) {
	imageResources := map[resource.ImageID]resource.ImageInfo{
		ImageUIButtonDisabled: {Path: "images/ebitenui/button-disabled.png"},
		ImageUIButtonIdle:     {Path: "images/ebitenui/button-idle.png"},
		ImageUIButtonHover:    {Path: "images/ebitenui/button-hover.png"},
		ImageUIButtonPressed:  {Path: "images/ebitenui/button-pressed.png"},

		ImageCaveTiles:   {Path: "images/cave_tiles.png"},
		ImageForestTiles: {Path: "images/forest_tiles.png"},

		ImageCellSelector: {Path: "images/cell_selector.png"},

		ImageDroneCore:      {Path: "images/drone_core.png"},
		ImageDroneHarvester: {Path: "images/drone_harvester.png"},
		ImageDronePatrol:    {Path: "images/drone_patrol.png"},

		ImageBuildingGenerator: {Path: "images/building_generator.png"},
		ImageBuildingBarricade: {Path: "images/building_wall.png"},
		ImageBuildingSmelter:   {Path: "images/building_smelter.png"},

		ImageIronResource: {Path: "images/iron_resource.png"},

		ImageMountains:   {Path: "images/mountains.png", FrameWidth: 48},
		ImageHardTerrain: {Path: "images/hard_terrain.png", FrameWidth: 32},
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

	ImageCaveTiles
	ImageForestTiles

	ImageCellSelector

	ImageDroneCore
	ImageDroneHarvester
	ImageDronePatrol

	ImageBuildingBarricade
	ImageBuildingGenerator
	ImageBuildingSmelter

	ImageIronResource

	ImageMountains
	ImageHardTerrain
)

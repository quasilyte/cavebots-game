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
)

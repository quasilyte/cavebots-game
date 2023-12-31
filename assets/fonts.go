package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
)

func registerFontResources(ctx *ge.Context) {
	fontResources := map[resource.FontID]resource.FontInfo{
		FontTiny:   {Path: "fonts/whiterabbit.ttf", Size: 10, LineSpacing: 2},
		FontSmall:  {Path: "fonts/whiterabbit.ttf", Size: 14, LineSpacing: 2},
		FontNormal: {Path: "fonts/whiterabbit.ttf", Size: 18},
		FontBig:    {Path: "fonts/whiterabbit.ttf", Size: 26},
	}

	for id, res := range fontResources {
		ctx.Loader.FontRegistry.Set(id, res)
		ctx.Loader.LoadFont(id)
	}
}

const (
	FontNone resource.FontID = iota

	FontTiny
	FontSmall
	FontNormal
	FontBig
)

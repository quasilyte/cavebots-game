package eui

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/cavebots-game/styles"
	resource "github.com/quasilyte/ebitengine-resource"
	"golang.org/x/image/font"
)

type Resources struct {
	button *buttonResource
}

type buttonResource struct {
	Image         *widget.ButtonImage
	Padding       widget.Insets
	TextColors    *widget.ButtonTextColor
	AltTextColors *widget.ButtonTextColor
	FontFace      font.Face
}

func PrepareResources(loader *resource.Loader) *Resources {
	result := &Resources{}

	normalFont := loader.LoadFont(assets.FontNormal).Face

	{
		disabled := nineSliceImage(loader.LoadImage(assets.ImageUIButtonDisabled).Data, 12, 0)
		idle := nineSliceImage(loader.LoadImage(assets.ImageUIButtonIdle).Data, 12, 0)
		hover := nineSliceImage(loader.LoadImage(assets.ImageUIButtonHover).Data, 12, 0)
		pressed := nineSliceImage(loader.LoadImage(assets.ImageUIButtonPressed).Data, 12, 0)
		buttonPadding := widget.Insets{
			Left:  30,
			Right: 30,
		}
		buttonColors := &widget.ButtonTextColor{
			Idle:     styles.ButtonTextColor,
			Disabled: styles.DisabledButtonTextColor,
		}
		result.button = &buttonResource{
			Image: &widget.ButtonImage{
				Idle:     idle,
				Hover:    hover,
				Pressed:  pressed,
				Disabled: disabled,
			},
			Padding:    buttonPadding,
			TextColors: buttonColors,
			AltTextColors: &widget.ButtonTextColor{
				Idle:     styles.ButtonTextColor,
				Disabled: styles.DisabledButtonTextColor,
			},
			FontFace: normalFont,
		}
	}

	return result
}

func nineSliceImage(i *ebiten.Image, centerWidth, centerHeight int) *image.NineSlice {
	w, h := i.Size()
	return image.NewNineSlice(i,
		[3]int{(w - centerWidth) / 2, centerWidth, w - (w-centerWidth)/2 - centerWidth},
		[3]int{(h - centerHeight) / 2, centerHeight, h - (h-centerHeight)/2 - centerHeight})
}

type ButtonConfig struct {
	Text         string
	TextAltColor bool
	OnClick      func()
	LayoutData   any
	MinWidth     int
	Font         font.Face
}

func NewButtonWithConfig(res *Resources, config ButtonConfig) *widget.Button {
	ff := config.Font
	if ff == nil {
		ff = res.button.FontFace
	}
	options := []widget.ButtonOpt{
		widget.ButtonOpts.Image(res.button.Image),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			if config.OnClick != nil {
				config.OnClick()
			}
		}),
	}
	colors := res.button.TextColors
	if config.TextAltColor {
		colors = res.button.AltTextColors
	}
	options = append(options,
		widget.ButtonOpts.Text(config.Text, ff, colors),
		widget.ButtonOpts.TextPadding(res.button.Padding))
	if config.LayoutData != nil {
		options = append(options, widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(config.LayoutData)))
	}
	if config.MinWidth != 0 {
		options = append(options, widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(config.MinWidth, 0)))
	}
	return widget.NewButton(options...)
}

func NewButton(res *Resources, text string, onclick func()) *widget.Button {
	return NewButtonWithConfig(res, ButtonConfig{
		Text:    text,
		OnClick: onclick,
	})
}

func NewCenteredLabelWithMaxWidth(text string, ff font.Face, width float64) *widget.Text {
	options := []widget.TextOpt{
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
		),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.Text(text, ff, styles.ButtonTextColor),
	}
	if width != -1 {
		options = append(options, widget.TextOpts.MaxWidth(width))
	}
	return widget.NewText(options...)
}

func NewCenteredLabel(text string, ff font.Face) *widget.Text {
	return NewCenteredLabelWithMaxWidth(text, ff, -1)
}

func NewSeparator(ld interface{}, clr color.RGBA) widget.PreferredSizeLocateableWidget {
	c := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}))),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(ld)))

	c.AddChild(widget.NewGraphic(
		widget.GraphicOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch:   true,
			MaxHeight: 2,
		})),
		widget.GraphicOpts.ImageNineSlice(image.NewNineSliceColor(clr)),
	))

	return c
}

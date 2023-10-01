package eui

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/cavebots-game/assets"
	"github.com/quasilyte/cavebots-game/styles"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/gmath"
	"golang.org/x/image/font"
)

type Resources struct {
	button         *buttonResource
	buttonSelected *buttonResource
	panel          *panelResource
}

type buttonResource struct {
	Image         *widget.ButtonImage
	Padding       widget.Insets
	TextColors    *widget.ButtonTextColor
	AltTextColors *widget.ButtonTextColor
	FontFace      font.Face
}

type panelResource struct {
	Image   *image.NineSlice
	Padding widget.Insets
}

func PrepareResources(loader *resource.Loader) *Resources {
	result := &Resources{}

	normalFont := loader.LoadFont(assets.FontNormal).Face

	{
		disabled := nineSliceImage(loader.LoadImage(assets.ImageUIButtonDisabled).Data, 12, 0)
		idle := nineSliceImage(loader.LoadImage(assets.ImageUIButtonIdle).Data, 12, 0)
		hover := nineSliceImage(loader.LoadImage(assets.ImageUIButtonHover).Data, 12, 0)
		pressed := nineSliceImage(loader.LoadImage(assets.ImageUIButtonPressed).Data, 12, 0)
		selectedIdle := nineSliceImage(loader.LoadImage(assets.ImageUISelectButtonIdle).Data, 12, 0)
		selectedHover := nineSliceImage(loader.LoadImage(assets.ImageUISelectButtonHover).Data, 12, 0)
		selectedPressed := nineSliceImage(loader.LoadImage(assets.ImageUISelectButtonPressed).Data, 12, 0)
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
		result.buttonSelected = &buttonResource{
			Image: &widget.ButtonImage{
				Idle:    selectedIdle,
				Hover:   selectedHover,
				Pressed: selectedPressed,
			},
			Padding:    buttonPadding,
			TextColors: buttonColors,
			FontFace:   normalFont,
		}
	}

	{
		idle := loader.LoadImage(assets.ImageUIPanelIdle).Data
		result.panel = &panelResource{
			Image: nineSliceImage(idle, 10, 10),
			Padding: widget.Insets{
				Left:   16,
				Right:  16,
				Top:    10,
				Bottom: 10,
			},
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

func NewColoredLabel(text string, ff font.Face, clr color.RGBA, options ...widget.TextOpt) *widget.Text {
	opts := []widget.TextOpt{
		widget.TextOpts.Text(text, ff, clr),
	}
	if len(options) != 0 {
		opts = append(opts, options...)
	}
	return widget.NewText(opts...)
}

func NewLabel(text string, ff font.Face, options ...widget.TextOpt) *widget.Text {
	return NewColoredLabel(text, ff, styles.ButtonTextColor, options...)
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

func NewPanelWithPadding(res *Resources, minWidth, minHeight int, padding widget.Insets) *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(res.panel.Image),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(padding),
		)),
		// widget.ContainerOpts.Layout(widget.NewRowLayout(
		// 	widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		// 	widget.RowLayoutOpts.Spacing(4),
		// 	widget.RowLayoutOpts.Padding(res.panel.Padding),
		// )),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				StretchHorizontal:  true,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.MinSize(minWidth, minHeight),
		),
	)
}

type SelectButtonConfig struct {
	Resources *Resources
	Input     *input.Handler

	Value      *int
	Label      string
	ValueNames []string

	LayoutData any

	OnPressed func()
}

func NewSelectButton(config SelectButtonConfig) *widget.Button {
	maxValue := len(config.ValueNames) - 1
	value := config.Value
	key := config.Label
	valueNames := config.ValueNames

	var slider gmath.Slider
	slider.SetBounds(0, maxValue)
	slider.TrySetValue(*value)
	makeLabel := func() string {
		if key == "" {
			return valueNames[slider.Value()]
		}
		return key + ": " + valueNames[slider.Value()]
	}

	buttonOpts := []widget.ButtonOpt{}
	if config.LayoutData != nil {
		buttonOpts = append(buttonOpts, widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(config.LayoutData)))
	}
	button := newButtonSelected(config.Resources, makeLabel(), buttonOpts...)

	button.ClickedEvent.AddHandler(func(args interface{}) {
		increase := false
		{
			cursorPos := config.Input.CursorPos()
			buttonRect := button.GetWidget().Rect
			buttonWidth := buttonRect.Dx()
			if cursorPos.X >= float64(buttonRect.Min.X)+float64(buttonWidth)*0.5 {
				increase = true
			}
		}

		if increase {
			slider.Inc()
		} else {
			slider.Dec()
		}
		*value = slider.Value()

		button.Text().Label = makeLabel()
		if config.OnPressed != nil {
			config.OnPressed()
		}
	})

	return button
}

func newButtonSelected(res *Resources, text string, opts ...widget.ButtonOpt) *widget.Button {
	options := []widget.ButtonOpt{
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.ButtonOpts.Image(res.buttonSelected.Image),
		widget.ButtonOpts.Text(text, res.buttonSelected.FontFace, res.buttonSelected.TextColors),
		widget.ButtonOpts.TextPadding(res.buttonSelected.Padding),
	}
	options = append(options, opts...)
	return widget.NewButton(options...)
}

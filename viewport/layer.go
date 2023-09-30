package viewport

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type LayerContainer struct {
	belowObjects         layer
	objects              layer
	slightlyAboveObjects layer
	aboveObjects         layer
}

type layer struct {
	sprites []*ge.Sprite
	objects []cameraObject
}

type cameraObject interface {
	DrawWithOffset(dst *ebiten.Image, offset gmath.Vec)
	IsDisposed() bool
	BoundsRect() gmath.Rect
}

func (c *LayerContainer) AddSprite(s *ge.Sprite) {
	c.objects.AddSprite(s)
}

func (c *LayerContainer) AddGraphics(o cameraObject) {
	c.objects.Add(o)
}

func (c *LayerContainer) AddSpriteSlightlyAbove(s *ge.Sprite) {
	c.slightlyAboveObjects.AddSprite(s)
}

func (c *LayerContainer) AddSpriteAbove(s *ge.Sprite) {
	c.aboveObjects.AddSprite(s)
}

func (c *LayerContainer) AddGraphicsSlightlyAbove(o cameraObject) {
	c.slightlyAboveObjects.Add(o)
}

func (c *LayerContainer) AddGraphicsAbove(o cameraObject) {
	c.aboveObjects.Add(o)
}

func (c *LayerContainer) AddSpriteBelow(s *ge.Sprite) {
	c.belowObjects.AddSprite(s)
}

func (l *layer) Add(o cameraObject) {
	l.objects = append(l.objects, o)
}

func (l *layer) AddSprite(s *ge.Sprite) {
	l.sprites = append(l.sprites, s)
}

func (l *layer) filter() {
	liveSprites := l.sprites[:0]
	for _, s := range l.sprites {
		if s.IsDisposed() {
			continue
		}
		liveSprites = append(liveSprites, s)
	}
	l.sprites = liveSprites

	if len(l.objects) != 0 {
		liveObjects := l.objects[:0]
		for _, o := range l.objects {
			if o.IsDisposed() {
				continue
			}
			liveObjects = append(liveObjects, o)
		}
		l.objects = liveObjects
	}
}

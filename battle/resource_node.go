package battle

import (
	"github.com/quasilyte/cavebots-game/assets"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gsignal"
)

type resourceKind int

const (
	resourceUnknown resourceKind = iota
	resourceIron
)

type resourceStats struct {
	name string
	kind resourceKind
	img  resource.ImageID
}

var ironResourceStats = &resourceStats{
	name: "Iron",
	kind: resourceIron,
	img:  assets.ImageIronResource,
}

type resourceNode struct {
	stats *resourceStats

	world *worldState

	sprite *ge.Sprite

	pos gmath.Vec

	amount int

	EventDisposed gsignal.Event[*resourceNode]
}

func newResourceNode(world *worldState, pos gmath.Vec, stats *resourceStats, amount int) *resourceNode {
	return &resourceNode{
		world:  world,
		pos:    pos,
		stats:  stats,
		amount: amount,
	}
}

func (r *resourceNode) Init(scene *ge.Scene) {
	r.sprite = scene.NewSprite(r.stats.img)
	r.sprite.Pos.Base = &r.pos
	r.world.stage.AddSprite(r.sprite)
}

func (r *resourceNode) IsDisposed() bool {
	return r.sprite.IsDisposed()
}

func (r *resourceNode) Dispose() {
	if r.IsDisposed() {
		return
	}
	r.EventDisposed.Emit(r)
	r.sprite.Dispose()
}

func (r *resourceNode) Update(delta float64) {}

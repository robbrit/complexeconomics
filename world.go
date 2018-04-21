package econerra

import (
	"math/rand"
)

type World struct {
	markets []Market

	r *rand.Rand
}

type WorldOptions struct {
	Rand *rand.Rand
}

func NewWorld(opts WorldOptions) *World {
	w := &World{
		markets: make([]Market, NumGoods),
		r:       opts.Rand,
	}

	for g := range w.markets {
		w.markets[g] = NewMarket(w, Good(g))
	}
	return w
}

func (w *World) Market(g Good) Market {
	return w.markets[g]
}

func (w *World) Clear() {
	for _, m := range w.markets {
		m.Clear()
	}
}

func (w *World) Rand() *rand.Rand  { return w.r }
func (w *World) Markets() []Market { return w.markets }

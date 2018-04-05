package econerra

type World struct {
	markets []Market
}

func NewWorld() *World {
	return &World{
		markets: make([]Market, NumGoods),
	}
}

func (w *World) Market(g Good) Market {
	return w.markets[g]
}

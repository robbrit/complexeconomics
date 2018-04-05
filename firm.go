package econerra

import (
	"fmt"
	"math"
)

// FirmID is a unique identifier for a firm.
type FirmID uint64

// A Firm is a firm in the econerra world.
type Firm struct {
	agent

	id        FirmID
	totShares Size
	good      Good
}

const (
	// Technology modifier for production function.
	tech = 10.0
	// Exponent for C-D production.
	labElast = 0.5
)

// NewFirm constructs a new firm.
func NewFirm(id FirmID, w World, totShares Size, adjust Price, good Good, initInv Size) (*Firm, error) {
	if good == Labour {
		return nil, fmt.Errorf("firms can't produce labour")
	}

	f := &Firm{
		id:        id,
		totShares: totShares,
		good:      good,
	}

	f.agent = *newAgent(w, adjust, f)

	f.inventory[good] = initInv

	return f, nil
}

// Act causes the firm to make its decisions for this cycle.
func (f *Firm) Act() {
	f.agent.act()

	sellPrice := f.prices[f.good]
	sumInputPrices := 0.0

	for _, in := range inputGoods[f.good] {
		sumInputPrices += float64(f.prices[in.g]) * float64(in.s)
	}
	denom := tech * labElast * (float64(sellPrice) - sumInputPrices)

	// Post order of labour based on labour demand function.
	lf := math.Pow(float64(f.prices[Labour])/denom, 1.0/(labElast-1.0))
	l := f.labour(lf)

	f.world.Market(Labour).Post(&MarketOrder{
		Price: f.prices[Labour],
		Size:  l,
		Side:  Buy,
		Owner: f,
	})

	// Post orders for input goods based on input demand function.
	for _, in := range inputGoods[f.good] {
		ind := float64(in.s) * tech * math.Pow(float64(l), labElast)
		f.world.Market(in.g).Post(&MarketOrder{
			Price: f.prices[in.g],
			// The floor of the input demand always maximizes profit.
			Size:  Size(math.Floor(ind)) - f.inventory[in.g],
			Side:  Buy,
			Owner: f,
		})
	}

	// Post order to sell what we produced last round.
	f.world.Market(f.good).Post(&MarketOrder{
		Price: f.prices[f.good],
		Size:  f.inventory[f.good],
		Side:  Sell,
		Owner: f,
	})
}

// labour returns the floor or the ceil of the profit-maximizing labour demand.
func (f *Firm) labour(l float64) Size {
	cl := math.Ceil(l)
	fl := math.Floor(l)

	pc := f.profit(cl)
	pf := f.profit(fl)
	if pc > pf {
		return Size(cl)
	}
	return Size(fl)
}

// profit returns the profit for a given labour level.
func (f *Firm) profit(l float64) float64 {
	ip := int64(0) // Total cost of inputs.
	for _, in := range inputGoods[f.good] {
		ip += int64(f.prices[in.g]) * int64(in.s)
	}
	p := int64(f.prices[f.good])
	w := float64(f.prices[Labour])
	la := math.Pow(l, labElast)
	return tech*la*float64(p-ip) - w*l
}

// OnFill is an implementation of MarketAgent.OnFill.
func (f *Firm) OnFill(g Good, side Side, p Price, s Size, sig MarketSignal) {
	f.agent.onFill(g, side, p, s, sig)
}

// OnUnfilled is an implementation of MarketAgent.OnUnfilled.
func (f *Firm) OnUnfilled(g Good, side Side, p Price, s Size, sig MarketSignal) {
	f.agent.onUnfilled(g, side, p, s, sig)
}

// IsBuyer is an implementation of AgentStrategy.IsBuyer.
func (f *Firm) IsBuyer(g Good) bool {
	// TODO - some firm types buy input goods
	return g == Labour
}

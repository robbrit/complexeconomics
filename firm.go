package econerra

import (
	"fmt"
	"log"
	"math"
)

var _ = log.Printf

// FirmID is a unique identifier for a firm.
type FirmID uint64

// A Firm is a firm in the econerra world.
type Firm struct {
	agent

	id   FirmID
	good Good
}

// FirmOptions are the arguments used to construct a firm.
type FirmOptions struct {
	AgentOptions

	// ID is the unique identifier for the firm.
	ID FirmID
	// Good is the good that this firm produces.
	Good Good
	// InitInventory is the initial inventory that the firm has.
	InitInventory Size
}

const (
	// Exponent for C-D production.
	labElast = 0.5
)

var (
	// Technology modifier for production function.
	tech = map[Good]float64{
		Grain:      40.0,
		Vegetables: 30.0,
		Cotton:     30.0,
		Meat:       8.0,
		Beer:       10.0,
		Clothing:   15.0,
	}
)

// NewFirm constructs a new firm.
func NewFirm(w *World, opts FirmOptions) (*Firm, error) {
	if opts.Good == Labour {
		return nil, fmt.Errorf("firms can't produce labour")
	}

	f := &Firm{
		id:   opts.ID,
		good: opts.Good,
	}

	f.agent = *newAgent(w, opts.AgentOptions, f)

	f.inventory[opts.Good] = opts.InitInventory

	return f, nil
}

// produce produces whatever good this firm needs.
func (f *Firm) produce() {
	// Figure out how much we can produce based on the labour we got last round.
	q := tech[f.good] * math.Pow(float64(f.inventory[Labour]), labElast)

	// For each input good, see if we can produce q with that much.
	for _, input := range inputGoods[f.good] {
		maxQ := float64(f.inventory[input.g]) / float64(input.s)
		if maxQ < q {
			q = maxQ
		}
	}

	// We now know how much we can produce. Add that much to our inventory, and
	// reduce our input stock.
	q = math.Floor(q)

	f.inventory[f.good] += Size(q)
	for _, input := range inputGoods[f.good] {
		f.inventory[input.g] -= Size(q) * input.s
	}
}

// Act causes the firm to make its decisions for this cycle.
func (f *Firm) Act() {
	t := tech[f.good]
	f.agent.act()

	// Based on what we had last round, produce goods for market.
	log.Printf("%s: before: %v", f.good, f.inventory)
	f.produce()
	log.Printf("%s: after: %v", f.good, f.inventory)

	// Reset any variables that should be reset.
	f.inventory[Labour] = 0

	sellPrice := f.prices[f.good]
	sumInputPrices := 0.0

	for _, in := range inputGoods[f.good] {
		sumInputPrices += f.prices[in.g] * float64(in.s)
	}
	denom := t * labElast * (sellPrice - sumInputPrices)

	// Post order of labour based on labour demand function.
	lf := math.Pow(f.prices[Labour]/denom, 1.0/(labElast-1.0))
	l := f.labour(lf)
	log.Printf("%s: want %d workers", f.good, l)

	f.world.Market(Labour).Post(&MarketOrder{
		Price: f.prices[Labour],
		Size:  l,
		Side:  Buy,
		Owner: f,
	})

	// Post orders for input goods based on input demand function.
	for _, in := range inputGoods[f.good] {
		ind := float64(in.s) * t * math.Pow(float64(l), labElast)
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
	ip := 0.0 // Total cost of inputs.
	for _, in := range inputGoods[f.good] {
		ip += float64(f.prices[in.g]) * float64(in.s)
	}
	p := float64(f.prices[f.good])
	w := float64(f.prices[Labour])
	la := math.Pow(l, labElast)
	return tech[f.good]*la*(p-ip) - w*l
}

// OnFill is an implementation of MarketAgent.OnFill.
func (f *Firm) OnFill(g Good, side Side, p float64, s Size, sig MarketSignal) {
	f.agent.onFill(g, side, p, s, sig)
}

// OnUnfilled is an implementation of MarketAgent.OnUnfilled.
func (f *Firm) OnUnfilled(g Good, side Side, p float64, s Size, sig MarketSignal) {
	f.agent.onUnfilled(g, side, p, s, sig)
}

// IsBuyer is an implementation of AgentStrategy.IsBuyer.
func (f *Firm) IsBuyer(g Good) bool {
	// TODO - some firm types buy input goods
	return g == Labour
}

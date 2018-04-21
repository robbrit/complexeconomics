package econerra

import (
	"log"
	"math"

	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/mat"
)

var _ = log.Printf

// IndividualID is a unique identifier for an individual.
type IndividualID uint64

// An Individual is an individual in the econerra world.
type Individual struct {
	agent

	id         IndividualID
	ascentRate float64
	numSteps   int

	lastSizes []float64
}

// IndividualOptions have all the values needed for an individual.
type IndividualOptions struct {
	AgentOptions

	// ID is the unique identifier for the individual.
	ID IndividualID
	// AscentRate is the gradient-ascent step value.
	AscentRate float64
	// NumSteps is the number of steps taken during gradient ascent.
	NumSteps int
}

// NewIndividual constructs a new individual.
func NewIndividual(w *World, opts IndividualOptions) *Individual {
	i := &Individual{
		id:         opts.ID,
		ascentRate: opts.AscentRate,
		numSteps:   opts.NumSteps,
		lastSizes:  make([]float64, NumGoods),
	}

	i.agent = *newAgent(w, opts.AgentOptions, i)

	// Set each utility good to a random value.
	for g := range UtilityGoods {
		i.lastSizes[g] = 1.0
	}
	return i
}

// Act causes the individual to make its decisions for this cycle.
func (ind *Individual) Act() {
	ind.agent.act()

	// Reset any variables that should be reset.
	ind.inventory[Labour] = 0

	// Based on the current set of prices, the individual does a few steps of a
	// constrained gradient ascent of the utility function.
	g := mat.NewVecDense(NumGoods, nil)
	x := mat.NewVecDense(NumGoods, ind.lastSizes)
	temp := mat.NewVecDense(NumGoods, nil)

	// Calculate unit vector of prices.
	p := mat.NewVecDense(NumGoods, ind.prices)
	pn := mat.NewVecDense(NumGoods, nil)
	n := mat.Norm(p, 2.0)
	if n == 0.0 {
		panic("price vector of zero")
	}
	pn.ScaleVec(1.0/n, p)

	// Calculate base point in the budget hyperplane. We just set an arbitary
	// good to budget / price of that good and the rest of the values to 0.
	// Negate it since we subtract later on.
	orig := mat.NewVecDense(NumGoods, nil)
	orig.SetVec(int(baselineGood), -ind.cash/p.AtVec(int(baselineGood)))

	for i := 0; i < ind.numSteps; i++ {
		// 1) Calculate the gradient at x.
		fd.Gradient(g.RawVector().Data, utility, x.RawVector().Data, nil)
		// 2) Step in the direction of the gradient.
		x.AddScaledVec(x, ind.ascentRate, g)
		// 3) Project x onto the hyperplane.
		temp.AddVec(x, orig)
		d := mat.Dot(temp, pn)
		x.AddScaledVec(temp, -d, pn)
		// 4) Clamp x to be non-negative.
		for g := 0; g < NumGoods; g++ {
			if !UtilityGoods[Good(g)] || x.AtVec(int(g)) < 0.0 {
				x.SetVec(int(g), 0.0)
			}
		}
		// 5) Scale any remaining positive values since we may be over budget now.
		cost := mat.Dot(x, p)
		x.ScaleVec(ind.cash/cost, x)
	}

	// Now we know how much of each good we want at these prices, send orders.
	for g := range UtilityGoods {
		ind.world.Market(g).Post(&MarketOrder{
			Price: ind.prices[g],
			Size:  Size(math.Round(x.AtVec(int(g)))),
			Side:  Buy,
			Owner: ind,
		})
	}

	// Lastly, sell some labour.
	ind.world.Market(Labour).Post(&MarketOrder{
		Price: ind.prices[Labour],
		Size:  1,
		Side:  Sell,
		Owner: ind,
	})
}

// OnFill is an implementation of MarketAgent.OnFill.
func (ind *Individual) OnFill(g Good, side Side, p float64, s Size, sig MarketSignal) {
	ind.agent.onFill(g, side, p, s, sig)
}

// OnUnfilled is an implementation of MarketAgent.OnUnfilled.
func (ind *Individual) OnUnfilled(g Good, side Side, p float64, s Size, sig MarketSignal) {
	ind.agent.onUnfilled(g, side, p, s, sig)
}

// IsBuyer is an implementation of AgentStrategy.IsBuyer.
func (ind *Individual) IsBuyer(g Good) bool { return g != Labour }

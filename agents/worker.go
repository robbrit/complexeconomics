package agents

import (
	"log"
	"math"

	"github.com/robbrit/econerra/goods"
	"github.com/robbrit/econerra/market"
)

var _ = log.Println

// A Worker is an agent that sells labour.
type Worker struct {
	unemployed    bool
	wage          market.Price
	prices        map[goods.Good]market.Price
	demand        map[goods.Good]market.Size
	purchasesMade map[goods.Good]market.Size
}

// NewWorker creates a new worker.
func NewWorker(initialWage, initialPrice market.Price) *Worker {
	w := &Worker{
		unemployed:    true,
		wage:          initialWage,
		prices:        map[goods.Good]market.Price{},
		demand:        map[goods.Good]market.Size{},
		purchasesMade: map[goods.Good]market.Size{},
	}

	for _, good := range goods.AllGoods {
		w.prices[good] = initialPrice
	}

	return w
}

// Act triggers the worker's decision process.
func (w *Worker) Act(p *Parameters) {
	// Step one: adjust prices based on market outcomes last round.
	w.adjustPrices(p)
	// Step two: calculate profit-maximizing values of goods to purchase.
	w.chooseTargets(p)
	// Step three: place orders.
	w.placeOrders(p)
	// Step four: reset information about last iteration.
	w.reset()
}

func (w *Worker) adjustPrices(p *Parameters) {
	if p.LabourMarket.Volume() == 0 {
		// This is the first iteration, just use the prices we have.
		return
	}

	if w.unemployed {
		// I was unemployed last round, undercut the market.
		w.wage = p.LabourMarket.Low() - p.Increment
	} else if w.wage <= p.LabourMarket.High() {
		// I was employed last round, ask for a higher price if I can.
		w.wage += p.Increment
	}

	for _, good := range goods.AllGoods {
		amountBought := w.purchasesMade[good]
		demand := w.demand[good]

		price := w.prices[good]
		market := p.Goods[good].Market

		if amountBought < demand {
			// Didn't get enough, raise price.
			w.prices[good] += p.Increment
			// TODO(rob): A more intelligent agent would probably look at the market to make a
			// decision.
		} else if price >= market.Low() {
			// Got enough last time, undercut the market a little bit.
			w.prices[good] = market.Low() - p.Increment
		}
	}
}

func (w *Worker) chooseTargets(p *Parameters) {
	if w.unemployed {
		// If we didn't work, we can't buy things.
		for _, good := range goods.AllGoods {
			w.demand[good] = 0
		}
		return
	}

	/* Based on the prices we set, choose the utility maximizing quantities that satisfy the budget
	constraint. Using CES Utility:

		U = A(sum beta_i * x_i^rho)^(k/rho)
		k = sum(beta_i)

	and subjecting that to the budget constraint:

		sum(p_i * x_i) <= m

	we get demand functions:

		x_i = (beta_i / p_i)^sigma * m / sum(beta_i^sigma * p_i^(1-sigma))
		sigma = 1 / (1 - rho)

	In code:
		sigma is p.Elasticity
		beta_i is p.Goods[i].Share
		m is w.wage
		p_i is w.prices[i]
	*/

	denominator := 0.0
	for _, good := range goods.AllGoods {
		denominator += p.Goods[good].Share * math.Pow(float64(w.prices[good]), 1.0-p.Elasticity)
	}
	for _, good := range goods.AllGoods {
		numerator := math.Pow(p.Goods[good].Share/float64(w.prices[good]), p.Elasticity)
		demand := numerator * float64(w.wage) / denominator
		w.demand[good] = market.Size(math.Floor(demand))
	}
}

func (w *Worker) placeOrders(p *Parameters) {
	// Workers will always work.
	p.LabourMarket.Post(&market.Order{
		Price: w.wage,
		Size:  1,
		Side:  market.Sell,
		Owner: w,
	})

	for _, good := range goods.AllGoods {
		if w.demand[good] == 0 {
			continue
		}

		p.Goods[good].Market.Post(&market.Order{
			Price: w.prices[good],
			Size:  w.demand[good],
			Side:  market.Buy,
			Owner: w,
		})
	}
}

func (w *Worker) reset() {
	w.unemployed = true
	for _, good := range goods.AllGoods {
		w.purchasesMade[good] = 0
	}
}

// OnFill is triggered when the worker is hired.
func (w *Worker) OnFill(good goods.Good, side market.Side, wage market.Price, size market.Size) {
	if good == goods.Labour {
		w.wage = wage
		w.unemployed = false
	} else {
		w.purchasesMade[good] += size
	}
}

// OnUnfilled is triggered at the end of the cycle if the worker was not hired.
func (w *Worker) OnUnfilled(good goods.Good, side market.Side, size market.Size) {
	// Nothing needed here.
}

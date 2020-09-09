package agents

import (
	"log"

	"github.com/robbrit/econerra/market"
)

var _ = log.Println

// A Worker is an agent that sells labour.
type Worker struct {
	unemployed bool
	wage       market.Price
}

// NewWorker creates a new worker.
func NewWorker(initialWage market.Price) *Worker {
	return &Worker{true, initialWage}
}

// Act triggers the worker's decision process.
func (w *Worker) Act(p *Parameters) {
	if p.LabourMarket.Volume() == 0 {
		// This is the first cycle, just use the wage we have.
	} else if w.unemployed {
		// I was unemployed last round, undercut the market.
		w.wage = p.LabourMarket.Low() - p.Increment
	} else if w.wage <= p.LabourMarket.High() {
		// I was employed last round, ask for a higher price if I can.
		w.wage += p.Increment
	}

	p.LabourMarket.Post(&market.Order{w.wage, 1, market.Sell, w})
}

// OnFill is triggered when the worker is hired.
func (w *Worker) OnFill(side market.Side, wage market.Price, size market.Size) {
	w.wage = wage
	w.unemployed = false
}

// OnUnfilled is triggered at the end of the cycle if the worker was not hired.
func (w *Worker) OnUnfilled(market.Side, market.Size) {
	w.unemployed = true
}

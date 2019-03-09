package econerra

import (
	"log"
)

var _ = log.Println

// A Worker is an agent that sells labour.
type Worker struct {
	unemployed bool
	wage       Price
}

// NewWorker creates a new worker.
func NewWorker(initialWage Price) *Worker {
	return &Worker{true, initialWage}
}

// Act triggers the worker's decision process.
func (w *Worker) Act(p *Parameters) {
	if p.LabourMarket.Volume() == 0 {
		// This is the first cycle, just use the wage we have.
	} else if w.unemployed {
		// I was unemployed last round, undercut the market.
		w.wage = p.LabourMarket.Low() - p.Increment
	} else {
		// I was employed last round, ask for a higher price.
		w.wage += p.Increment
	}

	p.LabourMarket.Post(&MarketOrder{w.wage, 1, Sell, w})
}

// OnFill is triggered when the worker is hired.
func (w *Worker) OnFill(side Side, wage Price, size Size) {
	w.wage = wage
	w.unemployed = false
}

// OnUnfilled is triggered at the end of the cycle if the worker was not hired.
func (w *Worker) OnUnfilled(Side, Size) {
	w.unemployed = true
}

package agents

import (
	"math"

	"github.com/robbrit/econerra/market"
)

// A Firm is an agent responsible for buying labour and producing goods.
type Firm struct {
	wage          market.Price
	workersHired  market.Size
	targetWorkers market.Size
}

// NewFirm creates a new firm with the given production parameters.
func NewFirm(initialWage market.Price) *Firm {
	return &Firm{initialWage, 0, 0}
}

// TargetWorkers gets the number of workers that this firm is trying to hire
// this period.
func (f *Firm) TargetWorkers() market.Size { return f.targetWorkers }

// Act triggers the firm's decision process.
func (f *Firm) Act(p *Parameters) {
	if f.targetWorkers > 0 {
		// We've past the first round, so now do the actual calculation.
		if f.workersHired < f.targetWorkers {
			// Didn't hire enough people, offer a better wage than the market.
			f.wage = p.LabourMarket.High() + p.Increment
		} else if f.wage >= p.LabourMarket.Low() {
			// Got enough people, lower wages if possible
			f.wage -= p.Increment
		}
	}

	// We've calculated our wage, now give labour demand.
	target := math.Pow(float64(f.wage)/p.Tech/p.Scale, 1.0/(p.Scale-1.0))

	if profits(p, f.wage, math.Ceil(target)) > profits(p, f.wage, math.Floor(target)) {
		f.targetWorkers = market.Size(math.Ceil(target))
	} else {
		f.targetWorkers = market.Size(math.Floor(target))
	}
	f.workersHired = 0

	p.LabourMarket.Post(&market.Order{f.wage, f.targetWorkers, market.Buy, f})
}

// profits calculates how much profit a firm makes given a wage and target labour.
func profits(p *Parameters, wage market.Price, labour float64) float64 {
	return p.Tech*math.Pow(labour, p.Scale) - float64(wage)*labour
}

// OnFill is triggered when the firm hires a worker.
func (f *Firm) OnFill(side market.Side, wage market.Price, size market.Size) {
	f.workersHired++
}

// OnUnfilled is triggered if the firm doesn't hire enough workers.
func (f *Firm) OnUnfilled(market.Side, market.Size) {
	// No need to actually do anything here.
}

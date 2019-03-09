package econerra

import "math"

// A Firm is an agent responsible for buying labour and producing goods.
type Firm struct {
	wage          Price
	workersHired  Size
	targetWorkers Size
}

// NewFirm creates a new firm with the given production parameters.
func NewFirm(initialWage Price) *Firm {
	return &Firm{initialWage, 0, 0}
}

// TargetWorkers gets the number of workers that this firm is trying to hire
// this period.
func (f *Firm) TargetWorkers() Size { return f.targetWorkers }

// Act triggers the firm's decision process.
func (f *Firm) Act(p *Parameters) {
	if f.targetWorkers > 0 {
		// We've past the first round, so now do the actual calculation.
		if f.workersHired < f.targetWorkers {
			// Didn't hire enough people, offer a better wage than the market.
			f.wage = p.LabourMarket.High() + p.Increment
		} else {
			// Got enough people, lower wages.
			f.wage -= p.Increment
		}
	}

	// We've calculated our wage, now give labour demand.
	target := math.Pow(float64(f.wage)/p.Tech/p.Scale, 1.0/(p.Scale-1.0))

	if profits(p, f.wage, math.Ceil(target)) > profits(p, f.wage, math.Floor(target)) {
		f.targetWorkers = Size(math.Ceil(target))
	} else {
		f.targetWorkers = Size(math.Floor(target))
	}
	f.workersHired = 0

	p.LabourMarket.Post(&MarketOrder{f.wage, f.targetWorkers, Buy, f})
}

// profits calculates how much profit a firm makes given a wage and target labour.
func profits(p *Parameters, wage Price, labour float64) float64 {
	return p.Tech*math.Pow(labour, p.Scale) - float64(wage)*labour
}

// OnFill is triggered when the firm hires a worker.
func (f *Firm) OnFill(side Side, wage Price, size Size) {
	f.workersHired++
}

// OnUnfilled is triggered if the firm doesn't hire enough workers.
func (f *Firm) OnUnfilled(Side, Size) {
	// No need to actually do anything here.
}

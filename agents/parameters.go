package agents

import (
	"github.com/robbrit/econerra/market"
)

// Parameters is a structure of simulation-wide parameters that agents use to make calculations.
type Parameters struct {
	Increment    market.Price   // Agents' undercutting factor.
	Tech         float64 // Cobb-Douglas technology factor.
	Scale        float64 // Cobb-Douglas returns to scale.
	LabourMarket market.Market  // Where agents can buy/sell labour.
}

package agents

import (
	"github.com/robbrit/econerra/goods"
	"github.com/robbrit/econerra/market"
)

// GoodParameters defines the various parameters that are specific to a single good.
type GoodParameters struct {
	// Cobb-Douglas production technology factor.
	Tech float64
	// Cobb-Douglas production scale factor.
	Scale float64
	// CES utility share factor.
	Share float64
	// Where agents can buy this good.
	Market market.Market
}

// Parameters is a structure of simulation-wide parameters that agents use to make calculations.
type Parameters struct {
	// How much agents will adjust their price each iteration.
	Increment market.Price
	// Where agents can buy/sell labour.
	LabourMarket market.Market
	// CES elasticity of substition parameter
	Elasticity float64

	Goods map[goods.Good]GoodParameters
}

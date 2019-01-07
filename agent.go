package econerra

import (
	"log"
)

var _ = log.Printf

// An AgentStrategy determines how an agent should behave in different
// circumstances.
type AgentStrategy interface {
	// Is this agent a buyer for this good?
	IsBuyer(Good) bool
}

// AgentOptions are options used to construct an agent.
type AgentOptions struct {
	// Adjust is how much an agent adjusts their prices each period.
	Adjust float64
	// InitPrices is how much we should charge for each good in the first iteration.
	InitPrices map[Good]float64
	// InitCash is the amount of cash the agent has in the first iteration.
	InitCash float64
}

type agent struct {
	world      *World
	inventory  []Size
	cash       float64
	lastSignal []MarketSignal
	prices     []float64
	adjust     float64
	strategy   AgentStrategy
}

func newAgent(w *World, opts AgentOptions, as AgentStrategy) *agent {
	a := &agent{
		world:      w,
		inventory:  make([]Size, NumGoods),
		cash:       opts.InitCash,
		lastSignal: make([]MarketSignal, NumGoods),
		prices:     make([]float64, NumGoods),
		adjust:     opts.Adjust,
		strategy:   as,
	}

	for g, p := range opts.InitPrices {
		a.prices[g] = p
	}
	return a
}

func (a *agent) onFill(g Good, side Side, p float64, s Size, sig MarketSignal) {
	if side == Buy {
		a.inventory[g] += s
		a.cash -= float64(s) * p
	} else {
		a.inventory[g] -= s
		a.cash += float64(s) * p
	}

	a.lastSignal[g] = sig
}

func (a *agent) onUnfilled(g Good, side Side, p float64, s Size, sig MarketSignal) {
	a.lastSignal[g] = sig
}

// Act is called when it is the agent's turn to act.
func (a *agent) act() {
	for i := range a.prices {
		g := Good(i)
		dir := -1.0
		if a.strategy.IsBuyer(g) {
			dir = 1.0
		}

		switch a.lastSignal[g] {
		case SignalStrong:
			// Too strong, adjust away from market
			a.prices[i] -= a.adjust * dir
		case SignalWeak, SignalFairUnfilled:
			// Too weak, adjust towards the market
			a.prices[i] += a.adjust * dir
		}
		if a.prices[i] <= 0.0 {
			a.prices[i] = a.adjust
		}
	}
}

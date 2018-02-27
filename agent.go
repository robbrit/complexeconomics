package econerra

// An AgentStrategy determines how an agent should behave in different
// circumstances.
type AgentStrategy interface {
	// Is this agent a buyer for this good?
	IsBuyer(Good) bool
}

type agent struct {
	inventory  []Size
	cash       Price
	lastSignal []MarketSignal
	prices     []Price
	adjust     Price
	strategy   AgentStrategy
}

func newAgent(adjust Price, as AgentStrategy) *agent {
	return &agent{
		inventory:  make([]Size, NumGoods),
		cash:       0,
		lastSignal: make([]MarketSignal, NumGoods),
		prices:     make([]Price, NumGoods),
		adjust:     adjust,
		strategy:   as,
	}
}

func (a *agent) onFill(g Good, side Side, p Price, s Size, sig MarketSignal) {
	if side == BuyOrder {
		a.inventory[g] += s
		a.cash -= Price(s) * p
	} else {
		a.inventory[g] -= s
		a.cash += Price(s) * p
	}

	a.lastSignal[g] = sig
}

func (a *agent) onUnfilled(g Good, side Side, p Price, s Size, sig MarketSignal) {
	a.lastSignal[g] = sig
}

// Act is called when it is the agent's turn to act.
func (a *agent) act() {
	for i := range a.prices {
		g := Good(i)
		dir := -1
		if a.strategy.IsBuyer(g) {
			dir = 1
		}

		switch a.lastSignal[g] {
		case SignalStrong:
			// Too strong, adjust away from market
			a.prices[i] -= Price(int(a.adjust) * dir)
		case SignalWeak:
			// Too weak, adjust towards the market
			a.prices[i] += Price(int(a.adjust) * dir)
		}
	}
}

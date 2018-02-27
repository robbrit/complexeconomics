package econerra

// An Individual is an individual in the econerra world.
type Individual struct {
	agent

	shares []Size
}

// NewIndividual constructs a new individual.
func NewIndividual(initialShares []Size, adjust Price) *Individual {
	i := &Individual{
		shares: initialShares,
	}

	i.agent = *newAgent(adjust, i)
	return i
}

// Act causes the invidual to make its decisions for this cycle.
func (ind *Individual) Act() {
	ind.agent.act()

	// TODO: Figure out how much of each good we want based on our prices
}

// OnFill is an implementation of MarketAgent.OnFill.
func (ind *Individual) OnFill(g Good, side Side, p Price, s Size, sig MarketSignal) {
	ind.agent.onFill(g, side, p, s, sig)
}

// OnUnfilled is an implementation of MarketAgent.OnUnfilled.
func (ind *Individual) OnUnfilled(g Good, side Side, p Price, s Size, sig MarketSignal) {
	ind.agent.onUnfilled(g, side, p, s, sig)
}

// IsBuyer is an implementation of AgentStrategy.IsBuyer.
func (ind *Individual) IsBuyer(g Good) bool { return g != Labour }

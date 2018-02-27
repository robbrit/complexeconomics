package econerra

// FirmID is a unique identifier for a firm.
type FirmID uint64

// A Firm is a firm in the econerra world.
type Firm struct {
	agent

	id        FirmID
	totShares Size
}

// NewFirm constructs a new firm.
func NewFirm(id FirmID, totShares Size, adjust Price) *Firm {
	f := &Firm{
		id:        id,
		totShares: totShares,
	}

	f.agent = *newAgent(adjust, f)
	return f
}

// Act causes the firm to make its decisions for this cycle.
func (f *Firm) Act() {
	f.agent.act()

	// TODO: figure out what prices and wages we want to offer
}

// OnFill is an implementation of MarketAgent.OnFill.
func (f *Firm) OnFill(g Good, side Side, p Price, s Size, sig MarketSignal) {
	f.agent.onFill(g, side, p, s, sig)
}

// OnUnfilled is an implementation of MarketAgent.OnUnfilled.
func (f *Firm) OnUnfilled(g Good, side Side, p Price, s Size, sig MarketSignal) {
	f.agent.onUnfilled(g, side, p, s, sig)
}

// IsBuyer is an implementation of AgentStrategy.IsBuyer.
func (f *Firm) IsBuyer(g Good) bool {
	// TODO - some firm types buy input goods
	return g == Labour
}

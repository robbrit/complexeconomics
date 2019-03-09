package econerra

// A Side represents the side that an order is on (buy vs. sell)
//go:generate stringer -type=Side
type Side uint8

// A Price is how much it costs to buy a good.
type Price uint32

// A Size is a quantity of a good.
type Size uint32

const (
	// Buy is an order to buy things.
	Buy Side = iota
	// Sell is an order to sell things.
	Sell
)

// Market represents a market for buying and selling goods.
type Market interface {
	// Post an order to this market.
	Post(*MarketOrder)
	// Reset the market.
	Reset()
	// Get the highest price for unfilled buy orders.
	Bid() Price
	// Get the lowest price for unfilled sell orders.
	Ask() Price
	// Get the high of the last trading period.
	High() Price
	// Get the low of the last trading period.
	Low() Price
	// Get the volume of goods traded on this market in the last trading period.
	Volume() Size
}

// A MarketAgent is an agent that trades in the market, and can be notified of
// market events.
type MarketAgent interface {
	// OnFill is triggered when an order is filled.
	OnFill(Side, Price, Size)
	// OnUnfilled is called when the market is reset and order has not been filled.
	OnUnfilled(Side, Size)
}

// A MarketOrder is an order to trade something in the market for a given price.
type MarketOrder struct {
	Price Price
	Size  Size
	Side  Side
	Owner MarketAgent
}

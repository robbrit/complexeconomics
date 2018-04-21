package econerra

import (
	"container/heap"
)

// Size represents an order size.
type Size int64

// A MarketSignal tells how "good" an order is.
type MarketSignal uint8

const (
	// SignalWeak means that the order does not have much chance of being filled.
	SignalWeak MarketSignal = iota
	// SignalFair means the order is roughly at market, and has a chance of being
	// filled.
	SignalFair
	// SignalStrong means the order is very good, and will definitely get filled.
	SignalStrong
)

// A Side represents the side that an order is on (buy vs. sell)
type Side uint8

const (
	// Buy is an order to buy things.
	Buy Side = iota
	// Sell is an order to sell things.
	Sell
)

// Market represents a market for buying and selling goods.
type Market interface {
	// Determine what good is traded on this market.
	Good() Good
	// Post an order to this market.
	Post(*MarketOrder)
	// Clear the market.
	Clear()
	// Get the highest price for unfilled buy orders.
	Bid() float64
	// Get the lowest price for unfilled sell orders.
	Ask() float64
	// Get the last trade price on this market.
	Last() float64
	// Get the volume of goods traded on this market in the last cycle.
	Volume() Size
}

type marketImpl struct {
	world *World
	// The type of good that is sold in this market.
	good   Good
	orders []*MarketOrder
	bid    float64
	ask    float64
	last   float64
	volume Size
}

// A MarketAgent is an agent that trades in the market, and can be notified of
// market events.
type MarketAgent interface {
	// OnFill is triggered when an order is filled.
	OnFill(Good, Side, float64, Size, MarketSignal)
	// OnUnfilled is called when the market is cleared and order has not been
	// filled.
	OnUnfilled(Good, Side, float64, Size, MarketSignal)
}

// A MarketOrder is an order to trade something in the market for a given price.
type MarketOrder struct {
	Price float64
	Size  Size
	Side  Side
	Owner MarketAgent
}

// NewMarket constructs a new market for a given good.
func NewMarket(w *World, g Good) Market {
	return &marketImpl{w, g, nil, 0.0, 0.0, 0.0, 0}
}

func (m *marketImpl) Good() Good    { return m.good }
func (m *marketImpl) Bid() float64  { return m.bid }
func (m *marketImpl) Ask() float64  { return m.ask }
func (m *marketImpl) Last() float64 { return m.last }
func (m *marketImpl) Volume() Size  { return m.volume }

// Post adds an order to the market. Note that this order will not get filled
// right away, until the market is cleared.
func (m *marketImpl) Post(o *MarketOrder) {
	if o.Size == 0 {
		return
	}
	if o.Price <= 0.0 {
		return
	}
	m.orders = append(m.orders, o)
}

// Clear clears the market, by determining which orders get filled and which
// are not. Notifications are sent to the owners of each order.
func (m *marketImpl) Clear() {
	// Go through orders in random order.
	bids := orderMaxHeap{}
	offers := orderMinHeap{}
	heap.Init(&bids)
	heap.Init(&offers)

	type fill struct {
		buyOwner  MarketAgent
		sellOwner MarketAgent
		buyPrice  float64
		sellPrice float64
		price     float64
		size      Size
	}

	fills := []*fill{}
	for _, i := range m.world.Rand().Perm(len(m.orders)) {
		order := m.orders[i]

		switch order.Side {
		case Buy:
			if len(offers) == 0 || order.Price < offers[0].Price {
				heap.Push(&bids, order)
				continue
			}

			// Pop sell orders off the heap until we have filled the entire amount.
			size := order.Size
			for len(offers) > 0 && order.Price >= offers[0].Price && size > 0 {
				if offers[0].Size <= size {
					sell := heap.Pop(&offers).(*MarketOrder)
					fills = append(fills, &fill{order.Owner, sell.Owner, order.Price, sell.Price, sell.Price, sell.Size})
					size -= sell.Size
				} else {
					sell := offers[0]
					fills = append(fills, &fill{order.Owner, sell.Owner, order.Price, sell.Price, sell.Price, size})
					offers[0].Size -= size
					size = 0
				}
			}

			if size > 0 {
				order.Size = size
				heap.Push(&bids, order)
			}
		case Sell:
			if len(bids) == 0 || order.Price > bids[0].Price {
				heap.Push(&offers, order)
				continue
			}

			// Pop buy orders off the heap until we have filled the entire amount.
			size := order.Size
			for len(bids) > 0 && order.Price <= bids[0].Price && size > 0 {
				if bids[0].Size <= size {
					buy := heap.Pop(&bids).(*MarketOrder)
					fills = append(fills, &fill{buy.Owner, order.Owner, buy.Price, order.Price, buy.Price, buy.Size})
					size -= buy.Size
				} else {
					buy := bids[0]
					fills = append(fills, &fill{buy.Owner, order.Owner, buy.Price, order.Price, buy.Price, size})
					bids[0].Size -= size
					size = 0
				}
			}

			if size > 0 {
				order.Size = size
				heap.Push(&offers, order)
			}
		}
	}

	// Market is cleared now, send notifications to all agents.
	// Anything remaining did not get filled, and gets an unfilled notification.
	// Use this to calculate remaining bid/ask.
	if len(bids) > 0 {
		m.bid = bids[0].Price
		for _, o := range bids {
			s := SignalWeak
			if o.Price == m.bid {
				s = SignalFair
			}
			o.Owner.OnUnfilled(m.good, Buy, o.Price, o.Size, s)
		}
	}
	if len(offers) > 0 {
		m.ask = offers[0].Price
		for _, o := range offers {
			s := SignalWeak
			if o.Price == m.ask {
				s = SignalFair
			}
			o.Owner.OnUnfilled(m.good, Sell, o.Price, o.Size, s)
		}
	}

	// Anything that was filled gets a fill notification.
	m.volume = 0
	for _, f := range fills {
		p := f.price
		bs := SignalFair
		ss := SignalFair

		// If the price above the ask, then the buyer was weak.
		if p > m.ask {
			bs = SignalWeak
			ss = SignalStrong
		} else if p < m.bid || m.bid == 0 {
			bs = SignalStrong
			ss = SignalWeak
		}
		f.buyOwner.OnFill(m.good, Buy, p, f.size, bs)
		f.sellOwner.OnFill(m.good, Sell, p, f.size, ss)

		m.volume += f.size
		m.last = p
	}

	// Reset the market every cycle.
	m.orders = nil
}

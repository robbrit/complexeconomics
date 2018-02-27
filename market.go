package econerra

import (
	"container/heap"
	"math/rand"
)

// Price represents a price.
type Price int64

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
	// BuyOrder is an order to buy things.
	BuyOrder Side = iota
	// SellOrder is an order to sell things.
	SellOrder
)

// Market represents a market for buying and selling goods.
type Market interface {
	Good() Good
	Post(*MarketOrder)
	Clear()
}

type sequenceGenerator func(int) []int

type marketImpl struct {
	// The type of good that is sold in this market.
	good   Good
	orders []*MarketOrder
	seq    sequenceGenerator
}

// A MarketAgent is an agent that trades in the market, and can be notified of
// market events.
type MarketAgent interface {
	// OnFill is triggered when an order is filled.
	OnFill(Good, Side, Price, Size, MarketSignal)
	// OnUnfilled is called when the market is cleared and order has not been
	// filled.
	OnUnfilled(Good, Side, Price, Size, MarketSignal)
}

// A MarketOrder is an order to trade something in the market for a given price.
type MarketOrder struct {
	Price Price
	Size  Size
	Side  Side
	Owner MarketAgent
}

// NewMarket constructs a new market for a given good.
func NewMarket(g Good) Market {
	return &marketImpl{g, nil, rand.Perm}
}

// Good gives the good that is traded in this market.
func (m *marketImpl) Good() Good { return m.good }

// Post adds an order to the market. Note that this order will not get filled
// right away, until the market is cleared.
func (m *marketImpl) Post(o *MarketOrder) {
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
		buyPrice  Price
		sellPrice Price
		price     Price
		size      Size
	}

	fills := []*fill{}
	for _, i := range m.seq(len(m.orders)) {
		order := m.orders[i]

		switch order.Side {
		case BuyOrder:
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
		case SellOrder:
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
	// Anything that was filled gets a fill notification.
	for _, f := range fills {
		p := f.price
		bs := SignalFair
		if f.buyPrice > p {
			bs = SignalStrong
		}
		ss := SignalFair
		if f.sellPrice < p {
			ss = SignalStrong
		}
		f.buyOwner.OnFill(m.good, BuyOrder, p, f.size, bs)
		f.sellOwner.OnFill(m.good, SellOrder, p, f.size, ss)
	}

	// Anything remaining did not get filled, and gets an unfilled notification
	if len(bids) > 0 {
		bid := bids[0].Price
		for _, o := range bids {
			s := SignalWeak
			if o.Price == bid {
				s = SignalFair
			}
			o.Owner.OnUnfilled(m.good, BuyOrder, o.Price, o.Size, s)
		}
	}
	if len(offers) > 0 {
		ask := offers[0].Price
		for _, o := range offers {
			s := SignalWeak
			if o.Price == ask {
				s = SignalFair
			}
			o.Owner.OnUnfilled(m.good, SellOrder, o.Price, o.Size, s)
		}
	}

	// Reset the market every cycle.
	m.orders = nil
}

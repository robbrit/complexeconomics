package econerra

import (
	"testing"
)

type fakeAgent struct {
	fillGood   Good
	fillPrice  Price
	fillSize   Size
	fillSide   Side
	fillSignal MarketSignal

	unfilledGood   Good
	unfilledPrice  Price
	unfilledSize   Size
	unfilledSide   Side
	unfilledSignal MarketSignal
}

func (fa *fakeAgent) OnFill(g Good, s Side, p Price, q Size, sig MarketSignal) {
	fa.fillGood = g
	fa.fillPrice = p
	fa.fillSize = q
	fa.fillSide = s
	fa.fillSignal = sig
}

func (fa *fakeAgent) OnUnfilled(g Good, s Side, p Price, q Size, sig MarketSignal) {
	fa.unfilledGood = g
	fa.unfilledPrice = p
	fa.unfilledSize = q
	fa.unfilledSide = s
	fa.unfilledSignal = sig
}

func TestMarket(t *testing.T) {
	// Situation: one seller, three buyers - one buyer high, one medium, one low.
	b1 := &fakeAgent{}
	b2 := &fakeAgent{}
	b3 := &fakeAgent{}

	s := &fakeAgent{}

	m := NewMarket(Meat).(*marketImpl)
	m.seq = func(int) []int { return []int{0, 1, 2, 3} }

	m.Post(&MarketOrder{10.0, 100, SellOrder, s})
	m.Post(&MarketOrder{12.0, 10, BuyOrder, b1})
	m.Post(&MarketOrder{10.0, 200, BuyOrder, b2})
	m.Post(&MarketOrder{8.0, 1000, BuyOrder, b3})

	m.Clear()

	for _, test := range []struct {
		desc      string
		agent     *fakeAgent
		wantAgent *fakeAgent
	}{
		{
			"high price should get a strong signal",
			b1,
			&fakeAgent{Meat, 10.0, 10, BuyOrder, SignalStrong, 0, 0.0, 0, 0, 0},
		},
		{
			"mid price should get a fair signal",
			b2,
			&fakeAgent{Meat, 10.0, 90, BuyOrder, SignalFair,
				Meat, 10.0, 110, BuyOrder, SignalFair},
		},
		{
			"low price should get a weak signal",
			b3,
			&fakeAgent{0, 0.0, 0, 0, 0, Meat, 8.0, 1000, BuyOrder, SignalWeak},
		},
		{
			"sell should have latest fill values",
			s,
			&fakeAgent{Meat, 10.0, 90, SellOrder, SignalFair, 0, 0.0, 0, 0, 0},
		},
	} {
		if *test.agent != *test.wantAgent {
			t.Errorf("%s: got %v, want %v", test.desc, test.agent, test.wantAgent)
		}
	}
}

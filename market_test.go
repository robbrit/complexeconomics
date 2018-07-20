package econerra

import (
	"math/rand"
	"testing"
)

type fakeAgent struct {
	fillGood    Good
	fillfloat64 float64
	fillSize    Size
	fillSide    Side
	fillSignal  MarketSignal

	unfilledGood    Good
	unfilledfloat64 float64
	unfilledSize    Size
	unfilledSide    Side
	unfilledSignal  MarketSignal
}

func (fa *fakeAgent) OnFill(g Good, s Side, p float64, q Size, sig MarketSignal) {
	fa.fillGood = g
	fa.fillfloat64 = p
	fa.fillSize = q
	fa.fillSide = s
	fa.fillSignal = sig
}

func (fa *fakeAgent) OnUnfilled(g Good, s Side, p float64, q Size, sig MarketSignal) {
	fa.unfilledGood = g
	fa.unfilledfloat64 = p
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

	w := NewWorld(WorldOptions{
		Rand: rand.New(NewFakeRand([]int64{0, 1, 2, 3})),
	})

	m := NewMarket(w, Meat).(*marketImpl)

	m.Post(&MarketOrder{10.0, 100, Sell, s})
	m.Post(&MarketOrder{12.0, 10, Buy, b1})
	m.Post(&MarketOrder{10.0, 200, Buy, b2})
	m.Post(&MarketOrder{8.0, 1000, Buy, b3})

	m.Clear()

	for _, test := range []struct {
		desc      string
		agent     *fakeAgent
		wantAgent *fakeAgent
	}{
		{
			"high price should get a strong signal",
			b1,
			&fakeAgent{Meat, 10.0, 10, Buy, SignalStrong, 0, 0.0, 0, 0, 0},
		},
		{
			"mid price should get a fair signal",
			b2,
			&fakeAgent{Meat, 10.0, 90, Buy, SignalFair,
				Meat, 10.0, 110, Buy, SignalFair},
		},
		{
			"low price should get a weak signal",
			b3,
			&fakeAgent{0, 0.0, 0, 0, 0, Meat, 8.0, 1000, Buy, SignalWeak},
		},
		{
			"sell should have latest fill values",
			s,
			&fakeAgent{Meat, 10.0, 90, Sell, SignalFair, 0, 0.0, 0, 0, 0},
		},
	} {
		if *test.agent != *test.wantAgent {
			t.Errorf("%s: got %v, want %v", test.desc, test.agent, test.wantAgent)
		}
	}
}

package econerra

import (
	"testing"
)

type fakeAgent struct {
	fillPrice Price
	fillSize  Size
	fillSide  Side

	unfilledSize Size
	unfilledSide Side
}

func (fa *fakeAgent) OnFill(s Side, p Price, q Size) {
	fa.fillPrice = p
	fa.fillSize = q
	fa.fillSide = s
}

func (fa *fakeAgent) OnUnfilled(s Side, q Size) {
	fa.unfilledSize = q
	fa.unfilledSide = s
}

func TestMarket(t *testing.T) {
	// Situation: one seller, three buyers - one buyer high, one medium, one low.
	b1 := &fakeAgent{}
	b2 := &fakeAgent{}
	b3 := &fakeAgent{}

	s := &fakeAgent{}

	m := NewDoubleAuctionMarket()

	m.Post(&MarketOrder{10, 100, Sell, s})
	m.Post(&MarketOrder{12, 10, Buy, b1})
	m.Post(&MarketOrder{10, 200, Buy, b2})
	m.Post(&MarketOrder{8, 1000, Buy, b3})

	m.Reset()

	for _, test := range []struct {
		desc      string
		agent     *fakeAgent
		wantAgent *fakeAgent
	}{
		{
			"high buy should get filled at 10",
			b1,
			&fakeAgent{10, 10, Buy, 0, 0},
		},
		{
			"mid buy should get partially filled",
			b2,
			&fakeAgent{10, 90, Buy, 110, Buy},
		},
		{
			"low buy should not get filled at all",
			b3,
			&fakeAgent{0, 0, 0, 1000, Buy},
		},
		{
			"sell should have latest fill values",
			s,
			&fakeAgent{10, 90, Sell, 0, 0},
		},
	} {
		if *test.agent != *test.wantAgent {
			t.Errorf("%s: got %v, want %v", test.desc, test.agent, test.wantAgent)
		}
	}
}

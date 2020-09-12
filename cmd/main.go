package main

import (
	"log"
	"math/rand"

	"github.com/robbrit/econerra/agents"
	"github.com/robbrit/econerra/goods"
	"github.com/robbrit/econerra/market"
)

const (
	numWorkers = 5000
	numCycles  = 40
	initWage   = 25
	initPrice  = 10
	increment  = 1
	randSeed   = 123456
	elasticity = 0.8
	// All goods will use the same scale.
	scale = 0.5
)

var (
	numFirms = map[goods.Good]int{
		goods.Grain:      100,
		goods.Vegetables: 100,
		goods.Meat:       100,
	}
	technology = map[goods.Good]float64{
		goods.Grain:      100.0,
		goods.Vegetables: 70.0,
		goods.Meat:       10.0,
	}
	shares = map[goods.Good]float64{
		goods.Grain:      1.0,
		goods.Vegetables: 0.8,
		goods.Meat:       5.0,
	}
)

type actor interface {
	Act(p *agents.Parameters)
}

func main() {
	log.Printf("Starting simulation...\n")

	markets := []market.Market{}

	params := agents.Parameters{
		Increment:    increment,
		LabourMarket: market.NewDoubleAuction(goods.Labour),
		Goods:        map[goods.Good]agents.GoodParameters{},
	}
	markets = append(markets, params.LabourMarket)

	for _, good := range goods.AllGoods {
		params.Goods[good] = agents.GoodParameters{
			Tech:   technology[good],
			Scale:  scale,
			Share:  shares[good],
			Market: market.NewDoubleAuction(good),
		}
		markets = append(markets, params.Goods[good].Market)
	}

	var actors []actor
	for _, good := range goods.AllGoods {
		for i := 0; i < numFirms[good]; i++ {
			actors = append(actors, agents.NewFirm(good, initWage, initPrice))
		}
	}
	for i := 0; i < numWorkers; i++ {
		actors = append(actors, agents.NewWorker(initWage, initPrice))
	}

	r := rand.New(rand.NewSource(randSeed))

	for i := 0; i < numCycles; i++ {
		p := r.Perm(len(actors))
		for _, i := range p {
			actors[i].Act(&params)
		}
		for _, mkt := range markets {
			mkt.Reset()
		}

		/*demand := market.Size(0)
		for _, a := range actors {
			switch f := a.(type) {
			case *agents.Firm:
				demand += f.TargetWorkers()
			}
		}*/

		// TODO(rob): Dump to CSV file for processing with a more advanced tool like Pandas.
		//log.Printf("%3d Bid\tAsk\tLow\tHigh\tVolume\tLDemand", i)
		//log.Printf("%3d %d\t%d\t%d\t%d\t%d\t%d", i, mkt.Bid(), mkt.Ask(), mkt.Low(), mkt.High(), mkt.Volume(), demand)
	}
}

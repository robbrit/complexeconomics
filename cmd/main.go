package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/robbrit/econerra/agents"
	"github.com/robbrit/econerra/goods"
	"github.com/robbrit/econerra/market"
)

const (
	numWorkers = 1000
	numCycles  = 100
	initWage   = 100
	initPrice  = 2
	increment  = 1
	randSeed   = 123456
	elasticity = 0.8
	// All goods will use the same scale.
	scale = 0.5
)

var (
	numFirms = map[goods.Good]int{
		goods.Grain:      5,
		goods.Vegetables: 5,
		goods.Meat:       15,
	}
	technology = map[goods.Good]float64{
		goods.Grain:      1000.0,
		goods.Vegetables: 800.0,
		goods.Meat:       500.0,
	}
	shares = map[goods.Good]float64{
		goods.Grain:      2.0,
		goods.Vegetables: 1.0,
		goods.Meat:       5.0,
	}
)

type actor interface {
	Act(*agents.Parameters, int)
	TargetDemand(goods.Good) market.Size
	TargetSupply(goods.Good) market.Size
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

	filename := "output.csv"
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Unable to open CSV file %s for writing: %s", filename, err)
	}
	defer f.Close()
	w := csv.NewWriter(f)

	w.Write([]string{
		"Iteration",
		"Good",
		"Bid",
		"Ask",
		"Low",
		"High",
		"Volume",
		"Supply",
		"Demand",
	})

	for i := 0; i < numCycles; i++ {
		p := r.Perm(len(actors))
		for _, i := range p {
			actors[i].Act(&params, i)
		}
		for _, mkt := range markets {
			mkt.Reset()

			supply := market.Size(0)
			demand := market.Size(0)
			for _, a := range actors {
				supply += a.TargetSupply(mkt.Good())
				demand += a.TargetDemand(mkt.Good())
			}

			w.Write([]string{
				fmt.Sprintf("%d", i),
				fmt.Sprintf("%s", mkt.Good()),
				fmt.Sprintf("%d", mkt.Bid()),
				fmt.Sprintf("%d", mkt.Ask()),
				fmt.Sprintf("%d", mkt.Low()),
				fmt.Sprintf("%d", mkt.High()),
				fmt.Sprintf("%d", mkt.Volume()),
				fmt.Sprintf("%d", supply),
				fmt.Sprintf("%d", demand),
			})
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

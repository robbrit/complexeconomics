package main

import (
	"log"
	"math"
	"math/rand"

	"github.com/robbrit/econerra/agents"
	"github.com/robbrit/econerra/market"
)

const (
	numWorkers = 100
	numFirms   = 20
	numCycles  = 40
	initWage   = 25
	increment  = 1
	technology = 100.0
	scale      = 0.5
	randSeed   = 123456
)

type actor interface {
	Act(p *agents.Parameters)
}

func main() {
	log.Printf("Starting simulation...\n")

	var actors []actor
	mkt := market.NewDoubleAuction()

	params := agents.Parameters{
		Increment:    increment,
		Tech:         technology,
		Scale:        scale,
		LabourMarket: mkt,
	}

	w := technology * scale * math.Pow(float64(numWorkers)/float64(numFirms), scale-1.0)
	l := float64(numFirms) * math.Pow(w/technology/scale, 1.0/(scale-1.0))

	log.Printf("Expected market wage: %f", w)
	log.Printf("Total expected demand: %f", l)

	for i := 0; i < numFirms; i++ {
		actors = append(actors, agents.NewFirm(initWage))
	}

	for i := 0; i < numWorkers; i++ {
		actors = append(actors, agents.NewWorker(initWage))
	}

	r := rand.New(rand.NewSource(randSeed))

	for i := 0; i < numCycles; i++ {
		p := r.Perm(len(actors))
		for _, i := range p {
			actors[i].Act(&params)
		}
		mkt.Reset()

		demand := market.Size(0)
		for _, a := range actors {
			switch f := a.(type) {
			case *agents.Firm:
				demand += f.TargetWorkers()
			}
		}

		log.Printf("%3d Bid\tAsk\tLow\tHigh\tVolume\tLDemand", i)
		log.Printf("%3d %d\t%d\t%d\t%d\t%d\t%d", i, mkt.Bid(), mkt.Ask(), mkt.Low(), mkt.High(), mkt.Volume(), demand)
	}
}

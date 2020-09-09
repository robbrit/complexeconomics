package main

import (
	"log"
	"math"
	"math/rand"

	e "github.com/robbrit/econerra"
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
	Act(p *e.Parameters)
}

func main() {
	log.Printf("Starting simulation...\n")

	var actors []actor
	mkt := e.NewDoubleAuctionMarket()

	params := e.Parameters{
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
		actors = append(actors, e.NewFirm(initWage))
	}

	for i := 0; i < numWorkers; i++ {
		actors = append(actors, e.NewWorker(initWage))
	}

	r := rand.New(rand.NewSource(randSeed))

	for i := 0; i < numCycles; i++ {
		p := r.Perm(len(actors))
		for _, i := range p {
			actors[i].Act(&params)
		}
		mkt.Reset()

		demand := e.Size(0)
		for _, a := range actors {
			switch f := a.(type) {
			case *e.Firm:
				demand += f.TargetWorkers()
			}
		}

		log.Printf("%3d Bid\tAsk\tLow\tHigh\tVolume\tLDemand", i)
		log.Printf("%3d %d\t%d\t%d\t%d\t%d\t%d", i, mkt.Bid(), mkt.Ask(), mkt.Low(), mkt.High(), mkt.Volume(), demand)
	}
}

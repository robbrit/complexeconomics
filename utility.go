package econerra

import (
	"math"
)

// All goods within a substituteGroup are substitutes for one another.
type substituteGroup []struct {
	g Good
	s float64
}
// A complementGroup contains a set of substitute groups, all goods in
// substitute group i are complementary to goods in subsitute group j, i != j.
type complementGroup []struct {
	g substituteGroup
	s float64
}
// A goodGroup divides goods into groups that are distinct from one another
// with respect to substitutes.
type goodGroup []struct {
	g complementGroup
	s float64
}

// UtilityGoods are goods that provide utility for individuals.
var UtilityGoods = []Good{Vegetables, Meat, Beer, Clothing}

var utilityGroups = goodGroup{
	// Clothing Group
	{
		g: complementGroup{
			{
				g: substituteGroup{{Clothing, 1.0}},
				s: 1.0,
			},
		},
		s: 0.3,
	},
	// Edibles
	{
		g: complementGroup{
			// Drinks
			{
				g: substituteGroup{
					{Beer, 1.0},
				},
				s: 0.2,
			},
			// Foods
			{
				g: substituteGroup{
					{Meat, 0.7}, {Vegetables, 0.3},
				},
				s: 0.8,
			},
		},
		s: 0.7,
	},
}

const (
	compExp = -1.125 / 0.125
	subsExp = 0.5
	topExp = 0.5
)

func utility(q []float64) float64 {
	u := 0.0
	for _, g := range utilityGroups {
		gs := 0.0
		for _, cg := range g.g {
			cs := 0.0
			for _, sg := range cg.g {
				cs += sg.s * math.Pow(q[sg.g], subsExp)
			}
			gs += cg.s * math.Pow(cs, compExp / subsExp)
		}
		u += g.s * math.Pow(gs, topExp / compExp)
	}
	return u
}

package econerra

type input struct {
	// Which input good is used.
	g Good
	// How much of it is used to produce 1 unit of output.
	s Size
}

var inputGoods = map[Good][]input{
	Meat:     {{Grain, 10}},
	Beer:     {{Grain, 3}},
	Clothing: {{Cotton, 4}},
}

package goods

// A Good is something that a human consumes to gain utility.
//go:generate stringer -type=Good
type Good uint8

const (
	Grain Good = iota
	Vegetables
	Meat
	Labour
)

var AllGoods = []Good{Grain, Vegetables, Meat}

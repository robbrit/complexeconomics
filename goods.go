package econerra

// A Good is a unique type of good available in this market.
// go:generate stringer -type=Good
type Good uint8

const (
	// "Regular" goods
	Grain Good = iota
	Vegetables
	Cotton
	Meat
	Wool
	Beer
	Clothing

	// "Special" goods
	Labour

	NumGoods
)

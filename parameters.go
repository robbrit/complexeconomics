package econerra

// Parameters is a structure of simulation-wide parameters that agents use to make calculations.
type Parameters struct {
	Increment    Price   // Agents' undercutting factor.
	Tech         float64 // Cobb-Douglas technology factor.
	Scale        float64 // Cobb-Douglas returns to scale.
	LabourMarket Market  // Where agents can buy/sell labour.
}

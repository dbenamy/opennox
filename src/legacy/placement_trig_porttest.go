//go:build porttest

package legacy

// These observers call the same native helpers as the placement owners.
func portTestPlacementSin(x float64) float64 { return placementSin(x) }
func portTestPlacementCos(x float64) float64 { return placementCos(x) }

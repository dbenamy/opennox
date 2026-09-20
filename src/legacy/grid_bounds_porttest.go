//go:build porttest

package legacy

import (
	"github.com/opennox/libs/types"
	"math"
)

// PortTestGridBounds uses a null grid to assert rejection occurs before access.
// Inputs that expose the original C bug are run only in isolated child processes.
func PortTestGridBounds(xBits, yBits uint32) int {
	old := worldTileGrid
	worldTileGrid = nil
	defer func() { worldTileGrid = old }()
	return int(tileAtPoint(types.Pointf{X: math.Float32frombits(xBits), Y: math.Float32frombits(yBits)}))
}

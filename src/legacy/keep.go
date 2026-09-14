//go:build none

package legacy

import "C"

var (
	// maybe worth keeping mapgen
	_ = mapOrchestrationStartAlt
	// unused grip spell effect - could be useful
	_ = C.nox_xxx_unusedCheckGripEffect_4E03F0
	// some game setup code - unused game mode?
	_ = C.nox_xxx_unused_418840
)

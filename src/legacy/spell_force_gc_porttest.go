//go:build porttest

package legacy

import "C"
import "runtime"

var portTestSpellForceCollections int

//export portTestSpellForceCollect
func portTestSpellForceCollect() {
	runtime.GC()
	portTestSpellForceCollections++
}

func PortTestSpellForceCollections() int { return portTestSpellForceCollections }

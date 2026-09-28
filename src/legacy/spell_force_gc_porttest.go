//go:build porttest

package legacy

import "runtime"

var portTestSpellForceCollections int

func portTestSpellForceCollect() {
	runtime.GC()
	portTestSpellForceCollections++
}

func PortTestSpellForceCollections() int { return portTestSpellForceCollections }

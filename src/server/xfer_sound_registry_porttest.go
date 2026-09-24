//go:build porttest

package server

import "unsafe"

func PortTestXferSoundRegistry(name string, sound bool) unsafe.Pointer {
	table := xferFuncs
	if sound {
		table = damageSoundFuncs
	}
	p, ok := table[name]
	if !ok {
		panic("missing registry name: " + name)
	}
	return p
}

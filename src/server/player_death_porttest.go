//go:build porttest

package server

import "unsafe"

// Read the actual registration used when parsing a PlayerDie object definition.
func PortTestPlayerDeathCallback() unsafe.Pointer { return deathFuncs["PlayerDie"].Func }

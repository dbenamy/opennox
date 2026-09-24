//go:build porttest

package server

import (
	"github.com/opennox/libs/types"
)

// Identical original raw-pointer cases now exercise the production pointer API.
func PortTestCollisionWith(u, target *Object, normal *types.Pointf) {
	u.CallCollideWith(target, normal)
}

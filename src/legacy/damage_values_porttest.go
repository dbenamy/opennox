//go:build porttest

package legacy

import "github.com/opennox/opennox/v1/server"

// Enter the unchanged integer-result owner before and after typed dispatch.
func PortTestProjectileDamage(target, source, weapon *server.Object, amount, kind int32) int32 {
	return projectileDamage(target, source, weapon, amount, kind)
}

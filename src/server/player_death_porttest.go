//go:build porttest

package server

// Resolve the registered typed owner without changing the victim's Death slot.
// Native registry keys are opaque identities, not executable C callbacks.
func PortTestPlayerDeathCallback() DeathFunc {
	return objDeath.Get(deathFuncs["PlayerDie"].Func)
}

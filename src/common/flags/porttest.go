//go:build porttest

package noxflags

// PortTestGameFlags isolates flag queries without invoking unrelated engine
// lifecycle hooks. The caller must restore state and must not run in parallel.
func PortTestGameFlags(value GameFlag) func() {
	old := game
	game = value
	return func() { game = old }
}

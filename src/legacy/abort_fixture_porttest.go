//go:build porttest

package legacy

// PortTestInvalidFixtureDispatch reaches the existing fatal fallback without
// unrelated fixture setup. It is called only in a subprocess contract.
func PortTestInvalidFixtureDispatch(paint bool) {
	if paint {
		paintInvokeNative(-1, [6]uint32{})
	} else {
		mapRoomPortInvoke(-1, [4]uint32{})
	}
}

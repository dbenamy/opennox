//go:build porttest

package blobdata

// PortTestSummonConstants returns the production slot order and text offsets.
// They contain integers only; string pointers are owned separately by the fixture.
func PortTestSummonConstants() []byte {
	return append([]byte(nil), data587000[184456:184552]...)
}

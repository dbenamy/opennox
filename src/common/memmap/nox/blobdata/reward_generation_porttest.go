//go:build porttest

package blobdata

// PortTestRewardGoldConstants returns the original three gold scaling doubles.
func PortTestRewardGoldConstants() []byte { return append([]byte(nil), data581450[10248:10272]...) }

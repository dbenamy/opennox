//go:build porttest

package blobdata

// PortTestStatisticsTags returns the shipped report field names, including
// writable player suffixes. Do not qualify against zero-filled blob state.
func PortTestStatisticsTags() []byte { return append([]byte(nil), data587000[71480:72016]...) }

//go:build porttest

package blobdata

// PortTestWorldTrigConstants returns the shipped trigonometric table scales.
func PortTestWorldTrigConstants() []byte { return append([]byte(nil), data581450[7184:7208]...) }

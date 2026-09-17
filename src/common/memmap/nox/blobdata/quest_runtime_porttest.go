//go:build porttest

package blobdata

// Shipped score exponent and adjacent bytes, including the full original C
// long-double read width. This exposes data only, not a scoring implementation.
func PortTestQuestScoreConstant() []byte {
	return append([]byte(nil), data581450[10088:10100]...)
}

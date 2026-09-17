//go:build porttest

package blobdata

func PortTestCollisionCoreNormals() []byte {
	return append([]byte(nil), data587000[289928:289960]...)
}

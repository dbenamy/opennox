//go:build porttest

package blobdata

// Production light-radius coefficients normally installed by blob startup.
func PortTestClientParticleLightTable() []byte {
	return append([]byte(nil), data587000[154972:154984]...)
}

//go:build porttest

package blobdata

// PortTestFloorFacadeData supplies the actual startup names and pointer-table
// bytes. The fixture relocates the same six entries as blob_init.go.
func PortTestFloorFacadeData() []byte {
	return append([]byte(nil), data587000[26488:26660]...)
}

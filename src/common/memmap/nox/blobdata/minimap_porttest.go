//go:build porttest

package blobdata

func PortTestMinimapTables() []PortTestClientEffectsTable {
	return []PortTestClientEffectsTable{{Base: 0x587000, Offset: 149240, Data: append([]byte(nil), data587000[149240:149272]...)}}
}

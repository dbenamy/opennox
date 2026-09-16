//go:build porttest

package blobdata

func PortTestWorldWallTables() []PortTestClientEffectsTable {
	return []PortTestClientEffectsTable{
		{Base: 0x587000, Offset: 149364, Data: append([]byte(nil), data587000[149364:149428]...)},
		{Base: 0x587000, Offset: 85440, Data: append([]byte(nil), data587000[85440:85680]...)},
	}
}

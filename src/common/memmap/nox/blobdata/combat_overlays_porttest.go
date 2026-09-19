//go:build porttest

package blobdata

// Original console format and attachment direction vectors; fixtures own copies.
func PortTestCombatOverlayTables() []PortTestClientEffectsTable {
	var out []PortTestClientEffectsTable
	for _, r := range [][2]int{{161668, 12}, {194136, 2048}} {
		out = append(out, PortTestClientEffectsTable{Base: 0x587000, Offset: uintptr(r[0]), Data: append([]byte(nil), data587000[r[0]:r[0]+r[1]]...)})
	}
	return out
}

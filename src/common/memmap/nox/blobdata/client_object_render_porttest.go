//go:build porttest

package blobdata

func PortTestObjectRenderTables() []PortTestClientEffectsTable {
	var out []PortTestClientEffectsTable
	for _, r := range [][2]int{{185464, 4}, {196184, 256}} {
		out = append(out, PortTestClientEffectsTable{Base: 0x587000, Offset: uintptr(r[0]), Data: append([]byte(nil), data587000[r[0]:r[0]+r[1]]...)})
	}
	return out
}

//go:build porttest

package blobdata

func PortTestClientPresentationTables() []PortTestClientEffectsTable {
	var out []PortTestClientEffectsTable
	for _, r := range [][2]int{{151208, 64}, {151304, 128}, {161776, 72}, {163576, 36}, {194136, 2048}} {
		out = append(out, PortTestClientEffectsTable{Base: 0x587000, Offset: uintptr(r[0]), Data: append([]byte(nil), data587000[r[0]:r[0]+r[1]]...)})
	}
	return out
}

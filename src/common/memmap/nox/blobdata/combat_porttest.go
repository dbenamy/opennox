//go:build porttest

package blobdata

// PortTestCombatTables returns copies of the runtime direction and facing tables.
func PortTestCombatTables() map[uintptr][]byte {
	out := make(map[uintptr][]byte)
	for off, n := range map[uintptr]int{192088: 2048, 194136: 2048, 202504: 576} {
		out[off] = append([]byte(nil), data587000[int(off):int(off)+n]...)
	}
	return out
}

//go:build porttest

package blobdata

func PortTestWorldGeometryTables() map[uintptr][]byte {
	out := map[uintptr][]byte{}
	for off, n := range map[uintptr]int{192088: 2048, 194136: 2048, 230056: 36, 292496: 11, 292520: 264, 292784: 32} {
		out[off] = append([]byte(nil), data587000[off:off+uintptr(n)]...)
	}
	return out
}

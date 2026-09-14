//go:build porttest

package blobdata

func PortTestMapPopulationTables() map[uintptr][]byte {
	out := map[uintptr][]byte{}
	for off, n := range map[uintptr]int{254688: 16, 255032: 8} {
		out[off] = append([]byte(nil), data587000[off:off+uintptr(n)]...)
	}
	return out
}

//go:build porttest

package blobdata

func PortTestMapGrowthTables() map[uintptr][]byte {
	out := map[uintptr][]byte{}
	for off, n := range map[int]int{197812: 24, 197836: 24, 197924: 16} {
		out[uintptr(off)] = append([]byte(nil), data587000[off:off+n]...)
	}
	return out
}

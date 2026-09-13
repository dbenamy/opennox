//go:build porttest

package blobdata

// PortTestMapPaintingTables supplies immutable startup direction and edge data.
func PortTestMapPaintingTables() map[uintptr][]byte {
	out := map[uintptr][]byte{}
	for off, n := range map[uintptr]int{71276: 200, 230052: 40, 255052: 64, 282736: 576} {
		out[off] = append([]byte(nil), data587000[off:off+uintptr(n)]...)
	}
	return out
}

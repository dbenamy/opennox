//go:build porttest

package blobdata

// Numeric tables used by the real creature serializers. Pointer/name tables
// are relocated separately by the fixture, as they are during runtime setup.
func PortTestCreatureXferTables() map[uintptr][]byte {
	out := make(map[uintptr][]byte)
	for _, r := range [][2]uintptr{{66000, 66116}, {192088, 194136}, {230056, 230092}, {255604, 255604 + 72*16}, {262056, 262072}} {
		out[r[0]] = append([]byte(nil), data587000[r[0]:r[1]]...)
	}
	return out
}

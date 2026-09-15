//go:build porttest

package blobdata

// PortTestInventoryDisplayGeometry returns the actual paper-doll and tray hit
// rectangles from the embedded asset-independent game data.
func PortTestInventoryDisplayGeometry() []byte {
	return append([]byte(nil), data587000[136192:136384]...)
}

// PortTestInventoryWindowGeometry includes the two inventory event exclusion rectangles.
func PortTestInventoryWindowGeometry() []byte {
	return append([]byte(nil), data587000[136192:136416]...)
}

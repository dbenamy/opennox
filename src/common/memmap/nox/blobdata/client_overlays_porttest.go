//go:build porttest

package blobdata

// PortTestClientOverlayTables supplies the shipped names and exact startup pointer
// bindings without initializing unrelated mutable game state.
func PortTestClientOverlayTables() ([]PortTestClientEffectsTable, [][2]uintptr) {
	return []PortTestClientEffectsTable{
		{Base: 0x587000, Offset: 30980, Data: append([]byte(nil), data587000[30980:31008]...)},
		{Base: 0x587000, Offset: 85752, Data: append([]byte(nil), data587000[85752:85776]...)},
	}, [][2]uintptr{{29456, 30980}, {29460, 30988}, {29464, 30996}, {85712, 85752}, {85716, 85764}}
}
